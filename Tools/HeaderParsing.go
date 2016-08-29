package Tools

import (
	"net/http"
)

var allowedHeaders = map[string]bool{
	"Content-Type": true,
	"Expires": true,
	"Cache-Control": true,
}

func GetProperHeaders(h http.Header) map[string]string {
	var newHeaders map[string]string

	newHeaders = make(map[string]string)

	for k, v := range h {
		// todo: better header control, especially regarding the expires/cache-control
		// todo: use a black-list instead of a white-list
		if (allowedHeaders[k]) {
			newHeaders[k] = v[0]
		}
	}

	return newHeaders
}
