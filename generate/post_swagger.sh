#!/bin/bash

set -e
echo -n " ---> adding serverStartedCallback() to "
sed -i "s/wg.Wait()/serverStartedCallback()\nwg.Wait()/g" server.go
go fmt server.go
