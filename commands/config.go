package commands

import (
	"fmt"
	"github.com/apache/brooklyn-client/api/entity_config"
	"github.com/apache/brooklyn-client/command_metadata"
	"github.com/apache/brooklyn-client/error_handler"
	"github.com/apache/brooklyn-client/net"
	"github.com/apache/brooklyn-client/scope"
	"github.com/apache/brooklyn-client/terminal"
	"github.com/codegangsta/cli"
)

type Config struct {
	network *net.Network
}

func NewConfig(network *net.Network) (cmd *Config) {
	cmd = new(Config)
	cmd.network = network
	return
}

func (cmd *Config) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "config",
		Description: "Show the config for an application or entity",
		Usage:       "BROOKLYN_NAME SCOPE config",
		Flags:       []cli.Flag{},
	}
}

func (cmd *Config) Run(scope scope.Scope, c *cli.Context) {
	if err := net.VerifyLoginURL(cmd.network); err != nil {
		error_handler.ErrorExit(err)
	}
	if c.Args().Present() {
		configValue, err := entity_config.ConfigValue(cmd.network, scope.Application, scope.Entity, c.Args().First())

		if nil != err {
			error_handler.ErrorExit(err)
		}
		displayValue, err := stringRepresentation(configValue)
		if nil != err {
			error_handler.ErrorExit(err)
		}
		fmt.Println(displayValue)

	} else {
		config, err := entity_config.ConfigCurrentState(cmd.network, scope.Application, scope.Entity)
		if nil != err {
			error_handler.ErrorExit(err)
		}
		table := terminal.NewTable([]string{"Key", "Value"})
		for key, value := range config {
			table.Add(key, fmt.Sprintf("%v", value))
		}
		table.Print()
	}
}
