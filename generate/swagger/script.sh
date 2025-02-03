#!/bin/bash

set -e
GO_VERSION=$(go version)
echo " ---> go version: $GO_VERSION"
SPEC_DIR=$(mktemp -d)
echo " ---> source folder: $SPEC_DIR"
DST_DIR=$(mktemp -d)
echo " ---> generate folder: $DST_DIR"
# see if we have a replace directive
CN_VERSION=$(go mod edit -json | jq -c -r '.Replace | .[] | select(.Old.Path | contains("github.com/haproxytech/client-native/v6")) | .New.Version' 2>/dev/null | awk -F"-" '{print $NF}') || ""
REMOTE_VERSION=$(go mod edit -json | jq -c -r '.Replace | .[] | select(.Old.Path | contains("github.com/haproxytech/client-native/v6")) | .New.Version' 2>/dev/null | awk -F"/" '{print $1}') || ""
if [ "$REMOTE_VERSION" = "null" ]; then
   # we have a local version of CN
   CN_VERSION=$(go mod edit -json | jq -c -r '.Replace | .[] | select(.Old.Path | contains("github.com/haproxytech/client-native/v6")) | .New.Path' 2>/dev/null) || ""
fi
# if hash is to short take all of it (example v1.0.0-dev1)
[ "${#CN_VERSION}" -gt 0 ] && [ "${#CN_VERSION}" -lt 6 ] && CN_VERSION=$(go mod edit -json | jq -c -r '.Replace | .[] | select(.Old.Path | contains("github.com/haproxytech/client-native/v6")) | .New.Version')
# check if version is there, if not, use one from require
[ -z "$CN_VERSION" ] && CN_VERSION=$(go mod edit -json | jq -c -r '.Require | .[] | select(.Path | contains("github.com/haproxytech/client-native/v6")) | .Version' 2>/dev/null)
# check if version contains '-' and if it is a 12 char commit hash then use it, if not use the whole tag
suffix=$(echo $CN_VERSION | awk -F"-" '{print$NF}')
[ -z "${CN_VERSION##*"-"*}" ] && [ ${#suffix} -eq 12 ] && CN_VERSION=$suffix
echo " ---> version of client native used: $CN_VERSION"
# extract repository
REPO_PATH=$(go mod edit -json | jq -r '.Replace | .[] | select(.Old.Path | contains("github.com/haproxytech/client-native/v6")) | .New.Path'  2>/dev/null |  awk -F"/" '{print $2 "/" $3}') || ""
[ -z "$REPO_PATH" ] && REPO_PATH=haproxytech/client-native

# extract url, gitlab and github have different urls to raw content
URL_PATH=$(go mod edit -json | jq -r '.Replace | .[] | select(.Old.Path | contains("github.com/haproxytech/client-native/v6")) | .New.Path' 2>/dev/null |  awk -F"/" '{print $1}') || ""
EXTRA_PATH=""
if [[ $URL_PATH =~ "gitlab" ]]; then
   EXTRA_PATH="-/raw/"
else
  URL_PATH=raw.githubusercontent.com
fi

if [ "$REMOTE_VERSION" = "null" ]; then
  SPEC_URL=$(readlink -f $CN_VERSION/specification)
  echo " ---> using local version of specification: $SPEC_URL"
  echo " ---> copy specification to: $SPEC_DIR/haproxy_spec.yaml"
  cp $SPEC_URL/build/haproxy_spec.yaml $SPEC_DIR/haproxy_spec.yaml
  echo " ---> copy copyright to :    $SPEC_DIR/copyright.txt"
  cp $SPEC_URL/copyright.txt $SPEC_DIR/copyright.txt
else
  echo " ---> URL path: $URL_PATH"
  echo " ---> repository path: $REPO_PATH"
  echo " ---> client native version: $CN_VERSION"
  SPEC_URL=https://$URL_PATH/$REPO_PATH/$EXTRA_PATH$CN_VERSION/specification

  echo " ---> fetching specification: $SPEC_URL/build/haproxy_spec.yaml"
  wget -q -O $SPEC_DIR/haproxy_spec.yaml $SPEC_URL/build/haproxy_spec.yaml
  echo " ---> fetching copyright: $SPEC_URL/copyright.txt"
  wget -q -O $SPEC_DIR/copyright.txt $SPEC_URL/copyright.txt
fi

echo "module github.com/haproxytech" > $DST_DIR/go.mod
mkdir -p $DST_DIR/dataplaneapi/operations
cp configure_data_plane.go $DST_DIR/dataplaneapi/configure_data_plane.go

swagger generate server -f $SPEC_DIR/haproxy_spec.yaml \
    -A "Data Plane" \
    -t $DST_DIR \
    --existing-models github.com/haproxytech/client-native/v6/models \
    --exclude-main \
    --skip-models \
    -s dataplaneapi \
    --additional-initialism=FCGI \
    --additional-initialism=QUIC \
    --tags=Discovery \
    --tags=ServiceDiscovery \
    --tags=Information \
    --tags=Health \
    --tags=Specification \
    --tags=SpecificationOpenapiv3 \
    --tags=Transactions \
    --tags=Sites \
    --tags=Stats \
    --tags=Global \
    --tags=Frontend \
    --tags=Backend \
    --tags=Ring \
    --tags=DgramBind \
    --tags=LogForward \
    --tags=Bind \
    --tags=Server \
    --tags=Configuration \
    --tags=HTTPCheck \
    --tags=HTTPRequestRule \
    --tags=HTTPResponseRule \
    --tags=HTTPAfterResponseRule \
    --tags=HTTPErrorRule \
    --tags=HTTPErrors \
    --tags=BackendSwitchingRule \
    --tags=ServerSwitchingRule \
    --tags=TCPResponseRule \
    --tags=TCPRequestRule \
    --tags=TCPCheck \
    --tags=FCGIApp \
    --tags=Filter \
    --tags=StickRule \
    --tags=LogTarget \
    --tags=Reloads \
    --tags=ACL \
    --tags=Defaults \
    --tags=StickTable \
    --tags=Maps \
    --tags=Nameserver \
    --tags=Cluster \
    --tags=Peer \
    --tags=PeerEntry \
    --tags=ProcessManager \
    --tags=Resolver \
    --tags=Spoe \
    --tags=SpoeTransactions \
    --tags=Storage \
    --tags="ACL Runtime" \
    --tags=ServerTemplate \
    --tags=Cache \
    --tags=DeclareCapture \
    --tags=Userlist \
    --tags=User \
    --tags=Group \
    --tags=MailerEntry \
    --tags=Mailers \
    --tags=Table \
    --tags=CrtStore \
    --tags=CrtLoad \
    --tags=QUICInitialRule \
    --tags=Traces \
    --tags=LogProfile \
    -r $SPEC_DIR/copyright.txt \
    --template-dir generate/swagger/templates

echo " ---> removing doc.go"
rm doc.go || echo "doc.go does not exists"
echo " ---> removing embedded_spec.go"
rm embedded_spec.go  ||  echo "embedded_spec.go does not exists"
echo " ---> removing server.go"
rm server.go ||  echo "server.go does not exists"
echo " ---> removing operations/*"
rm -rf operations/* ||  echo "operations/ does not exists"

echo " ---> copy generated files to destination"
cp -a $DST_DIR/dataplaneapi/. .
