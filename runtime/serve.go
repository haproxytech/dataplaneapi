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
	"net"
	"strings"

	"github.com/haproxytech/dataplaneapi/log"
)

func serve(comm *Commands, conn net.Conn) {
	buf := make([]byte, 512)
	nr, err := conn.Read(buf)
	defer conn.Close()
	if err != nil {
		log.Error("-- command socket: " + err.Error())
		return
	}

	data := buf[0:nr]
	log.Debugf("-- command socket got: %s", data)

	cmd := strings.Fields(string(data))
	c, ok := comm.Get(cmd[0])
	if cmd[0] == "exit" {
		return
	}
	if !ok {
		c, _ = comm.Get("help")
	}

	rsp, err := c.Command(cmd)
	if err != nil {
		_, e := conn.Write([]byte(err.Error()))
		if e != nil {
			log.Error("-- command socket write: " + e.Error())
		}
		return
	}
	_, e := conn.Write(rsp)
	if e != nil {
		log.Error("-- command socket write: " + e.Error())
	}
	if len(rsp) < 1 || rsp[len(rsp)-1] != '\n' {
		if _, err = conn.Write([]byte("\n")); err != nil {
			log.Error("-- command socket write: " + err.Error())
		}
	}
}
