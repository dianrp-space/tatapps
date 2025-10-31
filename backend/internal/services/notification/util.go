package notification

import "strings"

// NormalizeWhatsAppRecipients formats raw user input into the pipe-separated
// representation the backend expects (e.g. "6281|6282").
func NormalizeWhatsAppRecipients(input string) string {
	recipients := SplitWhatsAppRecipients(input)
	if len(recipients) == 0 {
		return ""
	}
	return strings.Join(recipients, "|")
}

// SplitWhatsAppRecipients returns the distinct list of sanitized phone numbers
// extracted from a pipe/comma/semicolon separated string.
func SplitWhatsAppRecipients(input string) []string {
	if strings.TrimSpace(input) == "" {
		return nil
	}

	raw := strings.FieldsFunc(input, func(r rune) bool {
		switch r {
		case '|', ',', ';', '\n', '\r':
			return true
		default:
			return false
		}
	})

	seen := make(map[string]struct{}, len(raw))
	result := make([]string, 0, len(raw))
	for _, part := range raw {
		trimmed := strings.TrimSpace(part)
		normalized := strings.ReplaceAll(trimmed, " ", "")
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}

	return result
}
