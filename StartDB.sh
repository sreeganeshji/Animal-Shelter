docker-compose down --remove-orphans
docker-compose up

# convenience commands below to copy back-end/ to your $GOPATH
cp -a back-end $GOPATH/src