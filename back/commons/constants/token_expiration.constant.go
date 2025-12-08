package constants

import "time"

// Access token expires in 15 minutes
const ACCESS_TOKEN_EXPIRATION = int((15 * time.Minute) / time.Second)

// Refresh token expires in 7 days
const REFRESH_TOKEN_EXPIRATION = int((168 * time.Hour) / time.Second)

// Legacy constant for backward compatibility
const TOKEN_EXPIRATION = ACCESS_TOKEN_EXPIRATION
