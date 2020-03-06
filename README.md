# modgraph

Modgraph is a tool to visualize the dependency graph of a Go module.

## dependencies

The tool requires the `go` tool and  `graphviz` to be installed. It also
assumes that `$BROWSER` env variable is set to the path of your default web
browser executable.

## usage

```
$ cd /path/to/modgraph
$ go build
$ cd path/to/module
$ /path/to/modgraph/modgraph
```

This will write the dependency graph to `/tmp/graph.svg` and open it in a
browser tab.
