package disco

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// PrivHead representa la estructura LDM PRIVATE HEADER (PRIVHEAD) (512 bytes).
type PrivHead struct {
	Firma            [8]byte // "PRIVHEAD"
	CheckSum         uint32
	VersionMayor     uint16
	VersionMenor     uint16
	Timestamp        uint64
	SecuenciaLogica  uint64
	IdDisco          string // GUID del disco en LDM
	TOCBlockLBA      uint64 // Apunta al primer TOCBLOCK
	SectoresConfig   uint64 // Tamaño de Base de datos
	BytesSectores    uint64 // Tamaño del sector, útilmente 512
}

// TocBlock representa la estructura LDM Table of Contents (TOCBLOCK) (512 bytes).
type TocBlock struct {
	Firma      [8]byte // "TOCBLOCK"
	CheckSum   uint32
	Secuencia  uint64
	ConfigLBA  uint64 // Inicio de la DB config (VBLK)
	ConfigTamanio uint64 // Tamaño del config section
	LogLBA     uint64 // Journal VBLK
}

// VblkHeader representa la cabecera de las estructuras VBLK (Volumen, Componente, Partición).
type VblkHeader struct {
	Firma      [4]byte // "VBLK"
	NumeroRef  uint32
	GrupoID    uint32
	IndiceFrag uint16
	CantFrag   uint16
}

// AnalizarPrivHead lee un arreglo de 512 bytes y mapea hacia la estructura PrivHead.
func AnalizarPrivHead(datos []byte) (*PrivHead, error) {
	if len(datos) < 512 {
		return nil, fmt.Errorf("los datos son insuficientes para PRIVHEAD")
	}
	
	if string(datos[0:8]) != "PRIVHEAD" {
		return nil, fmt.Errorf("firma PRIVHEAD no encontrada")
	}

	priv := &PrivHead{}
	copy(priv.Firma[:], datos[0:8])
	// Nota: El parseo crudo exacto cambia por la estructura de MS-LDM en el offset correcto (Big Endian)
	// Como buena práctica forense se utiliza BigEndian en LDM para algunas estructuras y LittleEndian en MBR
	buffer := bytes.NewReader(datos[8:])
	binary.Read(buffer, binary.BigEndian, &priv.CheckSum)
	binary.Read(buffer, binary.BigEndian, &priv.VersionMayor)
	binary.Read(buffer, binary.BigEndian, &priv.VersionMenor)
	
	// Saltando campos no documentados y llenando la base.
	// En un diseño de Grado Comercial se calcularía el CheckSum completo para validar integridad.
	return priv, nil
}
