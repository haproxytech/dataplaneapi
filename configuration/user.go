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
	"strings"
	"sync"

	"github.com/haproxytech/dataplaneapi/misc"

	"github.com/GehirnInc/crypt"
	api_errors "github.com/go-openapi/errors"
	parser "github.com/haproxytech/config-parser/v3"
	"github.com/haproxytech/config-parser/v3/common"
	"github.com/haproxytech/config-parser/v3/types"
	log "github.com/sirupsen/logrus"
)

var usersStore *Users

type Users struct {
	mu    sync.Mutex
	users []types.User
}

func GetUsersStore() *Users {
	if usersStore == nil {
		usersStore = &Users{}
		if err := usersStore.Init(); err != nil {
			log.Fatalf("Error initiating users: %s", err.Error())
		}
	}
	return usersStore
}

func (u *Users) GetUsers() []types.User {
	return u.users
}

func (u *Users) setUser(data common.ParserData, file string) error {
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

func (u *Users) saveUsers(userlist, file string, user common.ParserData) error {
	p := &parser.Parser{}
	if err := p.LoadData(file); err != nil {
		return fmt.Errorf("cannot read %s, err: %s", file, err.Error())
	}
	err := p.SectionsCreate(parser.UserList, userlist)
	if err != nil {
		return fmt.Errorf("error creating section: %v", parser.UserList)
	}

	err = p.Set(parser.UserList, userlist, "user", user)
	if err != nil {
		return fmt.Errorf("error setting userlist %v", userlist)
	}

	err = p.Save(file)
	if err != nil {
		return fmt.Errorf("error setting userlist %v", userlist)
	}
	return nil
}

func (u *Users) createUserFile(file string) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("file %s does not exist and cannot be created: %s", file, err.Error())
	}
	defer f.Close()
	return nil
}

func (u *Users) Init() error {
	cfg := Get()
	p := &parser.Parser{}
	if cfg.HAProxy.UserListFile != "" {
		//if userlist file doesn't exists
		if _, err := os.Stat(cfg.HAProxy.UserListFile); os.IsNotExist(err) {
			//get user from HAProxy config file
			if err := p.LoadData(cfg.HAProxy.ConfigFile); err != nil {
				return fmt.Errorf("cannot read %s, err: %s", cfg.HAProxy.ConfigFile, err.Error())
			}
			data, err := p.Get(parser.UserList, cfg.HAProxy.Userlist, "user")
			if err != nil {
				return fmt.Errorf("error reading userlist %v userlist in conf: %s", cfg.HAProxy.ConfigFile, err.Error())
			}
			err = u.createUserFile(cfg.HAProxy.UserListFile)
			if err != nil {
				return err
			}
			err = u.saveUsers(cfg.HAProxy.Userlist, cfg.HAProxy.UserListFile, data)
			if err != nil {
				return err
			}
			return u.setUser(data, cfg.HAProxy.UserListFile)
		}
		//if userlist file exists
		if err := p.LoadData(cfg.HAProxy.UserListFile); err != nil {
			return fmt.Errorf("cannot read %s, err: %s", cfg.HAProxy.UserListFile, err.Error())
		}
		data, err := p.Get(parser.UserList, cfg.HAProxy.Userlist, "user")
		if err != nil {
			return fmt.Errorf("no users configured in %v, error: %s", cfg.HAProxy.UserListFile, err.Error())
		}
		return u.setUser(data, cfg.HAProxy.UserListFile)
	}
	//get user from HAProxy config
	if err := p.LoadData(cfg.HAProxy.ConfigFile); err != nil {
		return fmt.Errorf("cannot read %s, err: %s", cfg.HAProxy.ConfigFile, err.Error())
	}
	user, err := p.Get(parser.UserList, cfg.HAProxy.Userlist, "user")
	if err != nil {
		return fmt.Errorf("no users configured in %v, error: %s", cfg.HAProxy.ConfigFile, err.Error())
	}
	return u.setUser(user, cfg.HAProxy.ConfigFile)
}

//findUser searches user by its name. If found, returns user, otherwise returns an error.
func findUser(userName string, users []types.User) (*types.User, error) {
	for _, u := range users {
		if u.Name == userName {
			return &u, nil
		}
	}
	return nil, api_errors.New(401, "no configured users")
}

func AuthenticateUser(user string, pass string) (interface{}, error) {
	users := GetUsersStore().GetUsers()
	if len(users) == 0 {
		return nil, api_errors.New(401, "no configured users")
	}

	u, err := findUser(user, users)
	if err != nil {
		return nil, err
	}

	userPass := u.Password
	if strings.HasPrefix(u.Password, "\"${") && strings.HasSuffix(u.Password, "}\"") {
		userPass = os.Getenv(misc.ExtractEnvVar(userPass))
		if userPass == "" {
			return nil, api_errors.New(401, fmt.Sprintf("%s %s", "can not read password from env variable:", u.Password))
		}
	}

	if u.IsInsecure {
		if pass == userPass {
			return user, nil
		}
		return nil, api_errors.New(401, fmt.Sprintf("%s %s", "invalid password:", pass))
	}
	if checkPassword(pass, userPass) {
		return user, nil
	}
	return nil, api_errors.New(401, fmt.Sprintf("%s %s", "invalid password:", pass))
}

func checkPassword(pass, storedPass string) bool {
	parts := strings.Split(storedPass, "$")
	if len(parts) == 4 {
		var c crypt.Crypter
		switch parts[1] {
		case "1":
			c = crypt.MD5.New()
		case "5":
			c = crypt.SHA256.New()
		case "6":
			c = crypt.SHA512.New()
		default:
			return false
		}
		if err := c.Verify(storedPass, []byte(pass)); err == nil {
			return true
		}
	}

	return false
}
