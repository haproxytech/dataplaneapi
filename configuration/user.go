// Copyright 2019 HAProxy Technologies
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

package configuration

import (
	"fmt"
	"os"
	"sync"

	"github.com/haproxytech/client-native/configuration"
	parser "github.com/haproxytech/config-parser/v2"
	"github.com/haproxytech/config-parser/v2/common"
	"github.com/haproxytech/config-parser/v2/types"
	"github.com/prometheus/common/log"
)

type User struct {
	mu     sync.Mutex
	users  []types.User
	Parser *parser.Parser
}

func (u *User) Get() []types.User {
	return u.users
}

func (u *User) setUser(data common.ParserData, file string) error {
	if data == nil {
		return fmt.Errorf("no users configured in %s", file)
	}

	users, ok := data.([]types.User)
	if !ok {
		return fmt.Errorf("error reading users from %s", file)
	}
	if len(users) == 0 {
		return fmt.Errorf("no users configured in %s", file)
	}
	u.users = users
	return nil
}

func (u *User) saveUser(userlist, file string, user common.ParserData) error {
	if err := u.Parser.LoadData(file); err != nil {
		return fmt.Errorf("cannot read %s, err: %s", file, err.Error())
	}
	err := u.Parser.SectionsCreate(parser.UserList, userlist)
	if err != nil {
		return fmt.Errorf("error creating section: %v", parser.UserList)
	}

	err = u.Parser.Set(parser.UserList, userlist, "user", user)
	if err != nil {
		return fmt.Errorf("error setting userlist %v", userlist)
	}

	err = u.Parser.Save(file)
	if err != nil {
		return fmt.Errorf("error setting userlist %v", userlist)
	}
	return nil
}

func (u *User) createUserFile(file string) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("file %s does not exist and cannot be created: %s", file, err.Error())
	}
	defer f.Close()
	return nil
}

func (u *User) Init(client *configuration.Client) error {
	if u.Parser == nil {
		return fmt.Errorf("parser not initialized")
	}

	cfg := Get()
	opt := cfg.HAProxy

	if opt.UserListFile != "" {
		//if userlist file doesn't exists
		if _, err := os.Stat(opt.UserListFile); os.IsNotExist(err) {
			//get user from HAProxy config file
			data, err := client.Parser.Get(parser.UserList, opt.Userlist, "user")
			if err != nil {
				return fmt.Errorf("error reading userlist %v userlist in conf: %s", opt.Userlist, err.Error())
			}
			err = u.createUserFile(opt.UserListFile)
			if err != nil {
				return err
			}
			err = u.saveUser(opt.Userlist, opt.UserListFile, data)
			if err != nil {
				return err
			}
			log.Infof("userlist saved to %s", opt.UserListFile)
			return u.setUser(data, opt.UserListFile)
		}
		//if userlist file exists
		if err := u.Parser.LoadData(opt.UserListFile); err != nil {
			return fmt.Errorf("cannot read %s, err: %s", opt.UserListFile, err.Error())
		}
		data, err := u.Parser.Get(parser.UserList, opt.Userlist, "user")
		if err != nil {
			return err
		}
		return u.setUser(data, opt.UserListFile)
	}

	//get user from HAProxy config
	user, err := client.Parser.Get(parser.UserList, opt.Userlist, "user")
	if err != nil {
		return fmt.Errorf("no users configured in %v, error: %s", opt.ConfigFile, err.Error())
	}
	return u.setUser(user, opt.ConfigFile)
}
