package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"taskmanager/internal"

	"github.com/urfave/cli/v2"
)

const (
	flagCompleted = "completed"
	flagPending   = "pending"
)

func parseId(idStr string) (uint, error) {
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func convertFilter(completed bool, pending bool) internal.ListFilter {
	if completed == true {
		return internal.FilterCompleted
	} else if pending == true {
		return internal.FilterPending
	} else {
		return internal.FilterAll
	}
}

func main() {

	storage := internal.NewStorageJson()

	app := &cli.App{
		Name:        "todo",
		Description: "A simple todo manager",
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a task to the list",
				Action: func(cCtx *cli.Context) error {
					description := cCtx.Args().First()
					storage.AddTask(description)
					fmt.Println("added task: ", description)
					return nil
				},
			},
			{
				Name:    "complete",
				Aliases: []string{"c"},
				Usage:   "complete a task on the list",
				Action: func(cCtx *cli.Context) error {
					taskId, err := parseId(cCtx.Args().First())
					if err != nil {
						return err
					}
					storage.CompleteTask(taskId)
					fmt.Println("completed task: ", taskId)
					return nil
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  flagCompleted,
						Usage: "List completed tasks",
					},
					&cli.BoolFlag{
						Name:  flagPending,
						Usage: "List pending tasks",
					},
				},
				Usage: "list all tasks",
				Action: func(cCtx *cli.Context) error {
					filter := convertFilter(cCtx.Bool(flagCompleted), cCtx.Bool(flagPending))
					tasks, err := storage.ListTasks(filter)
					if err != nil {
						return err
					}
					for _, task := range tasks {
						statusEmoji := "❌"
						if task.Status == internal.StatusCompleted {
							statusEmoji = "✅"
						}
						fmt.Printf("%d: %s %s\n", task.ID, task.Description, statusEmoji)
					}
					return nil
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "delete a task on the list",
				Action: func(cCtx *cli.Context) error {
					taskId, err := parseId(cCtx.Args().First())
					if err != nil {
						return err
					}
					storage.DeleteTask(taskId)
					fmt.Println("deleted task: ", taskId)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
