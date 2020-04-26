package tools

import "time"

// Return current time in UTC
var Now = func() time.Time { return time.Now().UTC() }
