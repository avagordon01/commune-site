go install || exit 1
pkill commune
sudo /usr/bin/setcap cap_net_bind_service=+ep /srv/go/bin/commune
$GOPATH/bin/commune &>log.txt & disown %
