package urls

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

// Shorten handles shortening a long (origial) URL. In production I would do much
// more research into the needs here (performance, should we encrypt instead of hash
// in case all our data was wiped, etc). Shorten returns the first 8 characters of
// an md5 hash of the longURL.
func Shorten(longURL string) string {
	h := md5.New()
	io.WriteString(h, longURL)
	return hex.EncodeToString(h.Sum(nil))[0:8]
}
