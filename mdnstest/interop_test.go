// Copyright (C) 2026 The go-mdns Authors All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mdnstest

// The tests in this file check the responder against the mDNS clients the
// operating system provides, when they are installed: dns-sd (Bonjour, on
// macOS and Windows) and avahi-browse and avahi-resolve (Avahi, usually on
// Linux). Each test publishes a service under a random name with an
// mdns.Server, runs the client, and checks what it prints. A test is
// skipped when its client, or the daemon the client talks to, is not
// available, and with -short.

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/cybergarage/go-mdns/mdns"
)

const (
	interopServiceType = "_gomdnstest._tcp"
	interopSubtype     = "_gomdnssub"
	interopTimeout     = 10 * time.Second
)

// interopService is a service published for one test.
type interopService struct {
	server *mdns.Server
	svc    *mdns.LocalService
}

func startInteropService(t *testing.T) *interopService {
	t.Helper()
	suffix := randomSuffix(t)
	n, err := strconv.ParseUint(suffix[:4], 16, 16)
	if err != nil {
		t.Fatal(err)
	}
	svc := &mdns.LocalService{
		Instance:  "go-mdns-test-" + suffix,
		Service:   interopServiceType,
		Domain:    "",
		Subtypes:  []string{interopSubtype},
		Host:      "go-mdns-test-" + suffix,
		Port:      20000 + int(n%20000),
		TXT:       []string{"test=" + suffix, "flag"},
		Addresses: nil,
	}
	server := mdns.NewServer()
	if err := server.Start(); err != nil {
		t.Skipf("the server cannot bind the mDNS sockets here: %v", err)
	}
	t.Cleanup(func() { _ = server.Stop() })
	if err := server.Register(svc); err != nil {
		t.Fatal(err)
	}
	return &interopService{server: server, svc: svc}
}

func randomSuffix(t *testing.T) string {
	t.Helper()
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(b)
}

// requireTool skips the test when name is not installed, and with -short.
func requireTool(t *testing.T, name string) {
	t.Helper()
	if testing.Short() {
		t.Skip("the interoperability tests use the network; skipped with -short")
	}
	if _, err := exec.LookPath(name); err != nil {
		t.Skipf("%s is not installed", name)
	}
}

// requireAvahiDaemon skips the test when avahi-daemon is not reachable,
// which is common where the Avahi tools are installed without the daemon
// running, such as on macOS or in a container.
func requireAvahiDaemon(t *testing.T) {
	t.Helper()
	requireTool(t, "avahi-browse")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "avahi-browse", "--all", "--terminate", "--parsable", "--no-db-lookup").CombinedOutput()
	if err != nil {
		t.Skipf("avahi-daemon is not reachable: %v: %s", err, strings.TrimSpace(string(out)))
	}
}

// toolRun is a client process whose output is read line by line.
type toolRun struct {
	cmd   *exec.Cmd
	mutex sync.Mutex
	lines []string
	added chan struct{}
	done  chan struct{}
}

// startTool starts a client which keeps running, such as a browse. It is
// stopped with SIGTERM, or killed where signals are not supported, when the
// test ends.
func startTool(t *testing.T, name string, args ...string) *toolRun {
	t.Helper()
	cmd := exec.Command(name, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		t.Fatalf("start %s: %v", name, err)
	}
	run := &toolRun{
		cmd:   cmd,
		mutex: sync.Mutex{},
		lines: []string{},
		added: make(chan struct{}, 1),
		done:  make(chan struct{}),
	}
	go func() {
		defer close(run.done)
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			run.mutex.Lock()
			run.lines = append(run.lines, scanner.Text())
			run.mutex.Unlock()
			select {
			case run.added <- struct{}{}:
			default:
			}
		}
	}()
	t.Cleanup(func() { run.stop() })
	return run
}

func (run *toolRun) stop() {
	if err := run.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		_ = run.cmd.Process.Kill()
	}
	select {
	case <-run.done:
	case <-time.After(3 * time.Second):
		_ = run.cmd.Process.Kill()
		<-run.done
	}
	_ = run.cmd.Wait()
}

func (run *toolRun) output() string {
	run.mutex.Lock()
	defer run.mutex.Unlock()
	return strings.Join(run.lines, "\n")
}

// waitFor waits until the client prints a line match accepts, and returns
// it. Some clients buffer their output when it is not a terminal and print
// it only when they exit, so the client is stopped at the deadline and its
// last output is checked too.
func (run *toolRun) waitFor(t *testing.T, what string, match func(string) bool) string {
	t.Helper()
	find := func(from int) (string, int) {
		run.mutex.Lock()
		defer run.mutex.Unlock()
		for i := from; i < len(run.lines); i++ {
			if match(run.lines[i]) {
				return run.lines[i], i
			}
		}
		return "", len(run.lines)
	}
	deadline := time.After(interopTimeout)
	next := 0
	for {
		line, n := find(next)
		if line != "" {
			return line
		}
		next = n
		select {
		case <-run.added:
		case <-run.done:
			if line, _ := find(next); line != "" {
				return line
			}
			t.Fatalf("%s: the client exited without printing it; output:\n%s", what, run.output())
		case <-deadline:
			run.stop()
			if line, _ := find(next); line != "" {
				return line
			}
			t.Fatalf("%s: not printed within %s; output:\n%s", what, interopTimeout, run.output())
		}
	}
}

// runTool runs a client which exits by itself, such as a resolve, and
// returns its output.
func runTool(t *testing.T, name string, args ...string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), interopTimeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	if err != nil && !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Fatalf("%s %s: %v\n%s", name, strings.Join(args, " "), err, out)
	}
	return string(out)
}

