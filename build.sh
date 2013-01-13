#!/bin/sh -

go build ../wombat \
	../wombat/backends \
	../wombat/backends/mongo \
	../wombat/config \
	../wombat/imgconv \
	../wombat/models/user \
	../wombat/template \
	../wombat/template/data \
