package HTTP

func HTMLParsing(HTML []byte) []byte {
	// Get string from bytes
	HTMLString := string(HTML)

	// @todo: filters/processors for the HTML (empty space removal for example)

	return []byte(HTMLString)
}
