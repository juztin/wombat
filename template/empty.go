package template

const Empty = "{{define \"title\"}}Template Doesn't Exist{{end}}\n" +
	"{{define \"content\"}}\n" +
	"<div style='margin: 50px auto;'>\n" +
	"This template doesn't exist, or hasn't been created yet.{{if .User.IsAdmin}} <a href='?edit'>Create Template</a>{{end}}" +
	"</div>\n" +
	"{{end}}\n" +
	"{{define \"scripts\"}}{{end}}\n"
