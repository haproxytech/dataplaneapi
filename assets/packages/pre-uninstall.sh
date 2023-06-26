#!/bin/sh

systemctl stop dataplaneapi || true
systemctl disable dataplaneapi || true
