# JWT Automatic Token Renewal Implementation

## Overview

This implementation provides automatic JWT token renewal using a refresh token pattern, ensuring users remain authenticated without manual re-login for up to 7 days.

## Architecture

### Two-Token System

1. **Access Token (JWT)**
   - Duration: 15 minutes
   - Stored in: HttpOnly cookie (`access_token`)
   - Purpose: Authenticates API requests
   - Auto-renewed: When < 5 minutes remain

2. **Refresh Token**
   - Duration: 7 days
   - Stored in: 
     - HttpOnly cookie (`refresh_token`)
     - Database (hashed with SHA-256)
   - Purpose: Renew expired access tokens
   - Security: Single-use, rotated on each refresh

## Security Features

### 1. HttpOnly Cookies
Both tokens are stored in HttpOnly cookies, preventing XSS attacks from accessing them via JavaScript.

### 2. Token Hashing
Refresh tokens are hashed (SHA-256) before database storage, protecting against database compromise.

### 3. Token Rotation
Each time a refresh token is used, it's revoked and a new pair of tokens is generated, preventing replay attacks.

### 4. Automatic Renewal
The JWT guard middleware automatically renews access tokens when they have < 5 minutes remaining, providing seamless user experience.

### 5. Revocation on Logout
All refresh tokens for a user are revoked on logout, ensuring clean session termination.

## API Endpoints

### 1. Sign In
```
POST /api/v1/auth/signin
```
- Generates both access and refresh tokens
- Sets both cookies

### 2. Refresh Token
```
POST /api/v1/auth/refresh
```
- Uses refresh token cookie to generate new tokens
- Revokes old refresh token
- Sets new cookies

### 3. Logout
```
POST /api/v1/auth/logout
```
- Revokes all refresh tokens for the user
- Clears both cookies

## Database Schema

### refresh_token Table
```sql
CREATE TABLE refresh_token (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL REFERENCES account(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP NULL
);
```

## Configuration

Token expiration times are defined in `back/commons/constants/token_expiration.constant.go`:

```go
const ACCESS_TOKEN_EXPIRATION = int((15 * time.Minute) / time.Second)  // 900 seconds
const REFRESH_TOKEN_EXPIRATION = int((168 * time.Hour) / time.Second)  // 7 days
```

## Usage Flow

### Initial Login
1. User submits credentials via `/api/v1/auth/signin`
2. Backend validates credentials
3. Backend generates access token (15 min) and refresh token (7 days)
4. Backend stores hashed refresh token in database
5. Backend sets both tokens in HttpOnly cookies
6. User is authenticated

### Automatic Renewal (Transparent to Frontend)
1. User makes API request with access token
2. JWT guard middleware checks token expiration
3. If token expires in < 5 minutes:
   - Generate new access token
   - Set new access token cookie
4. Request proceeds normally

### Manual Refresh (If Access Token Expired)
1. Frontend detects 401 error
2. Frontend calls `/api/v1/auth/refresh`
3. Backend validates refresh token from cookie
4. Backend generates new token pair
5. Backend revokes old refresh token
6. Backend sets new cookies
7. Frontend retries original request

### Logout
1. User initiates logout
2. Frontend calls `/api/v1/auth/logout`
3. Backend revokes all refresh tokens for user
4. Backend clears cookies
5. User is logged out

## Code Examples

### Generating Tokens (Signin Service)
```go
func (s *SigninService) GenerateTokens(claims *guard.Claims) (tokenResponse TokenResponseDto, err error) {
    // Generate access token (15 min)
    accessToken, err := s.GenerateAccessToken(claims)
    if err != nil {
        return tokenResponse, err
    }

    // Generate refresh token (7 days)
    refreshToken, err := s.refreshTokenRepository.Create(
        claims.Id,
        time.Now().Add(168*time.Hour),
    )
    if err != nil {
        return tokenResponse, err
    }

    tokenResponse.AccessToken = accessToken
    tokenResponse.RefreshToken = refreshToken

    return tokenResponse, nil
}
```

### Automatic Renewal (JWT Guard)
```go
// Auto-renew token if it's close to expiration (less than 5 minutes)
if ShouldRenewToken(claims) {
    newToken, err := GenerateAccessToken(claims)
    if err == nil {
        lib.SetAccessTokenCookie(c, newToken, 0)
    }
    // Continue even if renewal fails - current token is still valid
}
```

## Testing

Run tests:
```bash
cd back
go test ./pkg/signin/... ./commons/guard/... -v
```

Tests cover:
- Refresh token generation
- Token hashing (deterministic and collision-free)
- Auto-renewal logic (timing thresholds)

## Migration

The refresh token table is automatically created on application startup via GORM auto-migration.

## Cleanup

Consider implementing a periodic cleanup job to remove expired and revoked refresh tokens:

```go
func (*RefreshTokenRepository) DeleteExpired() error {
    err := db.GetDB().
        Where("expires_at < ? OR is_revoked = ?", time.Now(), true).
        Delete(&model.RefreshToken{}).Error
    return err
}
```

## Best Practices

1. **Always use HTTPS** in production to protect cookies in transit
2. **Set appropriate cookie domain** in production environment
3. **Monitor refresh token usage** for suspicious patterns
4. **Implement rate limiting** on refresh endpoint
5. **Log token refresh events** for audit trail

## Troubleshooting

### Access Token Not Renewing
- Check if `ShouldRenewToken` threshold is appropriate
- Verify private key is accessible
- Check logs for token generation errors

### Refresh Token Invalid
- Ensure refresh token hasn't expired (7 days)
- Check if token was revoked (logout, password change)
- Verify token hash matches database entry

### Database Issues
- Ensure `refresh_token` table exists
- Check account foreign key constraint
- Verify database connection pool settings

## Security Considerations

1. **Never expose refresh tokens** in API responses (body)
2. **Always use HttpOnly cookies** for token storage
3. **Rotate refresh tokens** on each use
4. **Revoke all tokens** on password change
5. **Monitor for token abuse** (multiple simultaneous sessions)
6. **Use secure random** for token generation (crypto/rand)

## Future Enhancements

- [ ] Implement device/session tracking
- [ ] Add "remember me" option (30-day refresh tokens)
- [ ] Implement token family tracking for better security
- [ ] Add refresh token usage analytics
- [ ] Implement sliding window refresh token expiration
