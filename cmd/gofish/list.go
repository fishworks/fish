package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/fishworks/gofish"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list installed fish food. If an argument is provided, list all installed versions of that fish food",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			table := uitable.New()
			if len(args) == 0 {
				table.AddRow("NAME")
				for _, food := range findFood() {
					table.AddRow(food)
				}
			} else {
				table.AddRow("NAME", "VERSION", "LINKED")
				for _, ver := range findFoodVersions(args[0]) {
					f := gofish.Food{
						Name:    args[0],
						Version: ver,
					}
					table.AddRow(f.Name, f.Version, f.Linked())
				}
			}
			fmt.Println(table)
			return nil
		},
	}
	return cmd
}

func findFood() []string {
	barrelPath := gofish.Home(gofish.HomePath).Barrel()
	var fudz []string
	files, err := ioutil.ReadDir(barrelPath)
	if err != nil {
		return []string{}
	}

	for _, f := range files {
		if f.IsDir() {
			files, err := ioutil.ReadDir(filepath.Join(barrelPath, f.Name()))
			if err != nil {
				continue
			}
			if len(files) > 0 {
				fileName := f.Name()
				rigConf, err := ioutil.ReadFile(filepath.Join(barrelPath, f.Name()) + "/rig.conf")
				if err == nil {
					location := strings.TrimSpace(string(rigConf))
					fileName = strings.Join([]string{location, fileName}, "/")
				}
				fudz = append(fudz, fileName)
			}
		}
	}
	return fudz
}

func findFoodVersions(name string) []string {
	barrelPath := gofish.Home(gofish.HomePath).Barrel()
	var versions []string
	files, err := ioutil.ReadDir(filepath.Join(barrelPath, name))
	if err != nil {
		return []string{}
	}

	for _, f := range files {
		if f.IsDir() {
			versions = append(versions, f.Name())
		}
	}
	return versions
}
