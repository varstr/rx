package main

import (
	"kylelemons.net/go/rx/graph"
)

var tagsCmd = &Command{
	Name:    "tags",
	Usage:   "<repo>",
	Summary: "List known repository tags.",
	Help: `The tags command scans the specified repository and lists
information about its tags.  The <repo> can be the suffix of the repository
root path, as long as it is unique.

The -f option takes a template as a format.  The data passed into the
template invocation is an (rx/graph) TagList, and the default format is:

` + ind2sp(tagsTemplate),
}

var (
	tagsFormat = tagsCmd.Flag.String("f", "", "tags output format")
	tagsLong   = tagsCmd.Flag.Bool("long", false, "Use long output format")
	tagsUp     = tagsCmd.Flag.Bool("up", false, "Only show updates (overrides --down)")
	tagsDown   = tagsCmd.Flag.Bool("down", false, "Only show downgrades")
)

func tagsFunc(cmd *Command, args ...string) {
	switch len(args) {
	case 0:
		args = append(args, "all")
	case 1:
	default:
		cmd.BadArgs("too many arguments")
	}
	pathSuffix := args[0]

	// Scan before accessing Graph
	if err := Scan(); err != nil {
		cmd.Fatalf("scan: %s", err)
	}

	repo, err := Deps.FindRepo(pathSuffix)
	if err != nil {
		cmd.Fatalf("<repo>: %s", err)
	}

	var tags graph.TagList
	switch {
	case *tagsUp:
		tags, err = repo.Upgrades()
	case *tagsDown:
		tags, err = repo.Downgrades()
	default:
		tags, err = repo.Tags()
	}
	if err != nil {
		cmd.Fatalf("list tags for %q: %s", repo.Root, err)
	}

	switch {
	case *tagsFormat != "":
		render(stdout, *tagsFormat, tags)
	default:
		render(stdout, tagsTemplate, tags)
	}
}

func init() {
	tagsCmd.Run = tagsFunc
}

var (
	tagsTemplate = `{{range .}}{{.Rev}} {{.Name}}
{{end}}`
)
