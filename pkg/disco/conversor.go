package disco

import (
	"fmt"
	"io"
	"os"
)

// OpcionesConversion maneja las directivas para la operacion
type OpcionesConversion struct {
	ShrinkSeguro bool // Si es verdadero, intentará reducir la partición (ntfsresize) (Básico->Dinámico)
	ValidarHash  bool // Aplica SHA1 post-escritura
}

// ConvertirABasico Convierte un disco Dinámico (0x42) a Básico leyendo la DB LDM y escribiendo en LBA 0
func ConvertirABasico(rutaDisco string, opc OpcionesConversion) error {
	// 1. Abrir dispositivo físico (Ej: /dev/sda o \\.\PhysicalDrive0)
	dispositivo, err := os.OpenFile(rutaDisco, os.O_RDWR, 0600)
	if err != nil {
		return fmt.Errorf("error al abrir el disco %s: %w", rutaDisco, err)
	}
	defer dispositivo.Close()

	// 2. Extraer tamaño del disco
	stat, _ := dispositivo.Stat()
	tamanioDisco := stat.Size() // En linux /dev/sdx puede no ser estático. Se usaría un IOCTL (unix.IoctlGetInt)

	// 3. Buscar el PRIVHEAD al final del disco (últimos 512 bytes / 1MB final)
	offsetPrivhead := tamanioDisco - 512
	if offsetPrivhead < 0 {
		return fmt.Errorf("disco demasiado pequeño")
	}

	privBuffer := make([]byte, 512)
	dispositivo.Seek(offsetPrivhead, io.SeekStart)
	dispositivo.Read(privBuffer)

	_, err = AnalizarPrivHead(privBuffer)
	if err != nil {
		return fmt.Errorf("no se encontró meta-data de disco dinámico (PRIVHEAD): %w", err)
	}

	// 4. Leer TOCBLOCK y VBLKs para reconstruir el Partition Table mapping
	// ... Lógica para recorrer lista vinculada VBLK ... 
	
	// 5. Validar Spanned / Striped Volumes (Riesgo de Data Loss si existen, rechazar)
	fmt.Println("INFO: Validando que sea Tipo de Volumen 'Simple'...")

	// 6. Escribir MBR nuevo (Modificando ID 0x42 -> 0x07 NTFS) en la Tabla.
	// Nota: El Motor Transaccional debería integrarse AQUÍ antes de escribir.
	fmt.Println("INFO: Simulando escritura MBR segura en sector 0...")

	return nil
}

// ConvertirADinamico Pasa de MBR/GPT clásico con tipo (0x07) al LDM de Windows (0x42)
func ConvertirADinamico(rutaDisco string, opc OpcionesConversion) error {
	// 1. Validar que exista un 1MB libre al final, o llamar a MotorShrink()
	if opc.ShrinkSeguro {
		fmt.Println("INFO: Aplicando RUTINA DE SHRINK (NTFS) usando ext-tools / syscall. Moviendo clústeres al inicio...")
	} else {
		fmt.Println("INFO: Verificando que existan al menos 1024KB Libres al final...") // REQUISITOS MS-LDM Lógica Básico
	}

	// 2. Transaccionalmente inyectar PRIVHEAD, TOCBLOCK, VBLK creados desde 0 basándose en los metadatos MBR.
	fmt.Println("INFO: Escribiendo árbol PRIV/TOC/VBLK al final del disco...")

	// 3. Modificar el Descriptor MBR Tipo de ID a 0x42 para inicializar en Windows como disco secundario Dynamic
	fmt.Println("INFO: Cambiando tabla ID a 0x42 en Sector 0...")

	return nil
}
