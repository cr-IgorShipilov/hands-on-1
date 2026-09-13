---
name: add-cli-flag
description: Add a new command-line flag to the setup lab tool (main.go), following the project's existing long/short flag pattern. Use when asked to add a new --flag/-x option, a new simulated step, or to extend setup's behavior.
---

# Adding a flag to `setup`

Follow the pattern already used by `--init`, `--configure`, `--check`, and
`--deploy` in `main.go`. Don't introduce a different flag style.

## Steps

1. **Declare the variable** in `main()` alongside the other `bool` flag
   variables.

2. **Register the flag twice** — long form and short form — both bound to
   the same variable:

   ```go
   flag.BoolVar(&myFlag, "myflag", false, "what this flag does")
   flag.BoolVar(&myFlag, "mf", false, "shorthand for --myflag")
   ```

3. **Add it to the "nothing was passed" check** near the top of `main()`
   so `setup` with no flags still prints usage and exits 1.

4. **Add an `if myFlag { runMyFlag() }` block** in `main()`, in the same
   fixed position relative to the other flags if execution order matters
   (init → configure → check → deploy today). Flags always run in this
   fixed order regardless of the order the user typed them — don't make
   the new flag's position depend on argument order.

5. **Write `runMyFlag()`** following the existing style:
   - Print a header: `fmt.Printf("%s==> Doing the thing%s\n", colorCyan, colorReset)`
   - Print each step indented two spaces, using `colorGreen`/`colorRed`/
     `colorYellow` for status, matching `runInit`/`runConfigure`/
     `runCheck`/`runDeploy`.
   - If the step depends on `servers.yaml`, read it with the existing
     `loadServersOrExit()` helper — don't write a second YAML parser.
   - End with a `Next step:` hint pointing at whatever flag logically
     follows, the same way `runCheck` points at `--deploy`.

6. **Update `printUsage()`** to list the new flag in the flags table.

7. **Update `README.md`** to add a numbered step for the new flag,
   matching the existing "Step N — ..." format, plus a wrap-up question
   if it's the kind of design choice worth discussing.

8. **Build and smoke-test**: `go build -o setup main.go && ./setup --myflag`.
   Also run `gofmt -l .` and `go vet ./...` before considering the change
   done — this project has no other tests.

## Keep it simulation-only

This tool never performs real installs, SSH, or network actions. A new
flag should print realistic-looking, colorized, delayed output — not
call out to any real service.
