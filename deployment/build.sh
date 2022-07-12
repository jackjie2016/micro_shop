APP=$1
DOMAIN=$2
cd ../$APP/$DOMAIN
docker build -t mxshop/$APP/$DOMAIN -f ./$APP/$DOMAIN/Dockerfile .