#!/bin/sh

while ! (migrate -database postgres://${DB__USER}:${DB__PASSWORD}@${DB__HOST}:${DB__PORT}/${DB__DATABASE}?sslmode=disable -path /usr/var/migrations up)
do
        echo "Migration is not successful. Trying again..."
done
