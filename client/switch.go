package client

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	// flags names
	titleFlag    = "title"
	messageFlag  = "message"
	durationFlag = "duration"
	idFlag       = "id"
)

var (
	// flags variables
	idFlagVar       flagList
	titleFlagVar    string
	messageFlagVar  string
	durationFlagVar time.Duration
)

// flagList represents []int values for CLI flags
type flagList []string

func (i *flagList) String() string {
	return "my string representation"
}

func (i *flagList) Set(v string) error {
	*i = append(*i, v)
	return nil
}

// BackendHTTPClient represents the HTTP client for communicating with the Backend API
type BackendHTTPClient interface {
	Create(title, message string, duration time.Duration) Reminder
	Edit(id string, title, message string, duration time.Duration) Reminder
	Fetch(ids []string) []Reminder
	Delete(ids []string) error
}

// NewSwitch creates a new instance of command Switch
func NewSwitch(help bool, uri string) Switch {
	httpClient := NewHTTPClient(uri)
	s := Switch{client: httpClient}
	s.commands = map[string]func() func(string){
		"create": s.create,
		"edit":   s.edit,
		"fetch":  s.fetch,
		"delete": s.delete,
	}
	if help || len(os.Args) == 1 {
		s.help()
	}
	return s
}

// Switch represents CLI command switch
type Switch struct {
	client   BackendHTTPClient
	commands map[string]func() func(string)
}

// Switch analyses the CLI args and executes the given command
func (s Switch) Switch() {
	cmdName := os.Args[1]
	cmd, ok := s.commands[os.Args[1]]
	if !ok {
		s.help()
	}
	cmd()(cmdName)
}

func (s Switch) create() func(string) {
	return func(cmdName string) {
		createCmd := flag.NewFlagSet(cmdName, flag.ExitOnError)
		s.setReminderFlags(createCmd)
		s.checkArgs(3)
		s.parseCmd(createCmd)
		res := s.client.Create(titleFlagVar, messageFlagVar, durationFlagVar)
		fmt.Printf("reminder created successfully:\n%s\n", res)
	}
}

func (s Switch) edit() func(string) {
	return func(cmdName string) {
		editCmd := flag.NewFlagSet(cmdName, flag.ExitOnError)
		editCmd.Var(&idFlagVar, idFlag, "The ID (int) of the reminder to edit")
		s.setReminderFlags(editCmd)
		s.checkArgs(2)
		s.parseCmd(editCmd)
		lastID := idFlagVar[len(idFlagVar)-1]
		res := s.client.Edit(lastID, titleFlagVar, messageFlagVar, durationFlagVar)
		fmt.Printf("reminder edited successfully:\n%s\n", res)
	}
}

func (s Switch) fetch() func(string) {
	return func(cmdName string) {
		fetchCmd := flag.NewFlagSet(cmdName, flag.ExitOnError)
		fetchCmd.Var(&idFlagVar, idFlag, "List of reminder IDs (int) to fetch")
		s.checkArgs(1)
		s.parseCmd(fetchCmd)
		res := s.client.Fetch(idFlagVar)
		fmt.Printf("reminders fetched successfully:\n%s\n", res)
	}
}

func (s Switch) delete() func(string) {
	return func(cmdName string) {
		deleteCmd := flag.NewFlagSet(cmdName, flag.ExitOnError)
		deleteCmd.Var(&idFlagVar, idFlag, "List of reminder IDs (int) to delete")
		s.checkArgs(1)
		s.parseCmd(deleteCmd)
		err := s.client.Delete(idFlagVar)
		if err != nil {
			log.Fatalf("could not delete record(s):\n%v\n%v\n", idFlagVar, err)
		}
		fmt.Printf("successfully deleted record(s):\n%v\n", idFlagVar)
	}
}

func (s Switch) help() {
	var help string
	for name := range s.commands {
		help += name + "\t --help\n"
	}
	fmt.Printf("Usage of %s:\n<command> [<args>]\n%s", os.Args[0], help)
	os.Exit(2)
}

// reminderFlags configures reminder specific flags for a command
func (s Switch) setReminderFlags(f *flag.FlagSet) {
	f.StringVar(&titleFlagVar, titleFlag, "", "Reminder title")
	f.StringVar(&titleFlagVar, "t", "", "Reminder title")
	f.StringVar(&messageFlagVar, messageFlag, "", "Reminder message")
	f.StringVar(&messageFlagVar, "m", "", "Reminder message")
	f.DurationVar(&durationFlagVar, durationFlag, 0, "Reminder time")
	f.DurationVar(&durationFlagVar, "d", 0, "Reminder time")
}

// parseCmd parses sub-command flags
func (s Switch) parseCmd(cmd *flag.FlagSet) {
	err := cmd.Parse(os.Args[2:])
	if err != nil {
		log.Fatalf("could not parse flags for command: %s\n%v", cmd.Name(), err)
	}
}

// checkArgs checks if the number of passed in args is greater or equal to min args
func (s Switch) checkArgs(minArgs int) {
	if len(os.Args) == 3 && os.Args[2] == "--help" {
		return
	}
	if len(os.Args)-2 < minArgs {
		fmt.Printf(
			"incorect use of %s\n%s %s --help\n",
			os.Args[1], os.Args[0], os.Args[1],
		)
		fmt.Printf(
			"%s expects at least: %d arg(s), %d provided\n",
			os.Args[1], minArgs, len(os.Args)-2,
		)
		os.Exit(2)
	}
}
