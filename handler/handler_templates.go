package handler

import "os"
import "text/template"



func WriteStringToFileWithTemplates(  text string , partName string, file *os.File, aStruct interface{}) {
	tmpl := template.New( partName )
	template.Must(tmpl.Parse(text))
	err := tmpl.Execute(file, aStruct)
	if ( err != nil ) {
		panic(err)
	}

}




