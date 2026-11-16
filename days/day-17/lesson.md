# Day 17 — Authentication, authorization, and credential safety

**Time: 55 minutes** — 5 min recall · 10 min concepts · 15 min walkthrough · 15 min exercises · 10 min review

**Version target:** Go 1.27.3; examples use features available in Go 1.26 where possible.

## Objectives
- Distinguish authentication (who) from authorization (may do what).
- Store password verifiers using a vetted adaptive password-hashing package.
- Understand token lifetime, validation, revocation, and TLS requirements.
- Test an authorization boundary without presenting a toy as production security.

## Why this matters
CareLog's fictional patient and observation records are sensitive by design, even in a teaching example. Correctly identifying a caller is not the same as deciding whether they may access a particular patient's data. Security must be part of the API design, not an afterthought.

## Concepts and explanation
Never store or log plaintext passwords. Use a vetted password-hashing implementation such as Argon2id or bcrypt from a maintained security package with documented parameters, salts, and migration policy; do not invent a hash scheme. The CareLog reference uses `golang.org/x/crypto/bcrypt` and `github.com/golang-jwt/jwt/v5` to demonstrate a configured login and expiring signed tokens. Signing keys and opaque token values require cryptographic randomness. Tokens need narrow scope, expiry, TLS protection, and a revocation/key-rotation plan. Validate signature/issuer/audience/expiry; an unguessable ID is not authorization. Rate-limit credential attempts, avoid account enumeration, and keep secrets out of logs. The reference auth flow is educational, single-account sample code—not production-ready identity management.

## Runnable example
Run `go run ./days/day-17/examples` from the repository root. It uses `httptest` and a fixed demonstration identity to show ownership authorization. It has no password handling, token issuance, persistence, TLS, or production-ready auth.

## Project Change

### What We Are Adding
A teaching login flow using a vetted adaptive password hash and expiring signed tokens, plus authentication middleware and patient-resource authorization.

### Why We Need It
CareLog patient/observation routes must reject anonymous requests and enforce access policy. Authentication alone does not grant access to every patient's data.

### Files Added
- `project/carelog/internal/auth/auth_test.go` — issuance, expiry, and verification cases.
- `project/carelog/internal/auth/middleware_test.go` — protected-route test matrix.
- `project/carelog/cmd/hash-password/main.go` — local password-hash helper.

### Files Changed
- `project/carelog/internal/auth/auth.go` and `middleware.go` — token validation and principal context.
- `project/carelog/internal/httpapi/patient.go` and `observation.go` — route protection and auth checks.
- `project/carelog/internal/config/config.go` — secret/hash configuration validation.
- `project/carelog/go.mod` — vetted bcrypt/JWT dependencies.

### Go Concepts Used
Adaptive password hashing, cryptographic keys, signed claims, middleware, context values with private keys, expiry validation, and error mapping.

### Tests Added
Anonymous, valid, malformed, expired, wrong-role, and invalid-signature cases; allowed/denied patient ownership policy; ensure secrets never appear in error responses or logs.

### Architecture Impact
Identity is established at the authentication boundary and authorization is checked per resource. This single-account educational flow is not production-ready; prefer a vetted identity provider and security review before real credentials or sensitive data.

## Exercises
1. Add a patient-ownership policy to a protected observation read and ensure a request can access only patients permitted for the authenticated principal.
2. Document a password-hashing decision (vetted library, parameter policy, rehash-on-login plan); do not implement a custom hash.
3. Write an authorization test matrix for anonymous, owner, other user, expired credential, and revoked credential cases.

## Quiz
1. What is the distinction between authentication and authorization?
2. Why is a fast general-purpose hash such as SHA-256 unsuitable for password storage?
3. Does a valid signed token prove the caller can access every patient observation?

## Understanding check
Can you show where identity enters the request context, how patient access is checked, and how credentials are stored and revoked without logging their value?

## Common mistakes
- Storing plaintext/reversible passwords or rolling a custom password hash.
- Treating a valid token as permission for every resource.
- Using predictable tokens, omitting expiry/revocation, or accepting tokens over insecure transport.
- Returning detailed login failures that enable account enumeration.

## Production and architecture perspective
Prefer a reputable identity provider or maintained, vetted library and follow current security guidance. Threat-model account recovery, key rotation, CSRF (for cookie-based sessions), XSS/token exposure, abuse, and multi-tenant isolation. Keep authentication dependencies isolated behind interfaces, but do not confuse interface isolation with security assurance.

## Summary and outcome
You can articulate the identity/permission boundary and the minimum design decisions required before CareLog handles real credentials. **Outcome:** the project gains an authorization model and a safe implementation plan—not a production auth system from a classroom demo.

## Instructor answers
See [the Day 17 instructor guide](../../instructor-guide/day-17/solutions.md) after attempting the exercises.
