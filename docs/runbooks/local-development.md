# Local development

From the repository root run make doctor, then make check. Start with make dev. The server binds only to 127.0.0.1:8080. Use PORT=8081 make dev for a second worktree. No .env loader is installed.

Try curl -i 'http://127.0.0.1:8080/api/v1/estimate?points=3' and the invalid case points=0. Structured lifecycle logs appear on stderr. No request payloads are logged. Stop with Ctrl-C; active requests receive up to five seconds to finish.

For an occupied port, select another PORT. For formatting failures run make fmt and recheck. For race-test toolchain errors install the platform C compiler through your normal environment setup; race detection needs cgo. No external service should be needed for tests.

No persistent state exists, so restarting resets the entire sample. The default check builds packages without leaving a server binary in the repository.
