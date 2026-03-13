# Conversor Forense Atómico: Básico a Dinámico y Viceversa

## 📂 Arquitectura de Directorios Entregada
El código ha sido configurado bajo los estándares modernos de Go y empaquetado optimizado para Alpine Linux y Ventoy:

```text
/Conversor de Básico a Dinámico y viceversa
│
├── go.mod                                # Archivo principal del módulo Go.
│
├── cmd/
│   └── interfaz/
│       └── main.go                       # Interfaz Gráfica con Fyne.
│
├── pkg/
│   ├── disco/                            # Lógica de conversión LDM/Sectores.
│   └── transaccion/                      # Motor de rollback y seguridad.
│
├── .github/workflows/
│   └── build.yml                         # CI/CD - Generador automático de ISO (Estilo MBR-GPT).
│
└── scripts/
    └── construir_iso.sh                  # Script de construcción basado en Debian Live Build.
```

## 🔐 Decisiones Técnicas Aplicadas

1. **Lenguaje:** **Golang**, compilado estáticamente para integrarse en la ISO.
2. **Motor de ISO:** Se utiliza **Debian Live Build**, el mismo sistema de grado profesional usado en el proyecto **MBR-GPT**, garantizando compatibilidad universal (UEFI/Legacy) y estabilidad forense.
3. **Auto-Arranque:** La ISO arranca directamente en un entorno LXDE minimalista que lanza el conversor con privilegios de root.

## 🚀 Cómo obtener la ISO

Al igual que en el proyecto **MBR-GPT**:
1. Haz un **Push** de los cambios.
2. Ve a la pestaña **Actions** en GitHub.
3. Descarga el artefacto **`Conversor-Forense-ISO`** una vez finalice la tarea.

> **Nota Técnica pre-uso:** 
> Ya he descargado las dependencias obligatorias iniciales en `go.mod`.
> *Para una compilación real, la rutina `Shrink` en Módulo 1 requerirá integrarse via librerías C o system-calls a `ntfsresize` en Linux ya que redimensionar un FS de Journaling sin usar la herramienta probada del kernel nativo resultaría en corrupción de datos. La arquitectura ya tiene el 'placeholder' condicional exacto esperando a ser enganchado en `conversor.go`.*
