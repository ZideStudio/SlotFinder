package constants

import "time"

// Access token expires in 15 minutes
const ACCESS_TOKEN_EXPIRATION = int((15 * time.Minute) / time.Second)

// Refresh token expires in 7 days
const REFRESH_TOKEN_EXPIRATION = int((168 * time.Hour) / time.Second)

// Token renewal threshold - renew when less than 5 minutes remain
const TOKEN_RENEWAL_THRESHOLD = int((5 * time.Minute) / time.Second)

// Legacy constant for backward compatibility
const TOKEN_EXPIRATION = ACCESS_TOKEN_EXPIRATION

// Duration constants for use with time package
const ACCESS_TOKEN_DURATION = 15 * time.Minute
const REFRESH_TOKEN_DURATION = 168 * time.Hour
const TOKEN_RENEWAL_THRESHOLD_DURATION = 5 * time.Minute
