FROM alpine:3.20 AS cert-builder

RUN apk add --no-cache openssl

WORKDIR /certs

RUN --mount=type=secret,id=cert_key,target=/tmp/server.key \
    --mount=type=secret,id=cert_crt,target=/tmp/server.crt \
    mkdir -p /etc/nginx-ssl \
    && cp /tmp/server.key /etc/nginx-ssl/server.key \
    && cp /tmp/server.crt /etc/nginx-ssl/server.crt \
    && chmod 600 /etc/nginx-ssl/server.key

FROM nginx:1.27-alpine AS runtime

RUN rm -f /etc/nginx/conf.d/default.conf

COPY conf.d/default.conf /etc/nginx/conf.d/default.conf

COPY --from=cert-builder /etc/nginx-ssl/ /etc/nginx/ssl/

RUN chown -R nginx:nginx /etc/nginx/ssl \
    && chmod 644 /etc/nginx/ssl/server.crt \
    && chmod 600 /etc/nginx/ssl/server.key

EXPOSE 8443

HEALTHCHECK --interval=10s --timeout=3s --retries=5 \
    CMD wget -qO- https://127.0.0.1:8443/healthz --no-check-certificate || exit 1
