ARG SWAGGER_VERSION
FROM quay.io/goswagger/swagger:$SWAGGER_VERSION

WORKDIR /data
ARG UID
ARG GID
COPY script.sh /generate/swagger/script.sh
VOLUME ["/data"]

RUN apk add bash jq && addgroup -g "$GID" -S docker && adduser -u "$UID" -S user -G docker && \
    chmod +x /generate/swagger/script.sh && \
    chown -R "${UID}:${GID}" /data

USER user
ENTRYPOINT ["/generate/swagger/script.sh"]
