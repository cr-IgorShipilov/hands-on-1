# Part 2 — Developing This Go App With Claude Code

In Part 1 you built and ran `setup` by hand. In this part you'll configure
**Claude Code** (Anthropic's terminal coding agent) to work on this same
project, using project memory (`CLAUDE.md`), project settings
(`.claude/settings.json`), and a custom skill
(`.claude/skills/add-cli-flag/`). By the end you'll use Claude to add a new
flag to the tool yourself.

This lab assumes the `setup-lab/` folder from Part 1, containing `main.go`
and `README.md`.

---

## Step 1 — Install Claude Code

Check if it's already installed:

```bash
claude --version
```

If not, install it (macOS/Linux/WSL):

```bash
curl -fsSL https://claude.ai/install.sh | bash
```

Windows PowerShell:

```powershell
irm https://claude.ai/install.ps1 | iex
```

Confirm:

```bash
claude --version
```

You'll need a Claude account (Pro, Max, Team, or Enterprise) or Claude
Console API access to log in.

---

## Step 2 — Start a session in the project

```bash
cd ~/rke2-lab   # or wherever setup-lab lives
claude
```

The first time you run `claude`, it walks you through login in your
browser. Once you're in, you'll see a prompt inside your terminal.

Ask a quick question to confirm Claude can see your project:

```text
what does this project do?
```

Claude reads the files in the working directory as needed — you don't
need to paste code in.

---

## Step 3 — Add the project files

From claude_files copy `CLAUDE.md` and the `.claude/` directory (containing
`settings.json` and `skills/add-cli-flag/SKILL.md`) into your
`setup-lab/` folder, alongside `main.go`:

```text
setup-lab/
├── main.go
├── README.md
├── CLAUDE.md
└── .claude/
    ├── settings.json
    └── skills/
        └── add-cli-flag/
            └── SKILL.md
```

Restart or resume your session so Claude Code picks the new files up,
then confirm `CLAUDE.md` loaded:

```text
/context
```

Look for `CLAUDE.md` under **Memory files** in the output.

---

## Step 4 — Read what CLAUDE.md does

Open `CLAUDE.md` and read it. It tells Claude, every session:

- how to build, format, and vet this project
- the project's conventions (single file, standard library only,
  long/short flag pairs, fixed execution order, colorized output)
- what *not* to do (don't make it call real services; don't commit
  generated files)

**Try it yourself:** ask Claude a question that CLAUDE.md should
influence, and see if the answer reflects it:

```text
what's the right way to add a dependency to this project?
```

You should get pushback, since CLAUDE.md says standard-library-only.

---

## Step 5 — Check the project's settings

Open `.claude/settings.json`. It pre-approves the Bash commands this
project needs (`go build`, `go run`, `go vet`, `go test`, `gofmt`,
`./setup`) so Claude doesn't stop and ask permission for routine
commands, and it denies reading `.env` files as a safety habit.

Run `/status` in your session and confirm `Project settings` (or
similar) is listed as a loaded settings source.

**Try it yourself:** ask Claude to build the project:

```text
build the project and run --help
```

Because `Bash(go build *)` and `Bash(./setup *)` are pre-approved, this
should run without a permission prompt. Compare this to asking for a
command that *isn't* in the allow list, like `git push`, which should
still prompt you.

---

## Step 6 — Look at the custom skill

Open `.claude/skills/add-cli-flag/SKILL.md`. This is a **skill**: a
reusable, step-by-step procedure Claude can load automatically when it's
relevant, or that you can invoke directly by typing `/add-cli-flag`.

It encodes exactly how flags are added in this codebase: register long
and short forms, add the flag to the fixed execution order, follow the
existing `run*()` function style, update `printUsage()`, and update
`README.md`.

List your available skills:

```text
/skills
```

You should see `add-cli-flag` in the list.

---

## Step 7 — Use the skill to add a new flag

Ask Claude to add a new flag. Don't specify implementation details —
that's what the skill is for:

```text
add a --status/-s flag that prints a fake summary of the cluster's
current health (all 3 nodes, uptime, and version)
```

Watch how Claude approaches the task. It should:

1. Load the `add-cli-flag` skill (automatically, based on its
   description, or you can invoke it directly with `/add-cli-flag`
   first)
2. Register `--status`/`-s` the same way the other flags are registered
3. Add a `runStatus()` function matching the existing style and colors
4. Update `printUsage()` and `README.md`
5. Build and smoke-test it

When it's done, verify yourself:

```bash
go build -o setup main.go
./setup --status
```

---

## Step 8 — Compare with and without the skill

Temporarily rename the skill file so Claude can't use it:

```bash
mv .claude/skills/add-cli-flag/SKILL.md .claude/skills/add-cli-flag/SKILL.md.bak
```

Start a fresh session (`/clear` or a new `claude` session) and ask for
another flag, e.g. `--reset` to delete `hosts.ini` and `servers.yaml`.
Compare the result to Step 7 — is the flag registered the same way? Is
the execution order preserved? Restore the skill file afterward:

```bash
mv .claude/skills/add-cli-flag/SKILL.md.bak .claude/skills/add-cli-flag/SKILL.md
```

---

## Step 9 — Clean up

```bash
rm -f hosts.ini servers.yaml setup
```

`CLAUDE.md` and `.claude/` are meant to stay in the project — they're
part of what you'd commit for teammates (or future-you) to reuse.

---

## Wrap-up questions

1. What's the difference between putting an instruction in `CLAUDE.md`
   versus a skill in `.claude/skills/`? When would you use each?
2. Why does `.claude/settings.json` pre-approve specific `go` subcommands
   instead of just allowing all `Bash` commands?
3. In Step 8, what did Claude get "wrong" without the skill that the
   skill fixed? What does that tell you about when a skill is worth
   writing?
4. `CLAUDE.md` says "don't add third-party dependencies." What happens
   if you ask Claude to use a YAML library anyway — does it push back,
   comply, or ask you to confirm?
5. This project's settings file is `.claude/settings.json` (shared,
   commit it). What would you put in `.claude/settings.local.json`
   instead, and why would you keep that one out of version control?
