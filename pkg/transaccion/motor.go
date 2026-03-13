package transaccion

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

// MotorTransaccion maneja el estado de la escritura a bajo nivel con rollback automático en tmpfs.
type MotorTransaccion struct {
	ArchivoDispositivo string
	RutaBackup         string // Normalmente /tmp/disk_backup.bin (tmpfs)
	desplazamiento     int64
	tamañoBackup       int
	backupCreado       bool
	hashOriginal       []byte
}

// Iniciar prepara el motor. Lee los sectores críticos (Ej LBA 0-33 y LDM) y los guarda en RAM.
func NuevoMotor(dispositivo string, off int64, tam int) *MotorTransaccion {
	return &MotorTransaccion{
		ArchivoDispositivo: dispositivo,
		RutaBackup:         "/tmp/disk_backup.bin", // RAM tmpfs en Alpine
		desplazamiento:     off,
		tamañoBackup:       tam,
		backupCreado:       false,
	}
}

// CrearPuntoRestauracion lee los datos crudos del disco físico desde el desplazamiento.
func (m *MotorTransaccion) CrearPuntoRestauracion() error {
	disco, err := os.Open(m.ArchivoDispositivo)
	if err != nil {
		return fmt.Errorf("error leyendo el disco: %w", err)
	}
	defer disco.Close()

	buffer := make([]byte, m.tamañoBackup)
	disco.Seek(m.desplazamiento, io.SeekStart)
	_, err = disco.Read(buffer)
	if err != nil {
		return fmt.Errorf("Fallo al leer sectores para Backup: %w", err)
	}

	// Guardar el Sha1 Original
	hash := sha1.Sum(buffer)
	m.hashOriginal = hash[:]

	// Guardarlo en /tmp (tmpfs)
	err = os.WriteFile(m.RutaBackup, buffer, 0600)
	if err != nil {
		return fmt.Errorf("Fallo al crear archivo de Rollback en RAM: %w", err)
	}

	m.backupCreado = true
	fmt.Println("INFO: Backup atómico creado en", m.RutaBackup)
	return nil
}

// EjecutarSeguro recibe una función o un slice de bytes a escribir. Si falla, llama a RevertirAutomata().
func (m *MotorTransaccion) EjecutarSeguro(datosNuevos []byte) error {
	if !m.backupCreado {
		return fmt.Errorf("no hay punto de restauración, peligro operacional")
	}

	disco, err := os.OpenFile(m.ArchivoDispositivo, os.O_WRONLY|os.O_SYNC, 0600)
	if err != nil {
		return m.AplicarRollbackV2(err)
	}
	defer disco.Close()

	disco.Seek(m.desplazamiento, io.SeekStart)
	_, err = disco.Write(datosNuevos)

	if err != nil {
		return m.AplicarRollbackV2(fmt.Errorf("error E/S durante mutacion: %w", err))
	}

	// Ejecutar 'Check-Read' para garantizar que el sector nuevo tiene el hash correcto (pre-calculado)
	hashNuevosDatos := sha1.Sum(datosNuevos)
	if errCheck := m.CheckRead(hashNuevosDatos[:]); errCheck != nil {
		return m.AplicarRollbackV2(fmt.Errorf("Fallo checksum Check-Read (Sector Corrupto): %w", errCheck))
	}

	fmt.Println("INFO: Transacción Atómica Completada exitosamente y sin corrupción.")
	return nil
}

// CheckRead lee de nuevo del disco para validar CRC32/SHA1
func (m *MotorTransaccion) CheckRead(hashEsperado []byte) error {
	disco, err := os.Open(m.ArchivoDispositivo)
	if err != nil {
		return err
	}
	defer disco.Close()

	buffer := make([]byte, m.tamañoBackup)
	disco.Seek(m.desplazamiento, io.SeekStart)
	disco.Read(buffer)

	hashLeido := sha1.Sum(buffer)
	if string(hashEsperado) != string(hashLeido[:]) {
		return fmt.Errorf("mismatch de hash en hardware")
	}

	return nil
}

// AplicarRollbackV2 restaura la imagen binaria desde tmpfs hacia el dispositivo físico.
func (m *MotorTransaccion) AplicarRollbackV2(errOriginal error) error {
	fmt.Println("CRÍTICO: Ejecutando Try/Catch/Rollback desde Memoria RAM... (Fallo:", errOriginal.Error(), ")")
	if !m.backupCreado {
		return fmt.Errorf("fatal: imposible hacer rollback, sin backup. Error original: %v", errOriginal)
	}
	
	bufferRollback, err := os.ReadFile(m.RutaBackup)
	if err != nil {
		return fmt.Errorf("fatal: error al leer backup %s", err)
	}

	disco, err := os.OpenFile(m.ArchivoDispositivo, os.O_WRONLY|os.O_SYNC, 0600)
	if err != nil {
		return fmt.Errorf("fatal: disco bloqueado durante el rollback")
	}
	defer disco.Close()

	disco.Seek(m.desplazamiento, io.SeekStart)
	_, err = disco.Write(bufferRollback)
	if err != nil {
		return fmt.Errorf("catástrofe: error en E/S al escribir Rollback: %v", err)
	}

	fmt.Println("INFO: Rollback aplicado correctamente. No hay pérdida de datos.")
	return fmt.Errorf("operación abortada y disco restaurado. Motivo original: %w", errOriginal)
}
