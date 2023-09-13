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

package runtime

import (
	"sync"

	commands "github.com/haproxytech/dataplaneapi/runtime/commands"
)

type Commands struct {
	cmdS     map[string]commands.Command
	muCmds   sync.Mutex
	help     commands.Help
	init     sync.Once
	initHelp sync.Once
}

func (c *Commands) Get(key string) (commands.Command, bool) {
	if key == "help" {
		c.initHelp.Do(func() {
			c.help.Header = []byte("Dataplaneapi runtime commands:\n")
			c.help.Footer = []byte("\ntype help <command> for more info")
		})
		return c.help, true
	}
	cmd, ok := c.cmdS[key]
	return cmd, ok
}

func (c *Commands) Register(cmd commands.Command) {
	c.init.Do(func() {
		c.cmdS = map[string]commands.Command{}
	})
	c.muCmds.Lock()
	defer c.muCmds.Unlock()
	c.cmdS[cmd.Definition().Key] = cmd
	c.help.Commands = append(c.help.Commands, cmd.Definition())
}

func (c *Commands) UnRegister(cmdKey string) {
	c.muCmds.Lock()
	defer c.muCmds.Unlock()
	delete(c.cmdS, cmdKey)
}

type Command interface {
	Definition() definition
	Command(cmd []string) (response []byte, err error)
}

type definition struct {
	Key      string
	Info     string
	Commands []allCommands
}

type allCommands struct {
	Command string
	Info    string
}
