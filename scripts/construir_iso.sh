#!/bin/bash
set -e

echo "--- INICIANDO CONSTRUCCIÓN PROFESIONAL DEL ISO FORENSE (DEBIAN LIVE) ---"

# 1. Limpieza
rm -rf build_space
mkdir -p build_space
cd build_space

# 2. Configuración de la ISO (Debian Bookworm con LXDE para la GUI de Go)
lb config \
    --distribution bookworm \
    --binary-images iso-hybrid \
    --architectures amd64 \
    --archive-areas "main contrib non-free non-free-firmware" \
    --iso-application "Conversor-Basico-Dinamico" \
    --iso-volume "CONVERSOR-LIVE" \
    --bootappend-live "boot=live components locales=es_ES.UTF-8 keyboard-layouts=es"

# 3. Lista de paquetes necesarios
mkdir -p config/package-lists/
cat <<EOF > config/package-lists/desktop.list.chroot
lxde-core
xserver-xorg
xinit
sudo
util-linux
parted
# Dependencias para ejecutar el binario de Go con Fyne (X11)
libgl1-mesa-dri
libgl1-mesa-glx
libx11-6
libxcursor1
libxrandr2
libxinerama1
libxi6
EOF

# 4. Inyectar el binario compilado de Go
APP_DIR="config/includes.chroot/usr/local/bin"
mkdir -p "$APP_DIR"
# El binario se compila fuera y se pasa aquí
cp ../bin/conversor_gui "$APP_DIR/"
chmod +x "$APP_DIR/conversor_gui"

# 5. Configurar el Auto-Arranque de la interfaz
AUTOSTART_DIR="config/includes.chroot/etc/xdg/lxsession/LXDE"
mkdir -p "$AUTOSTART_DIR"
cat <<EOF > "$AUTOSTART_DIR/autostart"
@lxpanel --profile LXDE
@pcmanfm --desktop --profile LXDE
@xset s off
@sudo /usr/local/bin/conversor_gui
EOF

# 6. Permisos de Sudo para el usuario Live (necesario para acceso a bloques)
SUDO_DIR="config/includes.chroot/etc/sudoers.d"
mkdir -p "$SUDO_DIR"
echo "live ALL=(ALL) NOPASSWD: ALL" > "$SUDO_DIR/live"
chmod 0440 "$SUDO_DIR/live"

# 7. Ejecutar compilación
echo "Compilando Sistema Operativo Live..."
lb build

# 8. Mover resultado
if [ -f live-image-amd64.hybrid.iso ]; then
    mv live-image-amd64.hybrid.iso ../Conversor_Forense.iso
    echo "--- ¡ÉXITO! ISO GENERADA EN LA RAÍZ ---"
else
    echo "ERROR: Fallo al generar la ISO."
    exit 1
fi
