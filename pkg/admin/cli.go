package admin

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CLI struct {
	client *Client
}

func NewCLI(socketPath string) (*CLI, error) {
	client, err := NewClient(socketPath)
	if err != nil {
		return nil, err
	}

	return &CLI{
		client,
	}, nil
}

// GetAppsNames returns the list of all apps
func (c *CLI) GetAppsNames() error {
	appNames, err := c.client.GetAppsNames()
	if err != nil {
		return err
	}

	for _, name := range appNames {
		fmt.Println(name)
	}

	return nil
}

// DeleteApp deletes an app if a matching app exists
func (c *CLI) DeleteApp(appname string) error {
	// since this is a very destructive action
	// we ask the user to type it out the app name as a form of validation
	fmt.Println(fmt.Sprintf("Are you sure you want to delete the app '%s'? This action can not be reversed.", appname))
	fmt.Println("")
	fmt.Println("Keep in mind the following:")
	fmt.Println("a) If an agent is still running, the app will be recreated.")
	fmt.Println("b) The API is idempotent, ie. if the app already does NOT exist, this command will run just fine.")
	fmt.Println("")
	fmt.Println(fmt.Sprintf("Type '%s' to confirm (without quotes).", appname))
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	trimmed := strings.TrimRight(text, "\n")
	if trimmed != appname {
		return fmt.Errorf("The app typed does not match. Want '%s' but got '%s'", appname, trimmed)
	}

	// finally delete the app
	err = c.client.DeleteApp(appname)
	if err != nil {
		return fmt.Errorf("failed to delete app: %w", err)
	}

	fmt.Println(fmt.Sprintf("Deleted app '%s'.", appname))
	return nil
}