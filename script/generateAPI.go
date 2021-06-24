package script

import (
	"log"
	"os"
)

// GenerateAPI generates api with format
func GenerateAPI(format string) []string {
	commands := make([]string, 4)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	commands = append(commands, "export SWAGGERAPIPATH="+dir)

	prefix := "swagger -apiPackage=\"api/v1\" -mainApiFile=\"api/route.go\" -format=\"" + format + "\" "
	switch format {
	case "go":
		fallthrough
	case "swagger":
		commands = append(commands, prefix+"-output=\"document/\"")
	case "asciidoc":
		commands = append(commands, prefix+"-output=\"document/API.adoc\"")
		commands = append(commands, "asciidoctor -a icons -a toc2 -a stylesheet=golo.css -a stylesdir=./stylesheets document/API.adoc")
	case "markdown":
		commands = append(commands, prefix+"-output=\"document/API.md\"")
	case "confluence":
		commands = append(commands, prefix+"-output=\"document/API.confluence\"")
	}
	println(commands)
	return commands
}
