ARG SWAGGER_VERSION

FROM golang:1.25 AS golang

FROM quay.io/goswagger/swagger:0.32.3
COPY --from=golang /usr/local/go /usr/local/go

ENV PATH="$PATH:/usr/local/go/bin"

WORKDIR /data
ARG UID
ARG GID
COPY script.sh /generate/swagger/script.sh
COPY templates /generate/swagger/templates
VOLUME ["/data"]

RUN apk add bash jq && sed -i "s/999/99/" /etc/group &&  \
    addgroup -g "$GID" -S docker && adduser -u "$UID" -S user -G docker && \
    chmod +x /generate/swagger/script.sh && \
    chown -R "${UID}:${GID}" /data

USER user
ENTRYPOINT ["/generate/swagger/script.sh"]
