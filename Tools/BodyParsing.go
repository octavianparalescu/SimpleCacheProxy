package Tools

import "path/filepath"

var ignoredExtensions = map[string]bool{
	"jpg": true,
	"jpeg": true,
	"bmp": true,
	"png": true,
}

func GetProperBODY(path string, Body []byte) []byte {

	var extension = filepath.Ext(path)

	// todo: do a better check to see if this is a text file such as a HTML or CSS
	if (!ignoredExtensions[extension]) {
		// if we have an HTML file, we will parse the body
		Body = parseHTML(Body)
	}

	return Body
}