# piimask

[![Go Reference](https://pkg.go.dev/badge/github.com/fumbledlol/piimask.svg)](https://pkg.go.dev/github.com/fumbledlol/piimask)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Format-preserving PII redaction and token masking for structured logs and user interfaces in Go.

## Features

- Standard library only (`net/netip`, `strings`, `unicode/utf8`).
- Decodes runes directly to avoid breaking multi-byte UTF-8 sequences in international addresses (`björn@example.com` becomes `b•••n@example.com`).
- Preserves leading zeros in IPv6 addresses (`fe80:••••:••••:0001`).
- Runs in 60 ns per call with 1 allocation using `net/netip` and `strings.Builder`.
- Masks emails, IPv4/IPv6 addresses, tokens, API keys, device IDs, payment cards, and phone numbers.

## Installation

```bash
go get github.com/fumbledlol/piimask
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/fumbledlol/piimask"
)

func main() {
	// Email addresses
	fmt.Println(piimask.MaskEmail("alexander@example.com"))
	// al•••••er@example.com

	// IPv4 and IPv6
	fmt.Println(piimask.MaskIP("192.168.1.100"))
	// 192.••.••.100
	fmt.Println(piimask.MaskIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"))
	// 2001:••••:••••:7334

	// Tokens and device IDs
	fmt.Println(piimask.MaskToken("dev_tok_9f823a789bc123"))
	// dev_••••••••••••c123

	// Payment cards
	fmt.Println(piimask.MaskCreditCard("4111 2222 3333 4444"))
	// •••• •••• •••• 4444

	// Phone numbers
	fmt.Println(piimask.MaskPhone("+1 (555) 123-4567"))
	// +1••••••4567
}
```

### Structured logging with slog

```go
slog.Info("user login",
    slog.String("email", piimask.MaskEmail(user.Email)),
    slog.String("ip", piimask.MaskIP(clientIP)),
)
```

## Benchmarks

```
BenchmarkMaskEmail-12     21070707        61.46 ns/op        48 B/op       1 allocs/op
BenchmarkMaskIP-12        18037436        61.49 ns/op        24 B/op       1 allocs/op
BenchmarkMaskToken-12      4674932       245.10 ns/op       144 B/op       2 allocs/op
```

## License

[MIT](LICENSE)