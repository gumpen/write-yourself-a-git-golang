package wyag

import "log"

type AddCommand struct{}

func (c *AddCommand) Help() string {
	return "add command"
}

func (c *AddCommand) Run(args []string) int {
	log.Println("Add!")
	return 0
}

func (c *AddCommand) Synopsis() string {
	return "Print \"ADD\""
}
