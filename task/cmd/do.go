package cmd

import (
	"fmt"
	"strconv"
	"task/db"

	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Marks a task as completed.",
	Run: func(cmd *cobra.Command, args []string) {
		//task do 1 2 3
		//因此首先解析参数，完成的任务 id
		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to parse the given argument: ", arg)
			} else {
				ids = append(ids, id)
			}
		}
		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println("Something went wrong: ", err)
			return
		}
		for _, i := range ids {
			if i < 1 || i > len(tasks) {
				continue
			}
			task := tasks[i-1]
			err := db.DeleteTask(task.Key)
			if err != nil {
				fmt.Printf("Failed to mark task \"%s\" as completed.\n", task.Value)
				continue
			} else {
				fmt.Printf("Marked task \"%s\" as completed.\n", task.Value)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}
