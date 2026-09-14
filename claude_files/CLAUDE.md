# CLAUDE.md

Project instructions for Claude Code when working in this repository.

## What this project is

A single-file Go CLI (`main.go`) used as a beginner training lab. It
**simulates** the early steps of an RKE2 deployment — nothing here touches
the network or installs real software. Flags: `--init`/`-i`,
`--configure`/`-cfg`, `--check`/`-chk`, `--deploy`/`-d`.

## Build and test

- Build: `go build -o setup main.go`
- Run: `./setup --help`
- Format: `gofmt -l .` (fix with `gofmt -w main.go`)
- Vet: `go vet ./...`
- There are no automated tests yet. If you add Go source files beyond
  `main.go`, add table-driven tests alongside them (`_test.go`) and run
  with `go test ./...`.

## Conventions

- **Single file, standard library only.** Do not add third-party
  dependencies (e.g. a YAML library) or introduce a `go.mod` requiring
  `go get` — the point of this project is that a student can `go build`
  it with zero setup. Keep `servers.yaml` parsing/writing hand-rolled.
- **Flags come in long/short pairs.** Every flag is registered twice with
  `flag.BoolVar`, both pointing at the same variable, e.g. `--check` and
  `-chk`. Keep new flags consistent with this pattern (see
  `.claude/skills/add-cli-flag/SKILL.md`).
- **Flag execution order is fixed**, not argument order: init → configure →
  check → deploy, regardless of the order flags were typed. Preserve this
  in `main()`.
- **Output is colorized** using the ANSI constants already defined
  (`colorGreen`, `colorRed`, `colorYellow`, `colorCyan`, `colorBold`,
  `colorReset`). Reuse these rather than adding new raw escape codes.
- **`hosts.ini` and `servers.yaml` must stay consistent** — whenever one is
  regenerated, the other should be too, so the IPs match.
- Every user-facing action prints a short `==> Doing thing` header, then
  indented `- created`/`- wrote`/`[done]` lines, matching the existing
  style in `runInit`, `runConfigure`, `runCheck`, and `runDeploy`.

## What NOT to do

- Don't make this tool do anything real (no actual SSH, no actual
  network calls, no actual package installs). It's a training simulation.
- Don't commit generated `hosts.ini`, `setup` (the compiled binary), or
  `servers.yaml` — these are meant to be produced by `--init`/`--configure`
  each time a student runs the lab.

## Where things are

- `main.go` — the entire CLI tool
- `README.md` — student-facing hands-on lab instructions (build/run/test)
- `.claude/skills/add-cli-flag/SKILL.md` — the pattern for adding a new
  flag to this tool
