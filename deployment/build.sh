DOMAIN=$1
cd ../server || exit
docker build -t kucar/$DOMAIN -f ../deployment/$DOMAIN/Dockerfile .