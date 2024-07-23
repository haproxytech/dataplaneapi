// Copyright 2024 HAProxy Technologies
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

package storage

import (
	"path/filepath"
	"testing"

	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/dataplaneapi/storagetype"
	"github.com/stretchr/testify/require"
)

func TestAddUser(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "add 1 user",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			clusterStoragePath := filepath.Join(dir, ClusterModeDataFileName)

			// create log store with initial configuration
			fs := &fileStorage[storagetype.ClusterModeData]{clusterStoragePath}
			users := setup(t, fs)
			clusterStorage, err := NewClusterModeStorage(clusterStoragePath)
			require.NoError(t, err)

			// apply change on users slice
			addUser(&users)

			// now add this user to storage
			err = clusterStorage.AddUserAndStore(userToAdd())
			require.NoError(t, err, "failed to add user")

			// check change has been applied
			require.Equal(t, users, clusterStorage.GetUsers())
			// reload storage and check
			err = clusterStorage.Load()
			require.NoError(t, err, "failed to load users")
			require.Equal(t, users, clusterStorage.GetUsers())
		})
	}
}

func TestRemoveUser(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "remove one user",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			clusterStoragePath := filepath.Join(dir, ClusterModeDataFileName)

			// create log store with initial configuration
			fs := &fileStorage[storagetype.ClusterModeData]{clusterStoragePath}
			users := setup(t, fs)
			clusterStorage, err := NewClusterModeStorage(clusterStoragePath)
			require.NoError(t, err)

			// apply change on users slice
			removeUser(&users)

			// now add this user to storage
			err = clusterStorage.RemoveUserAndStore(userToRemove())
			require.NoError(t, err, "failed to remove user")

			// check change has been applied
			require.Equal(t, users, clusterStorage.GetUsers())
			// reload storage and check
			err = clusterStorage.Load()
			require.NoError(t, err, "failed to load users")
			require.Equal(t, users, clusterStorage.GetUsers())
		})
	}
}

func TestReplaceAllUsers(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "replace all users",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			clusterStoragePath := filepath.Join(dir, ClusterModeDataFileName)

			// create log store with initial configuration
			fs := &fileStorage[storagetype.ClusterModeData]{clusterStoragePath}
			_ = setup(t, fs)
			clusterStorage, err := NewClusterModeStorage(clusterStoragePath)
			require.NoError(t, err)

			// now add this user to storage
			err = clusterStorage.ReplaceAllUsersAndStore(usersToReplace())
			require.NoError(t, err, "failed to replace all users")

			// check change has been applied
			require.Equal(t, usersToReplace(), clusterStorage.GetUsers())
			// reload storage and check
			err = clusterStorage.Load()
			require.NoError(t, err, "failed to load users")
			require.Equal(t, usersToReplace(), clusterStorage.GetUsers())
		})
	}
}

func setup(t *testing.T, fs *fileStorage[storagetype.ClusterModeData]) storagetype.Users {
	t.Helper()
	users := getSetupUsers()
	clusterMode := storagetype.ClusterModeData{
		Users: users,
	}
	err := fs.Store(clusterMode)
	require.NoError(t, err, "failed to write setup users")
	return users
}

func getSetupUsers() storagetype.Users {
	return storagetype.Users{
		storagetype.User{
			Name:     "admin",
			Password: misc.StringP("admin"),
		},
		storagetype.User{
			Name:     "user1",
			Password: misc.StringP("user1"),
		},
		storagetype.User{
			Name:     "user2",
			Password: misc.StringP("user2"),
		},
	}
}

func addUser(users *storagetype.Users) {
	*users = append(*users, userToAdd())
}

func removeUser(users *storagetype.Users) {
	*users = (*users)[:len(*users)-1]
}

func userToAdd() storagetype.User {
	return storagetype.User{
		Name:     "useradded",
		Password: misc.StringP("useradded"),
		Insecure: misc.BoolP(true),
	}
}

func userToRemove() storagetype.User {
	return storagetype.User{
		Name:     "user2",
		Password: misc.StringP("user2"),
	}
}

func usersToReplace() storagetype.Users {
	return storagetype.Users{
		storagetype.User{
			Name:     "toreplace1",
			Password: misc.StringP("toreplace1"),
		},
		storagetype.User{
			Name:     "toreplace2",
			Password: misc.StringP("replace2"),
		},
	}
}
