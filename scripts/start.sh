clear
go install
sudo setcap 'cap_net_bind_service=+ep' ~/go/bin/commune
~/go/bin/commune
