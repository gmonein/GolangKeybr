go build &&
ssh pi@192.168.0.25 "
sudo systemctl stop keybr
rm -rf /home/pi/old_golangkeybr
mkdir -p /home/pi/golangkeybr
mv /home/pi/golangkeybr /home/pi/old_golangkeybr
"
scp -r /home/g/golangkeybr pi@192.168.0.25:/home/pi/
ssh pi@192.168.0.25 "
sudo chown -R pi:pi /home/pi/golangkeybr
cd /home/pi/golangkeybr
go build .
sudo systemctl start keybr
"
