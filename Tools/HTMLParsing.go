package Tools

import "strings"

func parseHTML(HTML []byte) []byte {
	// Get string from bytes
	HTMLString := string(HTML)

	// Replace paths
	HTMLString = strings.Replace(HTMLString, "develop.drewberry.co.uk", "localhost:8040", -1)
	HTMLString = strings.Replace(HTMLString, "drewberryinsurance.co.uk", "localhost:8040", -1)

	return []byte(HTMLString)
}