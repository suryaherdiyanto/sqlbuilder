package pkg

import (
	"fmt"
	"strings"
)

func ColumnSplitter(s, leftQuote, rightQuote string) string {
	if strings.Contains(s, ".") {
		parts := strings.SplitN(s, ".", 2)
		return fmt.Sprintf("%s%s%s.%s%s%s", leftQuote, parts[0], rightQuote, leftQuote, parts[1], rightQuote)
	}

	if s == "*" || strings.Contains(strings.ToLower(s), " as ") {
		return s
	}

	return fmt.Sprintf("%s%s%s", leftQuote, s, rightQuote)
}
