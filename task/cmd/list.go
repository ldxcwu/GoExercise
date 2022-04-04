package cmd

import (
	"fmt"
	"os"
	"task/db"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "L ist all task that we have now.",
	Run: func(cmd *cobra.Command, args []string) {
		//print all the tasks if the tasks list is not empty
		//1. call the AllTasks func
		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println("Something went wrong!")
			os.Exit(1)
		}
		if len(tasks) == 0 {
			fmt.Println("You have no task to completed. Why not take a vocation? 🏄🏾‍♂️")
			return
		}
		for i, task := range tasks {
			fmt.Printf("%d. %s\n", i+1, task.Value)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
