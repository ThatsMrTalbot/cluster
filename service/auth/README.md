# Auth Service

Simple auth service using prototoken.

## Notes:
- The client contains a method to validate tokens, meaning once the login is complete no state needs to be maintained.
- Once a token expires (or is about to) the refresh token can be used to generate a new token and refresh token.
- Refresh tokens are only validated by the service. The client does not know the key.