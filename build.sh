CURRENT_DIR=$(pwd)
SSL_DIR=$CURRENT_DIR/nginx/ssl

docker-compose down

cd $SSL_DIR
openssl genrsa 2048 > server.key
openssl req -new -key server.key > server.csr
openssl x509 -days 3650 -req -extfile subjectnames.txt -signkey server.key < server.csr > server.crt

cd $CURRENT_DIR
docker-compose build
docker-compose up -d
