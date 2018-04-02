package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fishworks/fish"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

type rigListCmd struct{}

func newRigListCmd() *cobra.Command {
	rcmd := &rigListCmd{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list rigs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return rcmd.run()
		},
	}
	return cmd
}

func (r *rigListCmd) run() error {
	rigPath := fish.Home(fish.HomePath).Rigs()
	rigs := findRigs(rigPath)

	table := uitable.New()
	table.AddRow("NAME")
	for _, rig := range rigs {
		table.AddRow(rig)
	}
	fmt.Println(table)
	return nil
}

func findRigs(dir string) []string {
	var rigs []string
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() && f.Name() == "Food" {
			rigName := strings.TrimPrefix(filepath.Dir(path), dir+string(os.PathSeparator))
			rigs = append(rigs, rigName)
		}
		return nil
	})
	return rigs
}
