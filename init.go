package wyag

import (
	"log"
)

type InitCommand struct{}

func (c *InitCommand) Help() string {
	return "Where to create the repository"
}

func (c *InitCommand) Run(args []string) int {
	var path string
	if len(args) == 0 {
		path = "."
	} else if len(args) > 1 {
		log.Fatal("Too many argument")
		return 1
	} else {
		path = args[0]
	}

	_, err := repoCreate(path)
	if err != nil {
		log.Fatal(err)
		return 1
	}

	log.Println("Initialize the repository!")
	return 0
}

func (c *InitCommand) Synopsis() string {
	return "Print \"Init\""
}
