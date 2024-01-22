FROM ubuntu:jammy

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y curl && \
    curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | bash && \
    apt-get update && \
    apt-get install -y migrate

RUN mkdir -p /usr/var/migrations
COPY ./migrations /usr/var/migrations
COPY ./migrate.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/migrate.sh

CMD ["/usr/local/bin/migrate.sh"]
