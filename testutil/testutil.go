package testutil

import (
	"html/template"
	"io"
)

var (
	simpleHTMLTemplate = template.Must(template.New("td").Parse(`
<!DOCTYPE html>
<html>
  <head>
  </head>
  <body>
    {{range .Links}}
      <a href='{{.}}'> 
    {{end}}
    
    {{range .Emails}}
      <a href='mailto:{{.}}'> 
    {{end}}
  </body>
</html>
`))
)

type SimpleTemplateData struct {
	Links  []string
	Emails []string
}

func RenderSimpleTemplate(w io.Writer, data *SimpleTemplateData) {
	simpleHTMLTemplate.Execute(w, data)
}
