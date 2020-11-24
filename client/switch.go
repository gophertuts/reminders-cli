package client

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// idsFlag represents []string values passed from CLI
type idsFlag []string

func (list *idsFlag) String() string {
	return strings.Join(*list, ",")
}

func (list *idsFlag) Set(v string) error {
	*list = append(*list, v)
	return nil
}

// BackendHTTPClient represents the HTTP client for communicating with the Backend API
type BackendHTTPClient interface {
	Create(title, message string, duration time.Duration) ([]byte, error)
	Edit(id string, title, message string, duration time.Duration) ([]byte, error)
	Fetch(ids []string) ([]byte, error)
	Delete(ids []string) error
	Healthy(host string) bool
}

// NewSwitch creates a new instance of command Switch
func NewSwitch(uri string) Switch {
	httpClient := NewHTTPClient(uri)
	s := Switch{client: httpClient, backendAPIURL: uri}
	s.commands = map[string]func() func(string) error{
		"create": s.create,
		"edit":   s.edit,
		"fetch":  s.fetch,
		"delete": s.delete,
		"health": s.health,
	}
	return s
}

// Switch represents CLI command switch
type Switch struct {
	client        BackendHTTPClient
	backendAPIURL string
	commands      map[string]func() func(string) error
}

// Switch analyses the CLI args and executes the given command
func (s Switch) Switch() error {
	cmdName := os.Args[1]
	cmd, ok := s.commands[cmdName]
	if !ok {
		return fmt.Errorf("invalid command '%s'", cmdName)
	}
	return cmd()(cmdName)
}

// Help prints a useful message about command usage
func (s Switch) Help() {
	var help string
	for name := range s.commands {
		help += name + "\t --help\n"
	}
	fmt.Printf("Usage of %s:\n<command> [<args>]\n%s", os.Args[0], help)
}

// create represents the create command which creates a new reminder
func (s Switch) create() func(string) error {
	return func(cmd string) error {
		createCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		t, m, d := s.reminderFlags(createCmd)

		if err := s.checkArgs(3); err != nil {
			return err
		}
		if err := s.parseCmd(createCmd); err != nil {
			return err
		}

		res, err := s.client.Create(*t, *m, *d)
		if err != nil {
			return wrapError("could not create reminder", err)
		}
		fmt.Printf("reminder created successfully:\n%s", string(res))
		return nil
	}
}

// edit represents the edit command which edit a reminder
func (s Switch) edit() func(string) error {
	return func(cmd string) error {
		ids := idsFlag{}
		editCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		editCmd.Var(&ids, "id", "The ID (int) of the reminder to edit")
		t, m, d := s.reminderFlags(editCmd)

		if err := s.checkArgs(2); err != nil {
			return err
		}
		if err := s.parseCmd(editCmd); err != nil {
			return err
		}

		lastID := ids[len(ids)-1]
		res, err := s.client.Edit(lastID, *t, *m, *d)
		if err != nil {
			return wrapError("could not edit reminder", err)
		}
		fmt.Printf("reminder edited successfully:\n%s", string(res))
		return nil
	}
}

// fetch represents the fetch command which fetches a list of reminders
func (s Switch) fetch() func(string) error {
	return func(cmd string) error {
		ids := idsFlag{}
		fetchCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		fetchCmd.Var(&ids, "id", "List of reminder IDs (int) to fetch")

		if err := s.checkArgs(1); err != nil {
			return err
		}
		if err := s.parseCmd(fetchCmd); err != nil {
			return err
		}

		res, err := s.client.Fetch(ids)
		if err != nil {
			return wrapError("could not fetch reminder(s)", err)
		}
		fmt.Printf("reminders fetched successfully:\n%s", string(res))
		return nil
	}
}

// delete represents the delete command which deletes a reminder
func (s Switch) delete() func(string) error {
	return func(cmd string) error {
		ids := idsFlag{}
		deleteCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		deleteCmd.Var(&ids, "id", "List of reminder IDs (int) to delete")

		if err := s.checkArgs(1); err != nil {
			return err
		}
		if err := s.parseCmd(deleteCmd); err != nil {
			return err
		}

		err := s.client.Delete(ids)
		if err != nil {
			return wrapError("could not delete reminder(s)", err)
		}
		fmt.Printf("successfully deleted record(s):\n%v\n", ids)
		return nil
	}
}

// health represents the health command which prints whether a host is healthy or not
func (s Switch) health() func(string) error {
	return func(cmd string) error {
		var host string
		healthCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		healthCmd.StringVar(&host, "host", s.backendAPIURL, "Host to ping for health")
		if err := s.parseCmd(healthCmd); err != nil {
			return err
		}
		if !s.client.Healthy(host) {
			fmt.Printf("host: %s is down\n", host)
		} else {
			fmt.Printf("host: %s is up and running\n", host)
		}
		return nil
	}
}

// reminderFlags configures reminder specific flags for a command
func (s Switch) reminderFlags(f *flag.FlagSet) (*string, *string, *time.Duration) {
	t, m, d := "", "", time.Duration(0)
	f.StringVar(&t, "title", "", "Reminder title")
	f.StringVar(&t, "t", "", "Reminder title")
	f.StringVar(&m, "message", "", "Reminder message")
	f.StringVar(&m, "m", "", "Reminder message")
	f.DurationVar(&d, "duration", 0, "Reminder time")
	f.DurationVar(&d, "d", 0, "Reminder time")
	return &t, &m, &d
}

// parseCmd parses sub-command flags
func (s Switch) parseCmd(cmd *flag.FlagSet) error {
	err := cmd.Parse(os.Args[2:])
	if err != nil {
		return wrapError("could not parse '"+cmd.Name()+"' command flags", err)
	}
	return nil
}

// checkArgs checks if the number of passed in args is greater or equal to min args
func (s Switch) checkArgs(minArgs int) error {
	if len(os.Args) == 3 && os.Args[2] == "--help" {
		return nil
	}
	if len(os.Args)-2 < minArgs {
		fmt.Printf(
			"incorect use of %s\n%s %s --help\n",
			os.Args[1], os.Args[0], os.Args[1],
		)
		return fmt.Errorf(
			"%s expects at least: %d arg(s), %d provided",
			os.Args[1], minArgs, len(os.Args)-2,
		)
	}
	return nil
}
