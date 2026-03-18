# POC: GaussGo Modular Learning CLI

## Goal
Build a local-first CLI learning app in Go using Bubble Tea and Lipgloss, with Gonum for math engines.
The app is modular: new learning modules are dropped into `/mods/<mod-id>` and discovered at runtime.

## Scope (v1)
- Focus area: linear algebra content via mod `linearAlgebra`
- Runtime: local binary and container
- Data-driven: lessons/exercises in JSON/YAML inside each mod
- Persistent learning progress via `state.json`

## Core Product Rules
- Every mod is a directory in `/mods`
- Every mod must include `manifest.json`
- A state tracks user progress and enabled/disabled mods
- User can create a new state or load an existing one from main menu
- A state records which mods (and versions) existed when it was created

## Directory Layout
```text
gaussgo/
├── cmd/gaussgo/main.go
├── internal/
│   ├── app/                    # startup flow and menu controller
│   ├── tui/                    # Bubble Tea screens
│   ├── mods/                   # manifest parsing, dependency resolution, hooks API
│   ├── content/                # lesson/exercise loaders from mod data folders
│   ├── exercises/              # exercise evaluators
│   ├── math/                   # gonum wrappers
│   └── state/                  # state create/load/save and schema
├── mods/
│   ├── core/
│   │   ├── manifest.json
│   │   └── public/             # exported hook/function contracts docs for mod authors
│   └── linearAlgebra/
│       ├── manifest.json
│       └── data/
│           ├── units/
│           └── concepts/
├── states/                     # saved state files (state-*.json)
├── Dockerfile
└── POC.md
```

## Mod System

### Manifest (`/mods/<id>/manifest.json`)
Required fields:
- `id`: stable mod id (directory name should match)
- `name`: display name
- `version`: semver string
- `entry`: optional data entrypoint path
- `dependencies`: list of `{ "id": "...", "version": ">=x.y.z" }`
- `provides`: optional hook names exposed by this mod

Example:
```json
{
  "id": "linearAlgebra",
  "name": "Linear Algebra",
  "version": "0.1.0",
  "entry": "data/units",
  "dependencies": [
    { "id": "core", "version": ">=0.1.0" }
  ],
  "provides": ["learning.units", "exercise.compute"]
}
```

### Mod Loading Rules
- Scan `/mods/*/manifest.json`
- Validate manifest schema
- Build dependency graph and reject unresolved dependencies
- Resolve load order (topological: `core` before dependents)
- Expose enabled mod list to state and app menus

### Core Mod
`/mods/core` defines reusable hooks and public functions other mods call.

Initial hooks contract (v1):
- `learning.units.list` -> returns unit metadata
- `learning.unit.get` -> returns full unit content by id
- `exercise.evaluate` -> evaluates answer and returns feedback

## State System

### State Lifecycle
- Main menu offers:
  - Start Learning (choose enabled mod)
  - Create New State
  - Load Existing State
  - Select Language (English/Spanish for v1)
  - Mod Settings
  - Exit

After selecting a mod in Start Learning, show Learning Hub options:
- Statistics
- Lessons
- Random Quiz
- Helpers
- Back to Main Menu

### State Storage
- State files live in `/states`
- Naming: `state-<slug>.json`
- Active state can be tracked in `/states/active.txt` (optional)

### `state.json` Shape
```json
{
  "schemaVersion": 1,
  "stateId": "default",
  "profileName": "Sebastian",
  "createdAt": "2026-03-17T00:00:00Z",
  "updatedAt": "2026-03-17T00:00:00Z",
  "modsAtCreation": [
    { "id": "core", "version": "0.1.0" },
    { "id": "linearAlgebra", "version": "0.1.0" }
  ],
  "enabledMods": ["core", "linearAlgebra"],
  "ui": {
    "locale": "en"
  },
  "progress": {
    "linearAlgebra": {
      "readConcepts": ["dot-product"],
      "exerciseStats": {
        "attempted": 12,
        "correct": 9,
        "correctnessPercent": 75
      },
      "completedUnits": ["vectors"],
      "conceptScores": {
        "dot-product": 0.9
      }
    }
  }
}
```

Path convention:
- Active state file is selected from `/states`
- Default pointer file: `/states/active.txt`
- New state creation writes `/states/state-<slug>.json`

### State Rules
- `core` cannot be disabled
- Mod enable/disable in Mod Settings writes to current state
- Selected locale is saved in state (`ui.locale`), default `en`
- Random quiz concepts are limited to concepts already read once (`progress.<modId>.readConcepts`)
- If a previously enabled mod is missing, state loads with warning and mod marked unavailable
- If mod version differs from `modsAtCreation`, app shows compatibility notice

## Menus and Flows

### Startup Flow
1. Discover mods and validate dependency graph
2. Show Main Menu
3. User picks create/load state
4. Bind selected state as active runtime context
5. User starts learning from one enabled mod

### Mod Settings Flow
- Show all discovered mods with status:
  - enabled
  - disabled
  - missing dependency
  - incompatible version
- Toggle allowed mods and save to active state

### Learning Hub Flow (after Start Learning)
- Statistics:
  - Grouped by mod
  - Expand mod row to view concept-level table (for example matrices, vectors)
  - Per concept: attempted, correct, correctness percent
- Lessons:
  - Show concepts from selected mod (for example vectors/matrices from linear algebra)
- Random Quiz:
  - Two setup options: `Random` and `Selected Concepts`
  - `Selected Concepts` uses multi-select from read-only concept set (already read once)
  - Both setup options generate questions randomly from their eligible concept pool
  - Answers can be submitted by selecting an option or entering the option letter
  - End-of-quiz review includes final results and step-by-step solution per question
  - Wrong answers include clarification of mistake and correct method
- Helpers:
  - Call operation helper functions contributed by mods
- Back to Main Menu

Navigation rule:
- Every menu includes a return action that moves exactly one step back in the menu tree.

Formatting rule:
- Learning and quiz math content should be authored/rendered with LaTeX-style math when possible, with readable plain-text fallback for non-LaTeX renderers.

## Initial Implementation Plan
1. Implement `internal/mods` manifest structs + loader + dependency resolver
2. Implement `internal/state` create/load/save + schema validation
3. Implement Bubble Tea main menu and mod settings screens
4. Add `mods/core` and `mods/linearAlgebra` manifests + seed learning content
5. Wire `linearAlgebra` content loader and one Gonum-backed exercise evaluator

## Runtime Commands
- Local: `go run ./cmd/gaussgo`
- Container: `docker build -t gaussgo . && docker run -it gaussgo`

## Notes
- PDFs are reference material only while authoring the first mod.
- Runtime learning content is loaded from mod data files, not directly from PDFs.
