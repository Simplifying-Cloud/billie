package shared

import (
	"fmt"
	"strings"
	"time"
)

// FormatCurrency formats a float64 as USD currency with thousand separators
func FormatCurrency(amount float64) string {
	// Format with 2 decimal places
	str := fmt.Sprintf("%.2f", amount)

	// Split into integer and decimal parts
	parts := strings.Split(str, ".")
	intPart := parts[0]
	decPart := parts[1]

	// Add thousand separators
	var result []byte
	for i, c := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}

	return "$" + string(result) + "." + decPart
}

// FormatDate formats a time.Time as "Jan 2, 2006"
func FormatDate(t time.Time) string {
	return t.Format("Jan 2, 2006")
}

// GetInitials returns the initials from a name
func GetInitials(name string) string {
	if len(name) == 0 {
		return "?"
	}
	initials := ""
	words := []rune(name)
	initials += string(words[0])
	for i, r := range words {
		if r == ' ' && i+1 < len(words) {
			initials += string(words[i+1])
			break
		}
	}
	return initials
}
