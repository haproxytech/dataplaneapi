DATAPLANEAPI_PATH = ${PWD}
	
update:
	cd ${DATAPLANEAPI_PATH} && go get -v -insecure -u -f

build:
	make update
	cd ${DATAPLANEAPI_PATH}/cmd/dataplaneapi && go build -a

