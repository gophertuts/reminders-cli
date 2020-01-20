package client

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

// CLI constants
const (
	// commands
	CreateCmd = "create"
	EditCmd   = "edit"
	FetchCmd  = "fetch"
	DeleteCmd = "delete"

	// flags
	TitleFlag    = "title"
	MessageFlag  = "message"
	DurationFlag = "duration"
	IDFlag       = "id"
)

var (
	// HelpMsg represents the help message for the CLI client
	HelpMsg = fmt.Sprintf(
		"Usage of %s:\n<command> [<args>]\n%s\n%s\n%s\n%s\n",
		os.Args[0],
		"create\t --help",
		"edit\t --help",
		"fetch\t --help",
		"delete\t --help",
	)

	// commands
	createCmd = flag.NewFlagSet(CreateCmd, flag.ExitOnError)
	fetchCmd  = flag.NewFlagSet(FetchCmd, flag.ExitOnError)
	editCmd   = flag.NewFlagSet(EditCmd, flag.ExitOnError)
	deleteCmd = flag.NewFlagSet(DeleteCmd, flag.ExitOnError)

	// flags
	idFlag       flagList
	titleFlag    string
	messageFlag  string
	durationFlag time.Duration

	flagShortcuts = map[string]string{
		TitleFlag:    "t",
		MessageFlag:  "m",
		DurationFlag: "d",
	}
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

// Switch represents CLI command switch
type Switch struct {
	Args   []string
	Client BackendHTTPClient
}

// NewSwitch creates a new instance of command Switch
func NewSwitch(args []string, uri string) Switch {
	httpClient := NewHTTPClient(uri)
	return Switch{
		Args:   args,
		Client: httpClient,
	}
}

// Switch analyses the CLI args and executes the given command
func (s Switch) Switch() {
	// 1st arg 		- executable
	// 2nd arg 		- the command
	// 3rd+ args 	- the command flags
	switch s.Args[1] {
	case CreateCmd:
		s.reminderFlags(createCmd)
		s.checkArgs(3)
		s.parseCmd(createCmd)
		res := s.Client.Create(titleFlag, messageFlag, durationFlag)
		fmt.Printf("reminder created successfully:\n%s\n", res)
	case EditCmd:
		editCmd.Var(&idFlag, IDFlag, "The ID (int) of the reminder to edit")
		s.reminderFlags(editCmd)
		s.checkArgs(2)
		s.parseCmd(editCmd)
		res := s.Client.Edit(idFlag[len(idFlag)-1], titleFlag, messageFlag, durationFlag)
		fmt.Printf("reminder edited successfully:\n%s\n", res)
	case FetchCmd:
		fetchCmd.Var(&idFlag, IDFlag, "List of reminder IDs (int) to fetch")
		s.checkArgs(1)
		s.parseCmd(fetchCmd)
		res := s.Client.Fetch(idFlag)
		fmt.Printf("reminders fetched successfully:\n%s\n", res)
	case DeleteCmd:
		deleteCmd.Var(&idFlag, IDFlag, "List of reminder IDs (int) to delete")
		s.checkArgs(1)
		s.parseCmd(deleteCmd)
		err := s.Client.Delete(idFlag)
		if err != nil {
			log.Fatalf("could not delete record(s):\n%v\n%v\n", idFlag, err)
		}
		fmt.Printf("successfully deleted record(s):\n%v\n", idFlag)
	default:
		fmt.Printf("%q is not a valid command.\n", s.Args[1])
		os.Exit(2)
	}
}

// reminderFlags configures reminder specific flags for a command
func (s Switch) reminderFlags(f *flag.FlagSet) {
	f.StringVar(&titleFlag, TitleFlag, "", "Reminder title")
	f.StringVar(&titleFlag, flagShortcuts[TitleFlag], "", "Reminder title")
	f.StringVar(&messageFlag, MessageFlag, "", "Reminder message")
	f.StringVar(&messageFlag, flagShortcuts[MessageFlag], "", "Reminder message")
	f.DurationVar(&durationFlag, DurationFlag, 0, "Reminder time")
	f.DurationVar(&durationFlag, flagShortcuts[DurationFlag], 0, "Reminder time")
}

// checkArgs checks if the number of passed in args is greater or equal to min args
func (s Switch) checkArgs(minArgs int) {
	if len(s.Args) == 3 && s.Args[2] == "--help" {
		return
	}
	if len(s.Args)-2 < minArgs {
		fmt.Printf(
			"incorect use of %s\n%s %s --help\n",
			s.Args[1], s.Args[0], s.Args[1],
		)
		fmt.Printf(
			"%s expects at least: %d arg(s), %d provided\n",
			s.Args[1], minArgs, len(s.Args)-2,
		)
		os.Exit(2)
	}
}

// parseCmd parses sub-command flags
func (s Switch) parseCmd(cmd *flag.FlagSet) {
	err := cmd.Parse(s.Args[2:])
	if err != nil {
		log.Fatalf("could not parse flags for command: %s\n%v", cmd.Name(), err)
	}
}
