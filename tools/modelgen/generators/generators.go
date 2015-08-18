package generators

import (
	"text/template"

	"github.com/contiv/objmodel/tools/modelgen/texthelpers"
)

var templateMap = map[string]*template.Template{}

var funcMap = template.FuncMap{
	"initialCap": texthelpers.InitialCap,
	"initialLow": texthelpers.InitialLow,
	"depunct":    texthelpers.Depunct,
	"capFirst":   texthelpers.CapFirst,
}

func ParseTemplates() error {
	for name, content := range templates {
		var err error
		templateMap[name], err = template.New(name).Funcs(funcMap).Parse(content)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetTemplate(templateName string) *template.Template {
	return templateMap[templateName]
}
