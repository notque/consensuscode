# Build & Test Status: March 2026

| Project | Build | Tests | Issues |
|---------|-------|-------|--------|
| CollectiveFlow CLI | PASS | PASS (no test files) | Fixed `fmt.Println` redundant newline in config.go |
| Bluesky Tool | PASS | PASS (28/28) | None |
| CollectiveFlow Web | PASS | PASS (124/124) | Fixed corrupted UTF-8 in about.html, type guard in urgency_color filter |
| Collective Website | PASS | PASS (22/22) | Fixed app factory pattern, added missing routes, updated agent count |
| Root Makefile | PASS | N/A | `make help` works correctly |

## Issues Found

### CollectiveFlow CLI (Go)
- **`internal/cli/config.go:48`**: `fmt.Println` called with a string ending in `\n`, producing a redundant newline. The Go vet tool flags this as an error, preventing the test suite from running.

### CollectiveFlow Web (Flask)
- **`templates/about.html`**: Five emoji icon spans contained corrupted non-UTF-8 bytes (0xab, 0x94, 0x0f, 0x0b, etc.), causing `UnicodeDecodeError` when Jinja2 loaded the template. This broke the about page and navigation tests.
- **`app.py` urgency_color filter**: Passing a non-hashable type (e.g., a list) to `dict.get()` raised `TypeError`. Filter lacked a type guard.

### Collective Website (Flask)
- **`app.py` application factory**: `create_app()` returned a bare Flask app without registering any routes. Routes were decorated on the module-level `app` instance only, so `create_app('testing')` in the test suite produced an app with zero routes -- all 22 tests got 404s.
- **Missing routes**: `/projects`, `/how-we-work`, `/contribute` had templates but no route handlers in `app.py`.
- **`about` route**: Did not pass `agent_count` to the template, causing a rendering error.
- **Agent count hardcoded to 5**: The collective now has 7 agents. Updated to 7 in the index route, about route, and consensus API endpoint.

## Fixes Applied

### CollectiveFlow CLI
- `projects/collectiveflow/internal/cli/config.go`: Removed trailing `\n` from `fmt.Println` argument.

### CollectiveFlow Web
- `projects/collectiveflow/web/templates/about.html`: Replaced 5 corrupted emoji byte sequences with valid HTML entities.
- `projects/collectiveflow/web/app.py`: Added `isinstance(urgency, str)` type guard in `urgency_color` filter.

### Collective Website
- `projects/collective-website/app.py`: Restructured into a proper application factory pattern -- all routes, helper functions, and error handlers now register inside `create_app()` so that `create_app('testing')` returns a fully functional app.
- Added missing route handlers for `/projects`, `/how-we-work`, `/contribute`.
- Added `agent_count=7` parameter to the `about` route.
- Updated collective size from 5 to 7 across all endpoints.

## Environment Notes

- Python projects require a virtual environment (PEP 668 blocks system-level pip installs on this machine).
- pyenv provides Python 3.14.3; system python is 3.9. Always use `python3` explicitly.
- CollectiveFlow web pytest.ini includes `--cov` flags that require `pytest-cov` to be installed.
