package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/fishworks/gofish"
	"github.com/fishworks/gofish/pkg/ohai"
	"github.com/spf13/cobra"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

func newLintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint <file...>",
		Short: "lint fish food",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var failed bool
			for _, arg := range args {
				l := lua.NewState()
				defer l.Close()
				if err := l.DoFile(arg); err != nil {
					return err
				}
				var food gofish.Food
				if err := gluamapper.Map(l.GetGlobal(strings.ToLower(reflect.TypeOf(food).Name())).(*lua.LTable), &food); err != nil {
					return err
				}
				errs := food.Lint()
				for _, err := range errs {
					ohai.Warningln(err)
				}
				if len(errs) != 0 {
					failed = true
					return fmt.Errorf("%d errors encountered while linting %s", len(errs), food.Name)
				}
			}
			if failed {
				return errors.New("linting failed")
			}
			return nil
		},
	}
	return cmd
}
