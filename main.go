// setup - a training CLI tool for a beginner RKE2 lab.
//
// This tool simulates the early steps of a real deployment workflow:
//  1. --init       create empty hosts.ini + servers.yaml templates
//  2. --configure  fill those templates with fake (but realistic) data
//  3. --check      pretend to test connectivity to each server
//  4. --deploy     pretend to install RKE2 on each server
//
// Nothing here touches the network or installs real software - it's all
// simulated output with delays, meant to teach CLI flag handling, file I/O,
// and basic Go program structure.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	hostsFile   = "hosts.ini"
	serversFile = "servers.yaml"
)

// ANSI colors - purely cosmetic, makes the lab output easier to read.
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// Server represents one node that will (pretend to) join the RKE2 cluster.
type Server struct {
	FQDN    string
	IP      string
	Gateway string
	DNS1    string
	DNS2    string
}

func main() {
	var (
		initFlag      bool
		configureFlag bool
		checkFlag     bool
		deployFlag    bool
	)

	// Each option is registered twice (long + short) so both forms update
	// the same variable, e.g. --init and -i behave identically.
	flag.BoolVar(&initFlag, "init", false, "create hosts.ini and servers.yaml templates")
	flag.BoolVar(&initFlag, "i", false, "shorthand for --init")

	flag.BoolVar(&configureFlag, "configure", false, "fill hosts.ini and servers.yaml with fake data for 3 servers")
	flag.BoolVar(&configureFlag, "cfg", false, "shorthand for --configure")

	flag.BoolVar(&checkFlag, "check", false, "simulate a connectivity check against servers.yaml")
	flag.BoolVar(&checkFlag, "chk", false, "shorthand for --check")

	flag.BoolVar(&deployFlag, "deploy", false, "simulate deploying RKE2 on servers.yaml")
	flag.BoolVar(&deployFlag, "d", false, "shorthand for --deploy")

	flag.Usage = printUsage
	flag.Parse()

	if !initFlag && !configureFlag && !checkFlag && !deployFlag {
		printUsage()
		os.Exit(1)
	}

	// Flags run in a fixed, sensible order regardless of how they were
	// typed, so e.g. `setup -d -i -cfg` still inits then configures then deploys.
	if initFlag {
		runInit()
	}
	if configureFlag {
		runConfigure()
	}
	if checkFlag {
		runCheck()
	}
	if deployFlag {
		runDeploy()
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `%ssetup%s - RKE2 lab helper (training tool, no real actions performed)

Usage:
  setup [flags]

Flags:
  -i,   --init        create %s and %s templates
  -cfg, --configure    fill templates with fake data for 3 servers
  -chk, --check        simulate a connectivity check
  -d,   --deploy       simulate an RKE2 deployment

Flags can be combined, e.g.:
  setup --init --configure --check --deploy
`, colorBold, colorReset, hostsFile, serversFile)
}

// ---------------------------------------------------------------------
// --init
// ---------------------------------------------------------------------

func runInit() {
	fmt.Printf("%s==> Initializing lab files%s\n", colorCyan, colorReset)

	hostsTemplate := `; hosts.ini - Ansible-style inventory
; Run "setup --configure" to fill this in automatically,
; or edit it by hand with your own server details.

[rke2_servers]
; server1.example.local ansible_host=CHANGEME
; server2.example.local ansible_host=CHANGEME
; server3.example.local ansible_host=CHANGEME

[rke2_servers:vars]
ansible_user=root
`

	serversTemplate := `# servers.yaml - server inventory for the RKE2 lab
# Run "setup --configure" to fill this in automatically,
# or edit it by hand with your own server details.
servers:
  - fqdn: CHANGEME
    ip: CHANGEME
    gateway: CHANGEME
    dns1: CHANGEME
    dns2: CHANGEME
  - fqdn: CHANGEME
    ip: CHANGEME
    gateway: CHANGEME
    dns1: CHANGEME
    dns2: CHANGEME
  - fqdn: CHANGEME
    ip: CHANGEME
    gateway: CHANGEME
    dns1: CHANGEME
    dns2: CHANGEME
`

	writeFileIfAbsent(hostsFile, hostsTemplate)
	writeFileIfAbsent(serversFile, serversTemplate)

	fmt.Printf("%s\nNext step:%s run %ssetup --configure%s to fill these files with fake data.\n",
		colorGreen, colorReset, colorBold, colorReset)
}

