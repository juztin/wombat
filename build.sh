#!/bin/sh -

go build ../wombat \
	../wombat/backends \
	../wombat/config \
	../wombat/template \
	../wombat/template/data \
	../wombat/users \
	../wombat/users/backends \
	../wombat/users/backends/mongo \
	../wombat/wraps
