package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func main() {
	out, err := exec.Command("go", "mod", "graph").Output()
	if err != nil {
		fatalf(`executing "go mod graph": %v`, err)
	}
	g := make(graph)
	var root string
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for i, l := range lines {
		parts := strings.Fields(l)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, `unexpected "go mod graph" output: %q`, l)
			continue
		}
		g.add(parts[0], parts[1])
		if i == 0 {
			root = parts[0]
		}
	}
	g.sortDeps(root)
	dotGraph := g.dot(root)

	cmd := exec.Command("dot", "-Tsvg", "-o", "/tmp/deps.svg")
	cmd.Stdin = strings.NewReader(dotGraph)
	out, err = cmd.CombinedOutput()
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

type node struct {
	name   string
	deps   []*node
	weight int
	seen   bool
}

func (n *node) calcWeight() int {
	if n.weight == 0 {
		w := 1
		for _, dep := range n.deps {
			w += dep.calcWeight()
		}
		n.weight = w
	}

	return n.weight
}

type graph map[string]*node

func (g graph) add(from, to string) {
	fn, ok := g[from]
	if !ok {
		fn = &node{name: from}
	}
	tn, ok := g[to]
	if !ok {
		tn = &node{name: to}
	}
	fn.deps = append(fn.deps, tn)

	g[from] = fn
	g[to] = tn
}

func (g graph) sortDeps(root string) {
	g[root].calcWeight()

	for _, n := range g {
		sort.Slice(n.deps, func(i, j int) bool {
			// Sort by number of indirect dependencies.
			return n.deps[i].weight > n.deps[j].weight
		})
	}
}

func (g graph) dot(root string) string {
	sb := &strings.Builder{}
	fmt.Fprintln(sb, "digraph deps {")
	fmt.Fprintln(sb, "\trankstep = \"4.0\"")
	q := []*node{g[root]}
	for len(q) > 0 {
		var n *node
		n, q = q[0], q[1:]
		for _, dep := range n.deps {
			if dep.seen {
				continue
			}
			fmt.Fprintf(sb, "\t%q -> %q [label=\"%d\"]\n", n.name, dep.name, dep.weight)
			dep.seen = true
			q = append(q, dep)
		}
	}

	fmt.Fprintln(sb, "}")
	return sb.String()
}