func writeFileIfAbsent(path, content string) {
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("  %s- skipped%s %s (already exists)\n", colorYellow, colorReset, path)
		return
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Printf("  %s- failed%s to create %s: %v\n", colorRed, colorReset, path, err)
		os.Exit(1)
	}
	fmt.Printf("  %s- created%s %s\n", colorGreen, colorReset, path)
}

// ---------------------------------------------------------------------
// --configure
// ---------------------------------------------------------------------

func runConfigure() {
	fmt.Printf("%s==> Generating fake server data%s\n", colorCyan, colorReset)

	servers := fakeServers(3)

	writeServersYAML(serversFile, servers)
	fmt.Printf("  %s- wrote%s %s\n", colorGreen, colorReset, serversFile)

	writeHostsINI(hostsFile, servers)
	fmt.Printf("  %s- wrote%s %s\n", colorGreen, colorReset, hostsFile)

	fmt.Println()
	for _, s := range servers {
		fmt.Printf("  %s%-22s%s ip=%-15s gw=%-13s dns1=%-13s dns2=%s\n",
			colorBold, s.FQDN, colorReset, s.IP, s.Gateway, s.DNS1, s.DNS2)
	}

	fmt.Printf("%s\nNext step:%s run %ssetup --check%s to test connectivity.\n",
		colorGreen, colorReset, colorBold, colorReset)
}

// fakeServers generates n servers with believable, made-up networking details.
func fakeServers(n int) []Server {
	rand.Seed(time.Now().UnixNano())
	servers := make([]Server, 0, n)
	subnet := fmt.Sprintf("192.168.%d", 10+rand.Intn(20))
	for i := 1; i <= n; i++ {
		servers = append(servers, Server{
			FQDN:    fmt.Sprintf("rke2-node%d.lab.local", i),
			IP:      fmt.Sprintf("%s.%d", subnet, 100+i),
			Gateway: fmt.Sprintf("%s.1", subnet),
			DNS1:    "8.8.8.8",
			DNS2:    "8.8.4.4",
		})
	}
	return servers
}

func writeServersYAML(path string, servers []Server) {
	var b strings.Builder
	b.WriteString("# servers.yaml - server inventory for the RKE2 lab (auto-generated)\n")
	b.WriteString("servers:\n")
	for _, s := range servers {
		fmt.Fprintf(&b, "  - fqdn: %s\n", s.FQDN)
		fmt.Fprintf(&b, "    ip: %s\n", s.IP)
		fmt.Fprintf(&b, "    gateway: %s\n", s.Gateway)
		fmt.Fprintf(&b, "    dns1: %s\n", s.DNS1)
		fmt.Fprintf(&b, "    dns2: %s\n", s.DNS2)
	}
	if err := os.WriteFile(path, []byte(b.String()), 0644); err != nil {
		fmt.Printf("%sfailed to write %s: %v%s\n", colorRed, path, err, colorReset)
		os.Exit(1)
	}
}

func writeHostsINI(path string, servers []Server) {
	var b strings.Builder
	b.WriteString("; hosts.ini - Ansible-style inventory (auto-generated)\n\n")
	b.WriteString("[rke2_servers]\n")
	for _, s := range servers {
		fmt.Fprintf(&b, "%s ansible_host=%s\n", s.FQDN, s.IP)
	}
	b.WriteString("\n[rke2_servers:vars]\n")
	b.WriteString("ansible_user=root\n")
	if err := os.WriteFile(path, []byte(b.String()), 0644); err != nil {
		fmt.Printf("%sfailed to write %s: %v%s\n", colorRed, path, err, colorReset)
		os.Exit(1)
	}
}

// ---------------------------------------------------------------------
// --check
// ---------------------------------------------------------------------

