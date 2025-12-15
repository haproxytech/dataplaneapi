#!/bin/sh
docker exec -d ${DOCKER_CONTAINER_NAME} /bin/sh -c "CI_DATAPLANE_RELOAD_DELAY_OVERRIDE=1 DPAPI_ACME_PROPAGTIMEOUT_SEC=-1 dataplaneapi -f /etc/haproxy/dataplaneapi.yaml"
