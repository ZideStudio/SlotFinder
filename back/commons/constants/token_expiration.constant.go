package constants

import "time"

const ACCESS_TOKEN_EXPIRATION = 15 * time.Minute
const REFRESH_TOKEN_EXPIRATION = 8760 * time.Hour // 1 year
const TOKEN_RENEWAL_THRESHOLD_EXPIRATION = 5 * time.Minute
