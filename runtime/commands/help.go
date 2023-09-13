// Copyright 2023 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package commands

import (
	"fmt"
	"strconv"
	"strings"
)

type Help struct {
	Header   []byte
	Footer   []byte
	Commands []definition
}

func (g Help) Definition() definition {
	return definition{
		Key:  "help",
		Info: "show help",
		Commands: []allCommands{
			{"help", "show help"},
		},
	}
}

func (g Help) Command(cmd []string) (response []byte, err error) {
	if len(cmd) < 2 {
		return g.getGeneralHelp(), nil
	}
	return g.getCommandHelp(cmd[1]), nil
}

func (g Help) getGeneralHelp() (response []byte) {
	response = append(response, g.Header...)
	for _, c := range g.Commands {
		line := fmt.Sprintf("%15s   %s\n", c.Key, c.Info)
		response = append(response, []byte(line)...)
	}

	response = append(response, g.Footer...)
	return response
}

func (g Help) getCommandHelp(command string) (response []byte) {
	response = append(response, []byte("Dataplaneapi runtime\n\ncommand "+command+":\n")...)
	found := false
	for _, c := range g.Commands {
		cmd := strings.Split(c.Commands[0].Command, " ")
		if cmd[0] == command {
			maxLen := 0
			for _, cmd := range c.Commands {
				if len(cmd.Command) > maxLen {
					maxLen = len(cmd.Command)
				}
			}
			format := "%-" + strconv.Itoa(maxLen+2) + "s %s\n"
			for _, cmd := range c.Commands {
				found = true
				line := fmt.Sprintf(format, cmd.Command, cmd.Info)
				response = append(response, []byte(line)...)
			}
		}
	}
	if !found {
		return g.getGeneralHelp()
	}

	return response
}
