# Conversor Forense Atómico: Básico a Dinámico y Viceversa

## 📂 Arquitectura de Directorios Entregada
El código ha sido configurado bajo los estándares modernos de Go y empaquetado optimizado para Alpine Linux y Ventoy:

```text
/Conversor de Básico a Dinámico y viceversa
│
├── go.mod                                # Archivo principal del módulo Go (v1.21).
│
├── cmd/
│   └── interfaz/
│       └── main.go                       # (Módulo 3) - Interfaz Gráfica Minimalista con Fyne.io que llama a los módulos transaccionales.
│
├── pkg/
│   ├── disco/
│   │   ├── estructuras_ldm.go            # (Módulo 1) - Parseador y Writer para la base de datos binaria MS-LDM (PRIVHEAD, TOCBLOCK, VBLK).
│   │   └── conversor.go                  # (Módulo 1) - Lógica Bidireccional. Extracción, evaluación de Spanned/Striped Volumes y Shrink NTFS con sys-calls de Linux.
│   │
│   └── transaccion/
│       └── motor.go                      # (Módulo 2) - Try/Catch/Rollback atómico. Guarda copia de MBR/GPT y LDM a RAM (/tmp/tmpfs) antes de mutar. Ejecuta Check-Reads de hash.
│
└── scripts/
    └── construir_iso_hibrida.sh          # (Módulo 4) - Bash Script Universal (Grub-Mkrescue / UEFI+Legacy BIOS / Ventoy Compliant) para compilar en un ISO auto-ejecutable y aislado en RAM de Alpine.
```

## 🔐 Decisiones Técnicas Aplicadas

1. **Lenguaje:** **Golang**, aprovechando generadores de binarios estáticos cross-platform (útil para el chroot de Alpine) con `os.OpenFile` directo a dispositivos de bloques `/dev/sdX` o `\\.\PhysicalDrive0`.
2. **GUI (Módulo 3):** Se empleó *Fyne.io*. Su compilación es soportada de forma nativa en Alpine mediante `cage` (Kiosk en Wayland). Requerirá compilar en un entorno compatible para vincular las cabeceras `CGO` si compila para Linux.
3. **Manejo LDM (Módulo 1):** Como Windows utiliza orden *Big Endian* para registros específicos en su LDM, el parser se programó implementando `binary.BigEndian` en vez del Little Endian usual de `MBR`. Además previene el *Data Loss* comprobando clústeres extendidos antes del mapeo.
4. **Motor Transaccional Seguro (Módulo 2):** Se diseñó una clase especial `MotorTransaccion{}`. Cualquier inyección a disco pasa por `EjecutarSeguro`. Falla con Check-Read si SHA-1 del buffer nuevo difiere o hay un pico E/S. El *Rollback Automático* dispara una alerta para calmar al perito forense y reestablecer el `disk_backup.bin` originado en RAM (`tmpfs`).
5. **Bootloader Universal (Módulo 4):** Se entregó `construir_iso_hibrida.sh`. Configura Alpine para que omita el login local (`tty1::respawn...cage`), lo que activa el *Modo Quiosco* inamovible (cero distracciones, arranca directo en la GUI).

> **Nota Técnica pre-uso:** 
> Ya he descargado las dependencias obligatorias iniciales en `go.mod`.
> *Para una compilación real, la rutina `Shrink` en Módulo 1 requerirá integrarse via librerías C o system-calls a `ntfsresize` en Linux ya que redimensionar un FS de Journaling sin usar la herramienta probada del kernel nativo resultaría en corrupción de datos. La arquitectura ya tiene el 'placeholder' condicional exacto esperando a ser enganchado en `conversor.go`.*

## 🚀 Compilación Automática (En la Nube)
He preparado el proyecto para que **no tengas que lidiar con cross-compiling en tu Mac ni configurar Linux localmente**. Hemos integrado un flujo de trabajo CI/CD en la carpeta `.github/workflows/build_alpine_iso.yml`.

**Para obtener tu archivo `.iso`:**
1. Sube y haz Push de todo este código a un repositorio en **GitHub**.
2. Dirígete a la pestaña **Actions** en tu repositorio.
3. El proceso `"Build Alpine Hybrid ISO"` se iniciará automáticamente. GitHub compilará internamente la interfaz Go usando dependencias nativas de Linux, correrá el empaquetado y subirá el resultado.
4. Al finalizar, verás en la parte inferior de la ejecución de Action el archivo **`Conversor-Basico-Dinamico-ISO`** listo para descargar en zip. Descomprímelo y cópialo dentro de tu unidad Ventoy para arrancar.
