FROM alpine:3.8 as builder

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN adduser -D -g '' app

# ---
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

COPY git-dc-up /git-dc-up

USER app

ENTRYPOINT ["/git-dc-up"]
