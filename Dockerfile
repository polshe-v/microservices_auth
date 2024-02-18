FROM golang:1.22.0-alpine3.19 AS builder

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
RUN make build

FROM scratch

WORKDIR /opt/app/
COPY --from=builder /opt/app/bin/main .
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

EXPOSE 50000/tcp
USER auth:auth

CMD ["./main"]
