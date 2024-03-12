FROM golang:1.22.0-alpine3.19 AS builder
ARG ENV

RUN apk update && apk upgrade --available && \
    apk add make && \
    adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "10001" \
    "auth"

WORKDIR /opt/app/
COPY . .

RUN go mod download && go mod verify
RUN make build-app ENV=${ENV}
RUN chown -R auth:auth ./

FROM scratch
ARG CONFIG

WORKDIR /opt/app/
COPY --from=builder /opt/app/bin/main .
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /opt/app/${CONFIG} ./config
COPY --from=builder /opt/app/tls/ ./tls/

USER auth:auth

CMD ["./main", "-config=./config"]
