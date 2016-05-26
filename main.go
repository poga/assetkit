// TOOD: revision control
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var themePath = "themes/summit"

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var RootCmd = &cobra.Command{
	Use:   "suisui",
	Short: "SuiSui manage your assets beautifully",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of SuiSui",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0.0.1")
	},
}

var compileCmd = &cobra.Command{
	Use:   "compile path",
	Short: "Compile a project into standalone website",
	Run: func(cmd *cobra.Command, args []string) {
		projectPath := strings.TrimSuffix(args[0], string(os.PathSeparator))
		project, err := NewProject(projectPath)
		if err != nil {
			log.Fatal(err)
		}

		err = project.CompileTo(args[1])
		if err != nil {
			log.Fatal(err)
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status path",
	Short: "Status of a project",
	Run: func(cmd *cobra.Command, args []string) {
		projectPath := strings.TrimSuffix(args[0], string(os.PathSeparator))
		project, err := NewProject(projectPath)
		if err != nil {
			log.Fatal(err)
		}

		st, err := project.Status()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", st)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd, compileCmd, statusCmd)
}
