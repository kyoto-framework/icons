package icons

import (
	"fmt"
	"strings"
)

func Icon(set, name string, classes ...string) string {
	// Get icon
	icon := IconMap[set][name]
	// Add classes
	if len(classes) > 0 {
		icon = strings.ReplaceAll(icon, "<svg", fmt.Sprintf(`<svg class="%s"`, strings.Join(classes, " ")))
	}
	// Return icon
	return icon
}
