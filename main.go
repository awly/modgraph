package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	out, err := exec.Command("go", "mod", "graph").Output()
	if err != nil {
		fatalf(`executing "go mod graph": %v`, err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	sb := &strings.Builder{}

	fmt.Fprintln(sb, "digraph deps {")
	fmt.Fprintln(sb, "\tratio = \"4\"")
	for _, l := range lines {
		parts := strings.Fields(l)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, `unexpected "go mod graph" output: %q`, l)
			continue
		}
		fmt.Fprintf(sb, "\t%q -> %q\n", parts[0], parts[1])
	}
	fmt.Fprintln(sb, "}")

	out, err = exec.Command("dot", "-Tsvg", "-o", "/tmp/deps.svg").CombinedOutput()
	if err != nil {
		fatalf(`executing "dot -Tsvg -o /tmp/deps.svg": %v
output: %s`, err, out)
	}

	browser := os.Getenv("BROWSER")
	if len(browser) == 0 {
		fatalf("$BROWSER not set, can't view dependency graph\nopen /tmp/deps.svg any way you can")
	}
	out, err = exec.Command(browser, "/tmp/deps.svg").CombinedOutput()
	if err != nil {
		fatalf(`executing "$BROWSER /tmp/deps.svg": %v
output: %s`, err, out)
	}
}

func fatalf(f string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, f, args...)
	os.Exit(1)
}
