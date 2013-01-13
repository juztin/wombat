#!/bin/sh -

go build ../wombat \
	../wombat/backends \
	../wombat/backends/mongo \
	../wombat/config \
	../wombat/imgconv \
	../wombat/models/user \
	../wombat/template \
	../wombat/template/data \
	../wombat/apps/article \
	../wombat/apps/article/backends \
	../wombat/apps/article/backends/mongo \
	../wombat/apps/article/handlers \
