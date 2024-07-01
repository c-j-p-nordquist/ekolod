package tlsutils

import (
	"crypto/tls"
	"time"
)

func VersionToString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "Unknown"
	}
}

func CalculateCertExpiryDays(notAfter time.Time) int {
	return int(time.Until(notAfter).Hours() / 24)
}
