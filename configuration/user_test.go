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

	_ "github.com/GehirnInc/crypt/sha256_crypt"
	_ "github.com/GehirnInc/crypt/sha512_crypt"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

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
