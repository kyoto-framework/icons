package icons

import "html/template"

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"icon": Icon,
	}
}
