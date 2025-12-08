# JWT Token Renewal - Implementation Summary

## Problem Statement

The current JWT implementation had the following limitations:
- JWT tokens expired after 7 days
- Users had to re-authenticate after token expiration
- No mechanism to renew tokens automatically

## Solution

Implemented a **secure refresh token pattern** with automatic token renewal:

### Architecture Changes

1. **Two-Token System**
   - Short-lived access token (15 minutes) for API authentication
   - Long-lived refresh token (7 days) stored in database and cookie
   - Automatic renewal when access token has < 5 minutes remaining

2. **Security Enhancements**
   - HttpOnly cookies prevent XSS attacks
   - SHA-256 hashing of refresh tokens in database
   - Single-use tokens with rotation on each refresh
   - Complete revocation on logout

## Files Changed

### New Files
- `back/db/models/refresh_token.model.go` - Refresh token database model
- `back/db/repository/refresh_token.repo.go` - Refresh token CRUD operations
- `back/pkg/signin/signin_service_test.go` - Unit tests for token generation
- `back/commons/guard/jwt_guard_test.go` - Unit tests for auto-renewal
- `back/docs/JWT_TOKEN_RENEWAL.md` - Comprehensive documentation

### Modified Files

**Database & Models:**
- `back/db/migration.go` - Added refresh_token table to auto-migration

**Services:**
- `back/pkg/signin/signin.service.go` - Added `GenerateTokens()` and `RefreshAccessToken()`
- `back/pkg/account/account.service.go` - Updated to use dual tokens
- `back/pkg/provider/provider.service.go` - Updated to use dual tokens

**Controllers:**
- `back/pkg/signin/signin.controller.go` - Added `/auth/refresh` endpoint
- `back/pkg/account/account.controller.go` - Sets both cookies
- `back/pkg/provider/provider.controller.go` - Sets both cookies
- `back/pkg/auth/auth.controller.go` - Revokes refresh tokens on logout

**DTOs:**
- `back/pkg/signin/signin.dto.go` - Added RefreshToken to response
- `back/pkg/account/account.dto.go` - Added AccountTokensDto

**Infrastructure:**
- `back/commons/guard/jwt_guard.go` - Auto-renewal logic
- `back/commons/lib/cookie.go` - Refresh token cookie helper
- `back/commons/constants/token_expiration.constant.go` - New constants
- `back/server/router.go` - Added refresh endpoint route

## API Changes

### New Endpoint
```
POST /api/v1/auth/refresh
```
Refreshes the access token using the refresh token cookie.

### Modified Endpoints
All authentication endpoints now set both access and refresh token cookies:
- `POST /api/v1/auth/signin`
- `POST /api/v1/account` (create account)
- `PATCH /api/v1/account` (update account with username)
- `GET /api/v1/auth/:provider/callback` (OAuth callbacks)

### Logout Enhancement
```
POST /api/v1/auth/logout
```
Now revokes all refresh tokens for the user before clearing cookies.

## Database Schema

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

## Configuration Changes

New constants in `token_expiration.constant.go`:
```go
const ACCESS_TOKEN_EXPIRATION = 900    // 15 minutes
const REFRESH_TOKEN_EXPIRATION = 604800 // 7 days
```

## Testing

All tests pass:
- ✅ Refresh token generation test
- ✅ Token hashing test (deterministic)
- ✅ Auto-renewal logic test
- ✅ Existing event service tests (unchanged)
- ✅ Existing encryption tests (unchanged)

## Migration

No manual migration needed. GORM AutoMigrate will create the refresh_token table on application startup.

## Backward Compatibility

✅ **Fully backward compatible**
- Existing JWT signing/verification unchanged
- Same RSA keypair used
- Existing endpoints continue to work
- Frontend changes are optional

## Security Considerations

### What We Did Right
1. **HttpOnly cookies** - Prevents JavaScript access (XSS protection)
2. **Token hashing** - SHA-256 hash stored in DB (compromise protection)
3. **Token rotation** - Single-use tokens prevent replay attacks
4. **Automatic renewal** - Seamless UX without security compromise
5. **Complete revocation** - Logout invalidates all user tokens

### Production Checklist
- [ ] Ensure HTTPS is enabled (required for secure cookies)
- [ ] Verify cookie domain is set correctly
- [ ] Consider implementing refresh token cleanup job
- [ ] Monitor refresh token usage for suspicious patterns
- [ ] Implement rate limiting on refresh endpoint

## Frontend Impact

**No changes required!** The implementation is transparent:

1. **Automatic renewal**: Happens in middleware, invisible to frontend
2. **Cookies**: Automatically sent with requests
3. **Logout**: Works exactly as before

**Optional enhancement**: Frontend can handle expired tokens:
```javascript
if (response.status === 401) {
  await fetch('/api/v1/auth/refresh', { method: 'POST' });
  // Retry original request
}
```

## Performance Impact

**Minimal:**
- Database query on signin/refresh (1 insert)
- Database query on logout (1 update)
- In-memory token generation (no DB hit on auto-renewal)
- Cookie overhead: ~100 bytes per request

## Monitoring Recommendations

1. Track refresh token creation rate
2. Monitor refresh endpoint usage
3. Alert on high token revocation rates
4. Log token refresh failures

## References

- [RFC 6749 - OAuth 2.0 (Refresh Tokens)](https://tools.ietf.org/html/rfc6749#section-1.5)
- [OWASP Session Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)

## Next Steps

1. ✅ Merge this PR
2. ✅ Deploy to development environment
3. ✅ Test with frontend
4. ✅ Monitor token refresh patterns
5. ✅ Consider implementing device/session tracking
6. ✅ Add analytics for token usage

## Questions & Answers

**Q: Why 15 minutes for access token?**
A: Short-lived tokens reduce the window of opportunity if compromised. Auto-renewal makes this transparent.

**Q: Why store refresh token in DB?**
A: Enables revocation (logout, password change) and prevents token reuse.

**Q: Why hash refresh tokens?**
A: Protects users if database is compromised - attacker can't use stolen tokens.

**Q: Why rotate tokens?**
A: Prevents replay attacks and limits damage from token theft.

**Q: Is this OWASP compliant?**
A: Yes, follows OWASP session management best practices.
