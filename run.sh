go install || exit 1
sudo /usr/bin/setcap cap_net_bind_service=+ep /srv/go/bin/commune
$GOPATH/bin/commune
