FROM alpine:3.19
ARG ENV=$ENV

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.18.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

WORKDIR /opt/app

COPY migrations/*.sql migrations/
COPY migration.sh ./migration.sh
COPY ${ENV}.env ./.env

RUN chmod +x migration.sh

ENTRYPOINT ["bash", "migration.sh"]