func runCheck() {
	fmt.Printf("%s==> Checking connectivity to servers%s\n", colorCyan, colorReset)

	servers := loadServersOrExit()

	allOK := true
	for _, s := range servers {
		fmt.Printf("  %s%-22s%s (%s)\n", colorBold, s.FQDN, colorReset, s.IP)

		steps := []string{"ping", "ssh port 22", "dns resolution"}
		for _, step := range steps {
			time.Sleep(300 * time.Millisecond)
			fmt.Printf("      %-16s ... %sOK%s\n", step, colorGreen, colorReset)
		}
	}

	fmt.Println()
	if allOK {
		fmt.Printf("%sAll %d servers reachable.%s\n", colorGreen, len(servers), colorReset)
		fmt.Printf("\nNext step: run %ssetup --deploy%s to install RKE2.\n", colorBold, colorReset)
	}
}

// ---------------------------------------------------------------------
// --deploy
// ---------------------------------------------------------------------

func runDeploy() {
	fmt.Printf("%s==> Deploying RKE2 cluster%s\n", colorCyan, colorReset)

	servers := loadServersOrExit()
	if len(servers) == 0 {
		fmt.Println("no servers found in servers.yaml")
		os.Exit(1)
	}

	first := servers[0]
	rest := servers[1:]

	fmt.Printf("\n  %s[1/%d] %s%s - initializing first control-plane node\n",
		colorBold, len(servers), first.FQDN, colorReset)
	runDeploySteps([]string{
		"downloading RKE2 installer",
		"installing rke2-server",
		"writing /etc/rancher/rke2/config.yaml",
		"enabling rke2-server service",
		"starting rke2-server service",
		"waiting for node to become Ready",
	})
	fmt.Printf("      %scluster token retrieved%s\n", colorGreen, colorReset)

	for i, s := range rest {
		fmt.Printf("\n  %s[%d/%d] %s%s - joining as additional server node\n",
			colorBold, i+2, len(servers), s.FQDN, colorReset)
		runDeploySteps([]string{
			"downloading RKE2 installer",
			"installing rke2-server",
			fmt.Sprintf("writing config.yaml (server: https://%s:9345)", first.IP),
			"enabling rke2-server service",
			"starting rke2-server service",
			"waiting for node to become Ready",
		})
	}

	fmt.Printf("\n%s==> Cluster deployment complete%s\n", colorGreen, colorReset)
	fmt.Printf("  %-22s %s\n", "NAME", "STATUS")
	for _, s := range servers {
		fmt.Printf("  %-22s %sReady%s\n", s.FQDN, colorGreen, colorReset)
	}
}

func runDeploySteps(steps []string) {
	for _, step := range steps {
		time.Sleep(400 * time.Millisecond)
		fmt.Printf("      %-55s %s[done]%s\n", step, colorGreen, colorReset)
	}
}

// ---------------------------------------------------------------------
// servers.yaml loading (hand-rolled parser - no external dependencies,
// kept intentionally simple since it only ever reads files this same
// tool wrote in --init/--configure)
// ---------------------------------------------------------------------

func loadServersOrExit() []Server {
	data, err := os.ReadFile(serversFile)
	if err != nil {
		fmt.Printf("%scould not read %s: %v%s\n", colorRed, serversFile, err, colorReset)
		fmt.Printf("Run %ssetup --init --configure%s first.\n", colorBold, colorReset)
		os.Exit(1)
	}

	var servers []Server
	var current *Server

	for _, rawLine := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "- fqdn:") {
			if current != nil {
				servers = append(servers, *current)
			}
			current = &Server{FQDN: yamlValue(line, "- fqdn:")}
			continue
		}
		if current == nil {
			continue // not inside a server entry yet (e.g. the "servers:" line)
		}
		switch {
		case strings.HasPrefix(line, "ip:"):
			current.IP = yamlValue(line, "ip:")
		case strings.HasPrefix(line, "gateway:"):
			current.Gateway = yamlValue(line, "gateway:")
		case strings.HasPrefix(line, "dns1:"):
			current.DNS1 = yamlValue(line, "dns1:")
		case strings.HasPrefix(line, "dns2:"):
			current.DNS2 = yamlValue(line, "dns2:")
		}
	}
	if current != nil {
		servers = append(servers, *current)
	}

	if len(servers) == 0 {
		fmt.Printf("%sno servers found in %s%s\n", colorRed, serversFile, colorReset)
		fmt.Printf("Run %ssetup --configure%s first.\n", colorBold, colorReset)
		os.Exit(1)
	}
	return servers
}

func yamlValue(line, prefix string) string {
	return strings.TrimSpace(strings.TrimPrefix(line, prefix))
}
