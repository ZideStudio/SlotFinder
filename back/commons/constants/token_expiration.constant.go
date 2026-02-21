package constants

import "time"

const ACCESS_TOKEN_EXPIRATION = 15 * time.Minute
const REFRESH_TOKEN_EXPIRATION = 720 * time.Hour // 30 days
const TOKEN_RENEWAL_THRESHOLD_EXPIRATION = 5 * time.Minute
