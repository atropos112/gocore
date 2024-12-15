package logging

import (
	"fmt"
	"strings"
)

func formatMsgWithArgs(msg string, args ...any) string {
	// Create a slice to hold formatted key-value pairs
	var formattedArgs []string

	// Loop through the args slice in pairs
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			formattedArgs = append(formattedArgs, fmt.Sprintf("%s: %v", args[i], args[i+1]))
		}
	}

	// Join the formatted key-value pairs with a comma and space
	return msg + "|" + strings.Join(formattedArgs, ", ")
}
