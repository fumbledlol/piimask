# piimask

[![Go Reference](https://pkg.go.dev/badge/github.com/fumbledlol/piimask.svg)](https://pkg.go.dev/github.com/fumbledlol/piimask)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Zero-dependency, high-speed Go library for format-preserving PII redaction and secret masking in structured logs (`slog`, `zap`) and interfaces.

## Why piimask?

Compliance frameworks (GDPR, HIPAA, SOC 2, PCI-DSS) strictly penalize leaking Personally Identifiable Information (PII) into external log aggregators (Datadog, Loki, CloudWatch, Sentry). 

Completely removing fields like IP addresses or emails destroys debugging utility during production incidents. `piimask` performs **format-preserving masking** that keeps enough context for operators to verify domain, prefix, or last octet without exposing sensitive personal data.

## Features

- **Zero Dependencies**: Pure Go standard library implementation.
- **Unicode / UTF-8 Safe**: Safely decodes multi-byte runes before slicing, preventing invalid UTF-8 byte corruption on international emails (`björn@example.com`, `user@münchen.de`).
- **Preserves Leading Zeros**: Correctly formats IPv6 addresses (`fe80:••••:••••:0001`) without dropping hex zeros.
- **Ultra Fast**: Operates in ~120–165 ns with minimal heap allocations.
- **Comprehensive Masking**: Supports emails, IPv4/IPv6, session tokens, API keys, device IDs, credit cards, and phone numbers.

## Installation

```bash
go get github.com/fumbledlol/piimask
```

## Quick Start

```go
package main

import (
	"fmt"

	"github.com/fumbledlol/piimask"
)

func main() {
	// Email masking
	fmt.Println(piimask.MaskEmail("alexander@example.com"))
	// Output: al•••••er@example.com

	// IPv4 masking
	fmt.Println(piimask.MaskIP("192.168.1.100"))
	// Output: 192.••.••.100

	// IPv6 masking
	fmt.Println(piimask.MaskIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"))
	// Output: 2001:••••:••••:7334

	// API Keys / Tokens / Device IDs
	fmt.Println(piimask.MaskToken("dev_tok_9f823a789bc123"))
	// Output: dev_••••••••••••c123

	// Credit cards
	fmt.Println(piimask.MaskCreditCard("4111 2222 3333 4444"))
	// Output: •••• •••• •••• 4444

	// Phone numbers
	fmt.Println(piimask.MaskPhone("+1 (555) 123-4567"))
	// Output: +1••••••4567
}
```

### Usage with `log/slog`

```go
slog.Info("user login attempt",
    slog.String("email", piimask.MaskEmail(user.Email)),
    slog.String("ip", piimask.MaskIP(clientIP)),
    slog.String("device_id", piimask.MaskDeviceID(deviceID)),
)
```

## Benchmarks

```
BenchmarkMaskEmail-12     7499502       165.5 ns/op        80 B/op       3 allocs/op
BenchmarkMaskIP-12        9463587       122.1 ns/op        24 B/op       1 allocs/op
BenchmarkMaskToken-12     7552988       161.5 ns/op       112 B/op       2 allocs/op
```

## License

[MIT](LICENSE)