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
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/GehirnInc/crypt"
	api_errors "github.com/go-openapi/errors"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/config-parser/types"

	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/storagetype"
)

const mockPass = "$2a$10$e.I1x5KPNu7xy9u0zSzJaOcr5it8kR1Awnaf3boOtYno9y4DolER."

var usersStore *Users

var syncUserStore sync.Once

type Users struct {
	users []types.User
	mu    sync.Mutex
}

func GetUsersStore() *Users {
	syncUserStore.Do(func() {
		usersStore = &Users{}
	})
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
	u.mu.Lock()
	defer u.mu.Unlock()

	u.users = append(u.users, users...)
	return nil
}

func (u *Users) Init() error {
	configuration := Get()
	u.users = []types.User{}
	mode := configuration.Mode.Load()
	allUsers := configuration.GetUsers() // single + cluster mode
	if len(allUsers) > 0 {
		for _, storageUser := range allUsers {
			if mode != ModeCluster || strings.HasPrefix(storageUser.Name, storagetype.DapiClusterUserPrefix) {
				user := types.User{
					Name: storageUser.Name,
				}
				if storageUser.Password != nil {
					user.Password = *storageUser.Password
				}
				if storageUser.Insecure != nil {
					user.IsInsecure = *storageUser.Insecure
				}
				u.users = append(u.users, user)
			}
		}
		return nil
	}
	if mode == ModeCluster {
		return nil
	}
	if configuration.HAProxy.UserListFile != "" {
		errUserList := u.getUsersFromUsersListSection(configuration.HAProxy.UserListFile, configuration.HAProxy.Userlist)
		if errUserList != nil {
			return errUserList // file was specified, but errors exists, exit
		}
		return nil
	}
	return u.getUsersFromUsersListSection(configuration.HAProxy.ConfigFile, configuration.HAProxy.Userlist)
}

func (u *Users) AddUser(user types.User) error {
	clusterModeStorage := Get().GetClusterModeStorage()

	u.users = append(u.users, user)

	err := clusterModeStorage.AddUserAndStore(storagetype.User{
		Name:     user.Name,
		Insecure: &user.IsInsecure,
		Password: &user.Password,
	})
	return err
}

func (u *Users) RemoveUser(user types.User) error {
	clusterModeStorage := Get().GetClusterModeStorage()

	for i, v := range u.users {
		if v.Name == user.Name {
			u.users = slices.Delete(u.users, i, i+1)
			break
		}
	}
	err := clusterModeStorage.RemoveUserAndStore(storagetype.User{
		Name:     user.Name,
		Insecure: &user.IsInsecure,
		Password: &user.Password,
	})
	return err
}

func (u *Users) getUsersFromUsersListSection(filename, userlistSection string) error {
	// if file doesn't exists
	if _, err := os.Stat(filename); errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("cannot read %s, file does not exist", filename)
	}
	p, errP := parser.New(options.Path(filename))
	if errP != nil {
		return errP
	}
	data, err := p.Get(parser.UserList, userlistSection, "user")
	if err != nil {
		return fmt.Errorf("no users configured in %v, error: %w", filename, err)
	}

	return u.setUser(data, cfg.HAProxy.UserListFile)
}

// findUser searches user by its name. If found, returns user, otherwise returns an error.
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
		return nil, api_errors.New(http.StatusUnauthorized, "no configured users")
	}

	unatuhorized := false
	u, err := findUser(user, users)
	if err != nil {
		unatuhorized = true
	}

	userPass := mockPass
	if u != nil {
		userPass = u.Password
	}

	if strings.HasPrefix(userPass, "\"${") && strings.HasSuffix(userPass, "}\"") {
		userPass = os.Getenv(misc.ExtractEnvVar(userPass))
		if userPass == "" {
			unatuhorized = true
			userPass = mockPass
		}
	}

	if u != nil && u.IsInsecure {
		if pass != userPass {
			unatuhorized = true
		}
	} else {
		if !checkPassword(pass, userPass) {
			unatuhorized = true
		}
	}
	if unatuhorized {
		return nil, api_errors.New(http.StatusUnauthorized, "unauthorized")
	}
	return user, nil
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
