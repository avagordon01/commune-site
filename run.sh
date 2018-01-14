set -e

GOPATH=$(go env GOPATH)

go install
sudo /usr/bin/setcap cap_net_bind_service=+ep $GOPATH/bin/site
$GOPATH/bin/site
