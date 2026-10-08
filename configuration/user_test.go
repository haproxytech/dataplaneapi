// Copyright 2026 HAProxy Technologies
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
	"testing"
	"time"

	// register the SHA-256 crypt implementation, as the main package does
	_ "github.com/GehirnInc/crypt/sha256_crypt"
	_ "github.com/GehirnInc/crypt/sha512_crypt"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

func TestAuthenticateUser(t *testing.T) {
	store := GetUsersStore()
	saved := store.users
	t.Cleanup(func() { store.users = saved })
	store.users = []types.User{
		{Name: "insecure_user", Password: "plainpass", IsInsecure: true},
		// SHA-256 crypt of "secretpass"
		{Name: "crypted_user", Password: "$5$testsalt$eCSeYT8Aub0tlGnCdlCmGO4RrnbXQZlcDzFHJWzOPa6"},
	}

	tests := []struct {
		name    string
		user    string
		pass    string
		wantErr bool
	}{
		{name: "insecure user correct password", user: "insecure_user", pass: "plainpass"},
		{name: "insecure user wrong password", user: "insecure_user", pass: "wrong", wantErr: true},
		{name: "crypted user correct password", user: "crypted_user", pass: "secretpass"},
		{name: "crypted user wrong password", user: "crypted_user", pass: "wrong", wantErr: true},
		{name: "unknown user", user: "nosuchuser", pass: "plainpass", wantErr: true},
		// an insecure user's stored plaintext must not pass the crypt check for another user
		{name: "crypted user given plaintext of stored hash", user: "crypted_user", pass: "$5$testsalt$eCSeYT8Aub0tlGnCdlCmGO4RrnbXQZlcDzFHJWzOPa6", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := AuthenticateUser(tt.user, tt.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthenticateUser(%q, %q) error = %v, wantErr %v", tt.user, tt.pass, err, tt.wantErr)
			}
		})
	}
}

func TestAuthenticateUserNoConfiguredUsers(t *testing.T) {
	store := GetUsersStore()
	saved := store.users
	t.Cleanup(func() { store.users = saved })
	store.users = nil

	if _, err := AuthenticateUser("anyone", "anything"); err == nil {
		t.Error("AuthenticateUser succeeded with an empty user store")
	}
}

// Failing logins must take comparable time for known and unknown users.
func TestAuthenticateUser_NoUserEnumerationByTiming(t *testing.T) {
	store := GetUsersStore()
	saved := store.users
	t.Cleanup(func() { store.users = saved })
	store.users = []types.User{
		{Name: "crypted_user", Password: "$5$testsalt$eCSeYT8Aub0tlGnCdlCmGO4RrnbXQZlcDzFHJWzOPa6"},
	}
	measure := func(user string) time.Duration {
		start := time.Now()
		for range 50 {
			_, _ = AuthenticateUser(user, "wrong")
		}
		return time.Since(start) / 50
	}
	known, unknown := measure("crypted_user"), measure("nosuchuser")
	t.Logf("wrong password: known user %v, unknown user %v (ratio %.0fx)", known, unknown, float64(known)/float64(unknown))
	if known > 5*unknown {
		t.Fatal("known users are distinguishable from unknown users by response time")
	}
}
