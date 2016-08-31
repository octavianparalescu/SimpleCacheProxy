package Tools

import "strings"

func parseHTML(HTML []byte) []byte {
	// Get string from bytes
	HTMLString := string(HTML)

	// Replace paths
	HTMLString = strings.Replace(HTMLString, "develop.drewberry.co.uk", "develop.drewberry.co.uk", -1)
	HTMLString = strings.Replace(HTMLString, "drewberryinsurance.co.uk", "develop.drewberry.co.uk", -1)

	return []byte(HTMLString)
}