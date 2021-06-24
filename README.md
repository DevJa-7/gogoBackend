`export SWAGGERAPIPATH=/Volumes/D/Sources/gogo/webservice`
`swagger -apiPackage="api/v1" -mainApiFile="api/route.go" -format="asciidoc"`
`asciidoctor -a toc2 -a stylesheet=golo.css -a stylesdir=./stylesheets API.adoc`