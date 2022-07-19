#!/bin/sh
openssl req -new -subj "/CN=localhost" -newkey rsa:2048 -nodes -keyout ./api.server.key -out ./api.server.csr
openssl x509 -req -days 3650 -in ./api.server.csr -signkey ./api.server.key -out ./api.server.crt

openssl req -new -subj "/CN=localhost" -newkey rsa:2048 -nodes -keyout ./web.server.key -out ./web.server.csr
openssl x509 -req -days 3650 -in ./web.server.csr -signkey ./web.server.key -out ./web.server.crt

rm -rf api.server.csr web.server.csr
# move the *.crt and *.key files to assets/certs/ folder.
# enable below command when running the project on linux.
# mv ./api.server.key ./api.server.crt web.server.key web.server.crt ../assets/certs/