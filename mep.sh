go build &&
ssh pi@192.168.0.25 "
sudo systemctl stop keybr
sudo rm -rf /home/pi/old_keybr
mkdir -p /home/pi/go/src/keybr
sudo mv /home/pi/go/src/keybr /home/pi/old_keybr
"
scp -r /home/g/golangkeybr pi@192.168.0.25:/home/pi/go/src/keybr
ssh pi@192.168.0.25 "
sudo chown -R pi:pi /home/pi/go/src/keybr
cd /home/pi/go/src/keybr
go build .
sudo systemctl start keybr
"
ssh pi@192.168.0.25 journalctl -f -u keybr
