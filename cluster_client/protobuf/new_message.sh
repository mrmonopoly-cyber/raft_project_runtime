#!/bin/sh
protoc -I=. --go_out=. $1
