package utils

import "strings"

func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email // fallback if it's malformed
	}

	local := parts[0]
	domain := parts[1]

	maskedLocal := string(local[0]) + "***"
	maskedDomain := string(domain[0]) + "***.com"

	return maskedLocal + "@" + maskedDomain
}
