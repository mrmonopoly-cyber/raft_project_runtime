#!/bin/sh
protoc -I=. --go_out=../out $1
