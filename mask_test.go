package piimask

import (
	"strings"
	"testing"
)

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a@example.com", "•@example.com"},
		{"al@example.com", "a•@example.com"},
		{"bob@example.com", "b••@example.com"},
		{"alex@example.com", "a•••@example.com"},
		{"alice@example.com", "a•••e@example.com"},
		{"alexander@example.com", "al•••••er@example.com"},
		{"not-an-email", "••••••••"},
		// International Unicode handling (safe multi-byte runes)
		{"björn@example.com", "b•••n@example.com"},
		{"user@münchen.de", "u•••@münchen.de"},
	}

	for _, tc := range tests {
		actual := MaskEmail(tc.input)
		if actual != tc.expected {
			t.Errorf("MaskEmail(%q) = %q, expected %q", tc.input, actual, tc.expected)
		}
	}
}

func TestMaskIP(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"192.168.1.1", "192.••.••.1"},
		{"10.0.0.254", "10.••.••.254"},
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", "2001:••••:••••:7334"},
		{"fe80::1", "fe80:••••:••••:0001"}, // verifies leading zeros in 0001
		{"invalid-ip", "••••••••"},
	}

	for _, tc := range tests {
		actual := MaskIP(tc.input)
		if actual != tc.expected {
			t.Errorf("MaskIP(%q) = %q, expected %q", tc.input, actual, tc.expected)
		}
	}
}

func TestMaskToken(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"abc", "•••"},
		{"abcdef", "ab••••"},
		{"1234567890", "1234••7890"},
		{"dev_token_123456789abcdef", "dev_" + strings.Repeat("•", 17) + "cdef"},
	}

	for _, tc := range tests {
		actual := MaskToken(tc.input)
		if actual != tc.expected {
			t.Errorf("MaskToken(%q) = %q, expected %q", tc.input, actual, tc.expected)
		}
	}
}

func TestMaskCreditCard(t *testing.T) {
	card := "4111 2222 3333 4444"
	expected := "•••• •••• •••• 4444"
	if actual := MaskCreditCard(card); actual != expected {
		t.Errorf("MaskCreditCard(%q) = %q, expected %q", card, actual, expected)
	}
}

func TestMaskPhone(t *testing.T) {
	phone := "+1 (555) 123-4567"
	expected := "+1••••••4567"
	if actual := MaskPhone(phone); actual != expected {
		t.Errorf("MaskPhone(%q) = %q, expected %q", phone, actual, expected)
	}
}