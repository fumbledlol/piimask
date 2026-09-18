package piimask

import (
	"fmt"
	"net"
	"strings"
)

const defaultBullet = "•"

// MaskEmail masks the local part of an email address preserving readability and domain.
// Examples:
//   "a@domain.com"          -> "•@domain.com"
//   "al@domain.com"         -> "a•@domain.com"
//   "alex@domain.com"       -> "a•••@domain.com"
//   "alice@domain.com"      -> "a•••e@domain.com"
//   "alexander@domain.com"  -> "al•••••er@domain.com"
// Safely handles Unicode / non-ASCII characters without corrupting multi-byte runes.
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return strings.Repeat(defaultBullet, 8)
	}

	local := parts[0]
	domain := parts[1]
	runes := []rune(local)
	length := len(runes)

	switch {
	case length <= 1:
		return defaultBullet + "@" + domain
	case length <= 4:
		return string(runes[:1]) + strings.Repeat(defaultBullet, length-1) + "@" + domain
	case length <= 8:
		return string(runes[:1]) + strings.Repeat(defaultBullet, length-2) + string(runes[length-1:]) + "@" + domain
	default:
		return string(runes[:2]) + strings.Repeat(defaultBullet, length-4) + string(runes[length-2:]) + "@" + domain
	}
}

// MaskIP masks the middle octets/hextets of IPv4 and IPv6 addresses.
// Examples:
//   "192.168.1.50"          -> "192.••.••.50"
//   "2001:0db8:85a3:0000:0000:8a2e:0370:7334" -> "2001:••••:••••:7334"
func MaskIP(ipStr string) string {
	ip := net.ParseIP(strings.TrimSpace(ipStr))
	if ip == nil {
		return strings.Repeat(defaultBullet, 8)
	}

	if ipv4 := ip.To4(); ipv4 != nil {
		return fmt.Sprintf("%d.••.••.%d", ipv4[0], ipv4[3])
	}

	// IPv6: format first 2 bytes and last 2 bytes with leading zeros preserved
	firstSegment := fmt.Sprintf("%02x%02x", ip[0], ip[1])
	lastSegment := fmt.Sprintf("%02x%02x", ip[14], ip[15])
	return fmt.Sprintf("%s:••••:••••:%s", firstSegment, lastSegment)
}

// MaskToken masks sensitive API keys, session tokens, or device IDs.
// Keeps the first 4 and last 4 characters visible if length > 8.
// Example: "dev_tok_9f823a789bc123" -> "dev_••••••••••••c123"
func MaskToken(token string) string {
	token = strings.TrimSpace(token)
	runes := []rune(token)
	n := len(runes)

	if n <= 4 {
		return strings.Repeat(defaultBullet, n)
	}
	if n <= 8 {
		return string(runes[:2]) + strings.Repeat(defaultBullet, n-2)
	}
	return string(runes[:4]) + strings.Repeat(defaultBullet, n-8) + string(runes[n-4:])
}

// MaskDeviceID is an alias for MaskToken for device identifier strings.
func MaskDeviceID(token string) string {
	return MaskToken(token)
}

// MaskCreditCard masks credit card numbers, showing only the last 4 digits.
// Example: "4111222233334444" -> "•••• •••• •••• 4444"
func MaskCreditCard(cardNumber string) string {
	var digits strings.Builder
	for _, r := range cardNumber {
		if r >= '0' && r <= '9' {
			digits.WriteRune(r)
		}
	}
	str := digits.String()
	n := len(str)
	if n < 4 {
		return "•••• •••• •••• ••••"
	}
	last4 := str[n-4:]
	return fmt.Sprintf("•••• •••• •••• %s", last4)
}

// MaskPhone masks phone numbers while keeping the country code and last 2-4 digits.
// Example: "+1 (555) 123-4567" -> "+1••••••4567"
func MaskPhone(phone string) string {
	cleaned := strings.TrimSpace(phone)
	if len(cleaned) < 5 {
		return strings.Repeat(defaultBullet, len(cleaned))
	}
	hasPlus := strings.HasPrefix(cleaned, "+")
	digitsOnly := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, cleaned)

	n := len(digitsOnly)
	if n < 4 {
		return strings.Repeat(defaultBullet, 8)
	}

	visibleEnd := 4
	if n < 7 {
		visibleEnd = 2
	}

	prefix := ""
	if hasPlus {
		prefix = "+"
	}

	// NANP (+1) uses 1-digit country code; longer numbers use 2
	countryDigits := 1
	if n > 11 {
		countryDigits = 2
	}

	maskedLength := n - countryDigits - visibleEnd
	if maskedLength < 2 {
		maskedLength = 2
	}

	return fmt.Sprintf("%s%s%s%s", prefix, digitsOnly[:countryDigits], strings.Repeat(defaultBullet, maskedLength), digitsOnly[n-visibleEnd:])
}