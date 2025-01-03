package main

import (
	"flag"
	"fmt"
	"strings"
)

const (
	URIScheme      = "gemini"
	URIDefaultPort = 1965
	URILengthLimit = 1024
	CRLF           = "\r\n"
)

func getAbsoluteURI(authority string) string {
	authority = strings.TrimSuffix(authority, "/")
	return fmt.Sprintf("%s://%s:%d/%s", URIScheme, authority, URIDefaultPort, CRLF)
}

func main() {
	uriAuthority := flag.String("path", "", "sets the path to send the request")
	flag.Parse()

	fmt.Print(getAbsoluteURI(*uriAuthority))
}
