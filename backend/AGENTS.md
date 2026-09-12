# Backend

Run all project checks from the repository root with `make check`.

Handlers own HTTP parsing and response mapping. Services own business rules and must not import HTTP or handlers. cmd/server owns process startup and shutdown. Architecture tests enforce the internal layer allowlist.

Use the standard library in this sample. Wrap returned errors with context; use log/slog at process boundaries. Never log request bodies or secrets. HTTP errors use the documented JSON envelope. Inject time into handlers for deterministic tests.
