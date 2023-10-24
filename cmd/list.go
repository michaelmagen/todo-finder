package cmd

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/michaelmagen/todo-finder/util/todo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list <directory>",
	Short: "List all todo comments found in directory",
	Long: `List all todo comments found in directory.
	
A specific directory can be passed in to search. If no directory is passed in, then the current directory is search.
Files and directories ignored by git are not included in search, but this can be disabled by -g flag. 
Hidden files/directories are also ignored, and this can also be disable with -a.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			log.Error("Too many arguments!")
			return
		}

		var dirToSearch string
		if len(args) == 0 {
			wd, err := os.Getwd()
			if err != nil {
				log.Error("Failed to get current working directory:", err)
				return
			}
			dirToSearch = wd
		} else {
			dirToSearch = args[0]
		}
		todos, err := todoFinder.GetTodos(dirToSearch)

		if err != nil {
			log.Error(err)
			return
		}

		if len(todos) == 0 {
			fmt.Println("No Todos were found in ", dirToSearch)
			return
		}
		// Show list of todos in terminal
		todoFinder.CreateTodoList(todos, dirToSearch)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("hidden", "a", false, "Include hidden files/directories in search")
	viper.BindPFlag("hidden", listCmd.Flags().Lookup("hidden"))
	listCmd.Flags().BoolP("no-gitignore", "g", false, "Include files listed in .gitignore")
	viper.BindPFlag("no-gitignore", listCmd.Flags().Lookup("no-gitignore"))
}
