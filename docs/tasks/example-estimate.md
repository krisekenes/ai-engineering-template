# Return an effort estimate

Status: implemented example

## Outcome
A caller can get an illustrative hour estimate for a point value.

## Acceptance criteria
- GET /api/v1/estimate?points=3 returns 200 with data.points=3 and data.hours=12.
- Successful responses include a UTC RFC3339 meta.timestamp.
- Missing, non-integer, duplicated, or out-of-range points return 400 with error.code=VALIDATION_ERROR.
- Only values 1 through 100 are valid.
- Non-GET requests return 405; unknown paths return 404, both with JSON errors.
- GET /healthz returns 200 and data.status=ok.

## Scope and verification
Implement handlers and a pure estimating service; no storage or UI. Run make check. Tests cover boundaries, malformed input, method/path behavior, and the successful envelope. Consult api/openapi.yaml.
