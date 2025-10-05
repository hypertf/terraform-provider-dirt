package provider

import (
    "strings"
)

// isNotFound returns true when the underlying client reported a missing resource.
// We intentionally avoid changing the client API; it currently returns errors with
// messages like "<resource> not found" on HTTP 404.
func isNotFound(err error) bool {
    if err == nil {
        return false
    }
    // Unwrap in case of wrapped errors
    // Note: using errors.Is would need a sentinel in client; we keep it simple here.
    msg := err.Error()
    return strings.Contains(strings.ToLower(msg), "not found")
}

