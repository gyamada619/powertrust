package main

import (
	"github.com/gyamada619/powertrust/command"
	"github.com/mitchellh/cli"
)

func Commands(meta *command.Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"service": func() (cli.Command, error) {
			return &command.ServiceCommand{
				Meta: *meta,
			}, nil
		},
		"sign": func() (cli.Command, error) {
			return &command.SignCommand{
				Meta: *meta,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Meta:     *meta,
				Version:  Version,
				Revision: GitCommit,
				Name:     Name,
			}, nil
		},
	}
}
