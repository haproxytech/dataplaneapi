#!/bin/bash
docker exec -d ${DOCKER_CONTAINER_NAME} /bin/sh -c "CI_DATAPLANE_RELOAD_DELAY_OVERRIDE=1 dataplaneapi -f /etc/haproxy/dataplaneapi.yaml --maps-dir=/tmp/maps"
