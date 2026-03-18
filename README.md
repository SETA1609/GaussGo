# GaussGo

GaussGo is a local-first modular CLI/TUI learning framework for math topics (starting with linear algebra), built in Go with Bubble Tea.

## Current Status
- Phase 0 (parallel-readiness contracts) is in place.
- Phase 1 foundation is implemented (registry, mod discovery/validation/resolution, bootstrap service wiring).
- Phase 1 hardening pass applied from bug report (structured app errors, concurrency safety, and wider test coverage).
- Phase 2 state core is implemented (state model/store, normalization, active pointer, atomic persistence).
- Planning includes modular learning flows, per-mod/concept statistics, quiz modes, and step-by-step feedback.
- Logging and event bus services are defined at the contract level.

## Key Features (Planned + In Progress)
- Drop-in mods discovered from `mods/<mod-id>`.
- State persistence per learner in `states/`.
- Locale-aware UX with state-backed language selection.
- Learning hub with:
  - statistics (mod -> concept table)
  - lessons
  - quiz modes (`Random`, `Selected Concepts`)
  - helpers
- Quiz review with results, step-by-step solutions, and clarifications for wrong answers.

## Project Structure
- `phases/` implementation roadmap by phase.
- `docs/contracts/` versioned architecture and integration contracts.
- `internal/contracts/` compile-time interfaces used across phases.
- `mods/` modular content and manifests.

## Development
- Run: `make run`
- Lint: `make lint`
- Test: `make test`
- Full check: `make check`

## Run Options
- Local:
  - `make run`
  - or `go run ./cmd/gaussgo`
- Container:
  - Build: `docker build -t gaussgo .`
  - Run interactive: `docker run -it --rm gaussgo`

## Container Data and Mod Persistence
- The image includes seeded `data/`, `mods/`, `states/`, and `pdfs/`.
- For persistent state across runs, mount `states/` from host:
  - `docker run -it --rm -v "$(pwd)/states:/root/states" gaussgo`
- For local mod development with live host files, mount `mods/` too:
  - `docker run -it --rm -v "$(pwd)/mods:/root/mods" -v "$(pwd)/states:/root/states" gaussgo`

## Planning and Docs Policy
- After each phase implementation (or when explicitly requested), update `README.md` and relevant docs.
