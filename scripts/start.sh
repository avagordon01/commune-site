clear
go install
sudo setcap 'cap_net_bind_service=+ep' $GOPATH/bin/commune
$GOPATH/bin/commune
