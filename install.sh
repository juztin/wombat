#!/bin/sh -

go install bitbucket.org/juztin/wombat
go install bitbucket.org/juztin/wombat/backends
go install bitbucket.org/juztin/wombat/backends/mongo
go install bitbucket.org/juztin/wombat/config
go install bitbucket.org/juztin/wombat/imgconv
go install bitbucket.org/juztin/wombat/models/user
go install bitbucket.org/juztin/wombat/template
go install bitbucket.org/juztin/wombat/template/data

go install bitbucket.org/juztin/wombat/apps/article
go install bitbucket.org/juztin/wombat/apps/article/handlers
go install bitbucket.org/juztin/wombat/apps/article/backends