func containsAll(line string, parts ...string) bool {
	for _, p := range parts {
		if !strings.Contains(line, p) {
			return false
		}
	}
	return true
}

// Avahi

// avahiFields splits a line of avahi-browse --parsable output. A resolved
// line reads
// "=;eth0;IPv4;instance;type description;local;host.local;192.0.2.2;8080;"txt""
// and a removed one "-;eth0;IPv4;instance;type description;local".
func avahiFields(line string) []string {
	return strings.Split(line, ";")
}

func avahiResolvedLine(svc *mdns.LocalService) func(string) bool {
	return func(line string) bool {
		f := avahiFields(line)
		return 9 < len(f) && f[0] == "=" && f[3] == svc.Instance
	}
}

func TestAvahiBrowseAndResolve(t *testing.T) {
	requireAvahiDaemon(t)
	s := startInteropService(t)

	for _, browseType := range []string{interopServiceType, interopSubtype + "._sub." + interopServiceType} {
		t.Run(browseType, func(t *testing.T) {
			run := startTool(t, "avahi-browse", "--parsable", "--resolve", "--no-db-lookup", browseType)
			line := run.waitFor(t, "the resolved service", avahiResolvedLine(s.svc))
			f := avahiFields(line)
			if f[6] != s.svc.HostName() {
				t.Errorf("host = %q, want %q", f[6], s.svc.HostName())
			}
			if f[8] != strconv.Itoa(s.svc.Port) {
				t.Errorf("port = %q, want %d", f[8], s.svc.Port)
			}
			txt := strings.Join(f[9:], ";")
			for _, want := range s.svc.TXT {
				if !strings.Contains(txt, `"`+want+`"`) {
					t.Errorf("TXT %s does not hold %q", txt, want)
				}
			}
		})
	}
}

func TestAvahiResolveHost(t *testing.T) {
	requireAvahiDaemon(t)
	requireTool(t, "avahi-resolve")
	s := startInteropService(t)

	out := runTool(t, "avahi-resolve", "--name", s.svc.HostName())
	fields := strings.Fields(out)
	if len(fields) < 2 || fields[0] != s.svc.HostName() {
		t.Fatalf("avahi-resolve --name %s printed %q, want the host and an address", s.svc.HostName(), out)
	}
}

func TestAvahiSeesGoodbye(t *testing.T) {
	requireAvahiDaemon(t)
	s := startInteropService(t)

	run := startTool(t, "avahi-browse", "--parsable", "--no-db-lookup", interopServiceType)
	run.waitFor(t, "the added service", func(line string) bool {
		f := avahiFields(line)
		return 3 < len(f) && f[0] == "+" && f[3] == s.svc.Instance
	})
	if err := s.server.Deregister(s.svc); err != nil {
		t.Fatal(err)
	}
	run.waitFor(t, "the removed service", func(line string) bool {
		f := avahiFields(line)
		return 3 < len(f) && f[0] == "-" && f[3] == s.svc.Instance
	})
}

// dns-sd

func TestDNSSDBrowse(t *testing.T) {
	requireTool(t, "dns-sd")
	s := startInteropService(t)

	for _, browseType := range []string{interopServiceType, interopServiceType + "," + interopSubtype} {
		t.Run(browseType, func(t *testing.T) {
			// "... Add  2  4 local.  _gomdnstest._tcp.  go-mdns-test-0a1b2c3d"
			run := startTool(t, "dns-sd", "-B", browseType, "local.")
			run.waitFor(t, "the added service", func(line string) bool {
				return containsAll(line, " Add ", s.svc.Instance)
			})
		})
	}
}

func TestDNSSDResolve(t *testing.T) {
	requireTool(t, "dns-sd")
	s := startInteropService(t)

	// "go-mdns-test-0a1b2c3d._gomdnstest._tcp.local. can be reached at
	// go-mdns-test-0a1b2c3d.local.:20123 (interface 4)", followed by the TXT
	// strings.
	run := startTool(t, "dns-sd", "-L", s.svc.Instance, interopServiceType, "local.")
	reach := s.svc.HostName() + ".:" + strconv.Itoa(s.svc.Port)
	run.waitFor(t, "the host and port", func(line string) bool {
		return containsAll(line, "can be reached at", reach)
	})
	run.waitFor(t, "the TXT strings", func(line string) bool {
		return containsAll(line, s.svc.TXT[0])
	})
}

func TestDNSSDResolveHost(t *testing.T) {
	requireTool(t, "dns-sd")
	s := startInteropService(t)

	// "... Add  2  4 go-mdns-test-0a1b2c3d.local.  192.0.2.2  120"
	run := startTool(t, "dns-sd", "-G", "v4v6", s.svc.HostName())
	run.waitFor(t, "an address of the host", func(line string) bool {
		return containsAll(line, " Add ", s.svc.HostName())
	})
}

func TestDNSSDSeesGoodbye(t *testing.T) {
	requireTool(t, "dns-sd")
	s := startInteropService(t)

	run := startTool(t, "dns-sd", "-B", interopServiceType, "local.")
	run.waitFor(t, "the added service", func(line string) bool {
		return containsAll(line, " Add ", s.svc.Instance)
	})
	if err := s.server.Deregister(s.svc); err != nil {
		t.Fatal(err)
	}
	run.waitFor(t, "the removed service", func(line string) bool {
		return containsAll(line, " Rmv ", s.svc.Instance)
	})
}
