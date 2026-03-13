package main

import (
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ListarDiscosEstatico simula escanear los discos /dev/sdX en Linux Kiosk
func ListarDiscosEstatico() []string {
	// Implementación real recorrería filepath.Glob("/dev/sd*") o sysfs 
	// para omitir particiones loop de ventoy.
	rutas := []string{"/dev/sda", "/dev/sdb"}
	var validos []string
	for _, ruta := range rutas {
		if _, err := os.Stat(ruta); err == nil {
			validos = append(validos, ruta)
		}
	}
	if len(validos) == 0 {
		return []string{"Seleccionar Disco..."}
	}
	return validos
}

func main() {
	aplicacion := app.New()
	ventana := aplicacion.NewWindow("Conversor Forense Bidireccional de Discos")
	ventana.Resize(fyne.NewSize(600, 300))

	// Título Minimalista
	titulo := widget.NewLabelWithStyle("Básico ⮂ Dinámico", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	estadoLbl := widget.NewLabel("Estado: Listo para escanear.")
	feedbackLbl := widget.NewLabel("")

	discosDisponibles := ListarDiscosEstatico()
	selectorDisco := widget.NewSelect(discosDisponibles, func(s string) {
		estadoLbl.SetText("Disco seleccionado: " + s)
	})
	selectorDisco.SetSelectedIndex(0)

	btnConvertir := widget.NewButton("Convertir", func() {
		discoObjetivo := selectorDisco.Selected
		if discoObjetivo == "" || discoObjetivo == "Seleccionar Disco..." {
			dialog.ShowError(fmt.Errorf("por favor, seleccione un disco válido"), ventana)
			return
		}

		dialog.ShowConfirm("Confirmación Forense", 
			fmt.Sprintf("¿Está seguro de convertir %s? Se creará rollback en RAM.", discoObjetivo),
			func(b bool) {
				if b {
					go func() {
						feedbackLbl.SetText("Haciendo backup en RAM...")
						time.Sleep(1 * time.Second)
						
						feedbackLbl.SetText("Aplicando cambios en LBA 0 y LDM...")
						// Lógica hipotética llamando a pkg/disco y pkg/transaccion
						time.Sleep(1 * time.Second)

						feedbackLbl.SetText("Verificando hash (Check-Read)...")
						time.Sleep(1 * time.Second)

						// Simulando el exito:
						// feedbackLbl.SetText("Error E/S: Rollback Aplicado")
						dialog.ShowInformation("Éxito", "Conversión completada sin pérdida de datos.", ventana)
						feedbackLbl.SetText("Completado.")
					}()
				}
			}, ventana)
	})

	caja := container.NewVBox(
		titulo,
		widget.NewLabel("Detección Automática de Hardware (Linux sysfs):"),
		selectorDisco,
		estadoLbl,
		feedbackLbl,
		btnConvertir,
	)

	ventana.SetContent(caja)

	// Configuración Kiosk - Opcional en Desktop, pero Fyne soporta fullscreen
	// ventana.SetFullScreen(true)

	ventana.ShowAndRun()
}
