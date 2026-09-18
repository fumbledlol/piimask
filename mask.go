package piimask

import (
	"encoding/hex"
	"net/netip"
	"strconv"
	"strings"
	"unicode/utf8"
)

const defaultBullet = "•"

// MaskEmail masks the local part of an email address while preserving domain and casing.
// Safely handles UTF-8 runes without multibyte corruption.
func MaskEmail(email string) string {
	atIdx := strings.LastIndexByte(email, '@')
	if atIdx <= 0 || atIdx == len(email)-1 {
		return strings.Repeat(defaultBullet, 8)
	}

	local := email[:atIdx]
	domain := email[atIdx+1:]
	runeCount := utf8.RuneCountInString(local)

	var b strings.Builder
	b.Grow(len(email) + 12)

	switch {
	case runeCount <= 1:
		b.WriteString(defaultBullet)
	case runeCount <= 4:
		r, size := utf8.DecodeRuneInString(local)
		b.WriteRune(r)
		writeBullets(&b, runeCount-1)
		_ = size
	case runeCount <= 8:
		r1, size1 := utf8.DecodeRuneInString(local)
		b.WriteRune(r1)
		writeBullets(&b, runeCount-2)
		rLast, _ := utf8.DecodeLastRuneInString(local[size1:])
		b.WriteRune(rLast)
	default:
		r1, s1 := utf8.DecodeRuneInString(local)
		r2, _ := utf8.DecodeRuneInString(local[s1:])
		b.WriteRune(r1)
		b.WriteRune(r2)
		writeBullets(&b, runeCount-4)
		rLast, sLast := utf8.DecodeLastRuneInString(local)
		rPrev, _ := utf8.DecodeLastRuneInString(local[:len(local)-sLast])
		b.WriteRune(rPrev)
		b.WriteRune(rLast)
	}

	b.WriteByte('@')
	b.WriteString(domain)
	return b.String()
}

// MaskIP masks the middle segments of IPv4 and IPv6 addresses.
func MaskIP(ipStr string) string {
	addr, err := netip.ParseAddr(strings.TrimSpace(ipStr))
	if err != nil {
		return strings.Repeat(defaultBullet, 8)
	}

	if addr.Is4() {
		b := addr.As4()
		var buf [24]byte
		dst := strconv.AppendUint(buf[:0], uint64(b[0]), 10)
		dst = append(dst, ".••.••."...) // ".••."
		dst = strconv.AppendUint(dst, uint64(b[3]), 10)
		return string(dst)
	}

	b := addr.As16()
	var hexStart [4]byte
	var hexEnd [4]byte
	hex.Encode(hexStart[:], b[0:2])
	hex.Encode(hexEnd[:], b[14:16])

	var out strings.Builder
	out.Grow(24)
	out.Write(hexStart[:])
	out.WriteString(":••••:••••:")
	out.Write(hexEnd[:])
	return out.String()
}

// MaskToken masks session tokens or API keys, showing the first and last 4 characters if length > 8.
func MaskToken(token string) string {
	token = strings.TrimSpace(token)
	runeCount := utf8.RuneCountInString(token)
	if runeCount <= 4 {
		return strings.Repeat(defaultBullet, runeCount)
	}

	var b strings.Builder
	b.Grow(len(token) + 16)

	if runeCount <= 8 {
		// keep first 2 runes
		written := 0
		rem := token
		for written < 2 && len(rem) > 0 {
			r, sz := utf8.DecodeRuneInString(rem)
			b.WriteRune(r)
			rem = rem[sz:]
			written++
		}
		writeBullets(&b, runeCount-2)
		return b.String()
	}

	// keep first 4 runes
	rem := token
	for i := 0; i < 4 && len(rem) > 0; i++ {
		r, sz := utf8.DecodeRuneInString(rem)
		b.WriteRune(r)
		rem = rem[sz:]
	}

	writeBullets(&b, runeCount-8)

	// extract last 4 runes
	tail := token
	for utf8.RuneCountInString(tail) > 4 {
		_, sz := utf8.DecodeRuneInString(tail)
		tail = tail[sz:]
	}
	b.WriteString(tail)

	return b.String()
}

// MaskDeviceID is an alias for MaskToken for device identifier strings.
func MaskDeviceID(token string) string {
	return MaskToken(token)
}

// MaskCreditCard masks card numbers, preserving only the last 4 digits.
func MaskCreditCard(cardNumber string) string {
	var digits [32]byte
	count := 0
	for i := 0; i < len(cardNumber) && count < len(digits); i++ {
		c := cardNumber[i]
		if c >= '0' && c <= '9' {
			digits[count] = c
			count++
		}
	}
	if count < 4 {
		return "•••• •••• •••• ••••"
	}

	var b strings.Builder
	b.Grow(20)
	b.WriteString("•••• •••• •••• ")
	b.Write(digits[count-4 : count])
	return b.String()
}

// MaskPhone masks phone numbers while keeping the country code and last digits.
func MaskPhone(phone string) string {
	cleaned := strings.TrimSpace(phone)
	if len(cleaned) < 5 {
		return strings.Repeat(defaultBullet, len(cleaned))
	}

	hasPlus := strings.HasPrefix(cleaned, "+")
	var digits [32]byte
	count := 0
	for i := 0; i < len(cleaned) && count < len(digits); i++ {
		c := cleaned[i]
		if c >= '0' && c <= '9' {
			digits[count] = c
			count++
		}
	}

	if count < 4 {
		return strings.Repeat(defaultBullet, 8)
	}

	visibleEnd := 4
	if count < 7 {
		visibleEnd = 2
	}

	countryDigits := 1
	if count > 11 {
		countryDigits = 2
	}

	maskedLength := count - countryDigits - visibleEnd
	if maskedLength < 2 {
		maskedLength = 2
	}

	var b strings.Builder
	b.Grow(32)
	if hasPlus {
		b.WriteByte('+')
	}
	b.Write(digits[:countryDigits])
	writeBullets(&b, maskedLength)
	b.Write(digits[count-visibleEnd : count])
	return b.String()
}

func writeBullets(b *strings.Builder, n int) {
	for i := 0; i < n; i++ {
		b.WriteString(defaultBullet)
	}
}