#!/usr/bin/env bash

docker run -d -p 9090:9090 -v /var/run/docker.sock:/var/run/docker.sock ipedrazas/botd:latest
