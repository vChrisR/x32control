#!/bin/bash
echo "Running rpi-update"
sudo rpi-update

echo "Installing dependencies"
sudo apt-get update
sudo apt-get -y install fbi libharfbuzz0b libdouble-conversion1 libglapi-mesa libgles2-mesa libfontconfig1 libinput10 libxkbcommon0 libts-0.0-0 libts-0.0-0

echo "Installing systemd units"
sudo cp splashscreen.service /etc/systemd/system/
sudo cp x32control.service /etc/systemd/system/
sudo mv ./splash.png ../splash.png
sudo systemctl daemon-reload
sudo systemctl enable splashscreen.service
sudo systemctl enable x32control.service


echo "Changing boot configuration"
sudo systemctl mask plymouth-start.service
sudo systemctl disable getty@tty1.service
sudo bash -c "printf 'gpu_mem=256\nlcd_rotate=2\ndisable_splash=1\n' >> /boot/config.txt"
sudo sed -i 's/console=tty1/console=tty3/g' /boot/cmdline.txt
sudo sed -i '$ s/$/ logo.nologo vt.global_cursor_default=0 consoleblank=0 loglevel=3 quiet/' /boot/cmdline.txt

echo "Setup done."
