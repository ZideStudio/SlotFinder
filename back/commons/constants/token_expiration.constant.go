package constants

import "time"

const ACCESS_TOKEN_EXPIRATION = 10 * time.Second
const REFRESH_TOKEN_EXPIRATION = 168 * time.Hour
const TOKEN_RENEWAL_THRESHOLD_EXPIRATION = 5 * time.Minute
