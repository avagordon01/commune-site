#for a quicker development loop
#visudo to allow that specific sudo command without a password
#call run.sh in a while loop `while ./run.sh; do true; done`
#then press C-c to restart the server easily

set -e

SRV_BIN=$(go env GOPATH)/bin/site

go install
sudo /usr/bin/setcap cap_net_bind_service=+ep $SRV_BIN
exec $SRV_BIN
