package cmd

import (
	"fmt"
	"github.com/houyanzu/goforge/controllergen"
	"github.com/houyanzu/goforge/projectgen"
	"github.com/houyanzu/goforge/routergen"
	"github.com/houyanzu/goforge/toolfunc"
	"github.com/spf13/cobra"
	"os"
)

var root string
var methods string

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&root, "root", "app/api/home/", "root dir")
	rootCmd.PersistentFlags().StringVar(&methods, "methods", "", "methods")
}

func initConfig() {
	// Configuration initialization code if needed
}

var rootCmd = &cobra.Command{
	Use:   "goforge",
	Short: "A tool for generating Go project components",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: please input action")
			return
		}

		action := args[0]

		switch action {
		case "addController":
			if len(args) < 2 {
				fmt.Println("Error: please input controller route")
				return
			}
			route := args[1]
			if err := toolfunc.ValidateFileName(route); err != nil {
				fmt.Println(err)
				return
			}
			controllergen.OperateController(root, route, action, methods)
		case "addMethods":
			if len(args) < 2 {
				fmt.Println("Error: please input controller route")
				return
			}
			route := args[1]
			if err := toolfunc.ValidateFileName(route); err != nil {
				fmt.Println(err)
				return
			}
			controllergen.OperateController(root, route, action, methods)
		case "routergen":
			routergen.Routergen(root)
		case "init":
			if len(args) < 2 {
				fmt.Println("Error: please input project name")
				return
			}
			projectgen.InitProject(args[1])
		default:
			fmt.Println("action not found")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
