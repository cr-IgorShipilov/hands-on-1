# RKE2 Lab Helper — Hands-On Exercise

This is a training exercise for building, testing, and running a simple Go
CLI tool. The tool **simulates** the early steps of deploying an RKE2
Kubernetes cluster — it does not touch the network or install real
software. It is safe to run anywhere.

By the end of this lab you will have:
- Built a Go binary from source
- Generated inventory files (`hosts.ini`, `servers.yaml`)
- Run simulated connectivity checks and a simulated deployment
- Practiced basic CLI flag usage

## Prerequisites

- A terminal (Linux, macOS, or WSL/Git Bash on Windows)
- No prior Go experience required

---

## Step 1 — Install Go

Check if Go is already installed:

```bash
go version
```

If you see something like `go version go1.22.2 linux/amd64`, skip to Step 2.

If Go is **not** installed:

- **Ubuntu/Debian**: `sudo apt-get update && sudo apt-get install -y golang-go`
- **macOS (Homebrew)**: `brew install go`
- **Other platforms**: download from https://go.dev/dl/ and follow the installer instructions

Verify again with `go version` before continuing.

---

## Step 2 — Get the project files

Create a folder for the lab and place `main.go` inside it:

```bash
mkdir -p ~/rke2-lab
cd ~/rke2-lab
# copy main.go into this folder
```

Confirm the file is there:

```bash
ls -la
```

You should see `main.go`.

---

## Step 3 — Build the tool

Compile the source into an executable called `setup`:

```bash
go build -o setup main.go
```

This produces a single binary named `setup` in the current folder. If the
command finishes with no output, the build succeeded.

Make sure it's executable and runnable:

```bash
chmod +x setup
./setup
```

Since no flags were given, you should see a **usage/help message** printed
to the screen. This confirms the build works.

---

## Step 4 — Explore the help output

Run:

```bash
./setup --help
```

Read through the available flags:

| Long form     | Short form | Purpose                                         |
|---------------|------------|--------------------------------------------------|
| `--init`      | `-i`       | create `hosts.ini` and `servers.yaml` templates   |
| `--configure` | `-cfg`     | fill those templates with fake data for 3 servers |
| `--check`     | `-chk`     | simulate a connectivity check                     |
| `--deploy`    | `-d`       | simulate an RKE2 deployment                       |

---

## Step 5 — Initialize the lab files

```bash
./setup --init
```

This creates two files in your current folder:

- `hosts.ini` — an Ansible-style inventory template
- `servers.yaml` — a server inventory template

Inspect them:

```bash
cat hosts.ini
cat servers.yaml
```

Both files currently contain `CHANGEME` placeholders — nothing usable yet.

**Try it yourself:** run `./setup --init` a second time. What happens? Why?
(Hint: look at the "skipped" message — the tool won't overwrite existing files.)

---

## Step 6 — Generate fake server data

```bash
./setup --configure
```

This overwrites `hosts.ini` and `servers.yaml` with 3 fake servers, each
with a fully-qualified domain name, IP address, gateway, and two DNS
servers.

Inspect the results again:

```bash
cat servers.yaml
cat hosts.ini
```

Notice that the IP addresses in `hosts.ini` match the ones generated in
`servers.yaml` — the tool keeps both files consistent.

---

## Step 7 — Run a connectivity check

```bash
./setup --check
```

Watch the simulated output: for each server, the tool "checks" ping, SSH
port 22, and DNS resolution, with a short delay between each step (this
mimics what a real health-check script would look and feel like).

**Try it yourself:** what happens if you run `./setup --check` *before*
`--init`/`--configure` (in an empty folder)? Try it in a scratch folder to
see the error handling.

---

## Step 8 — Run the simulated deployment

```bash
./setup --deploy
```

Watch the output. The first server in `servers.yaml` is treated as the
initial control-plane node; the remaining two "join" the cluster using the
first server's IP address. At the end, a `kubectl get nodes`-style table
is printed showing all servers as `Ready`.

---

## Step 9 — Combine flags in one run

The tool supports chaining flags together, in either order:

```bash
rm -f hosts.ini servers.yaml   # start clean
./setup -i -cfg -chk -d
```

This runs init → configure → check → deploy in a single command, in that
fixed order, regardless of how you typed the flags.

---

## Step 10 — Clean up

When you're done experimenting:

```bash
rm -f hosts.ini servers.yaml setup
```

---

## Wrap-up questions

1. What is the difference between a long flag (`--init`) and its short form (`-i`) in this tool?
2. Which file does `--check` and `--deploy` read from? What happens if that file is missing?
3. Why might a real-world version of this tool use a proper YAML library instead of a hand-written parser?
4. What would you need to change in `main.go` to support a 4th or 5th server?
5. `--configure` overwrites existing files without asking. Is that good or risky design? How would you improve it?
