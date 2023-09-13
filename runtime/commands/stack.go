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
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/maruel/panicparse/v2/stack"
)

type Stack struct{}

func (g Stack) Definition() definition {
	return definition{
		Key:  "stack",
		Info: "output stack trace",
		Commands: []allCommands{
			{"stack", "output stack trace"},
			{"stack raw", "output stack trace - go stack trace"},
		},
	}
}

func (g Stack) Command(cmd []string) (response []byte, err error) {
	if len(cmd) < 2 || cmd[1] != "raw" {
		dmp, err := MakeStackDump()
		if err != nil {
			return []byte{}, err
		}
		return []byte(dmp), nil
	}
	buffSize := int(128 * 1e6)
	buff := make([]byte, buffSize)
	runtime.Stack(buff, true)
	trace := bytes.Trim(buff, "\x00") // this dramatically speeds up output

	return trace, nil
}

func MakeStackDump() (string, error) {
	buffSize := int(128 * 1e6)

	buff := make([]byte, buffSize)
	runtime.Stack(buff, true)

	var result strings.Builder
	trace := bytes.Trim(buff, "\x00")

	s, _, err := stack.ScanSnapshot(bytes.NewReader(trace), io.Discard, stack.DefaultOpts())
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	// Find out similar goroutine traces and group them into buckets.
	buckets := s.Aggregate(stack.AnyValue).Buckets

	// Calculate alignment.
	srcLen := 0
	pkgLen := 0
	for _, bucket := range buckets {
		for _, line := range bucket.Signature.Stack.Calls {
			if l := len(fmt.Sprintf("%s:%d", line.SrcName, line.Line)); l > srcLen {
				srcLen = l
			}
			if l := len(filepath.Base(line.Func.ImportPath)); l > pkgLen {
				pkgLen = l
			}
		}
	}

	for _, bucket := range buckets {
		// Print the goroutine header.
		extra := ""
		if s := bucket.SleepString(); s != "" {
			extra += " [" + s + "]"
		}
		if bucket.Locked {
			extra += " [locked]"
		}

		if len(bucket.CreatedBy.Calls) != 0 {
			extra += fmt.Sprintf(" [Created by %s.%s @ %s:%d]", bucket.CreatedBy.Calls[0].Func.DirName, bucket.CreatedBy.Calls[0].Func.Name, bucket.CreatedBy.Calls[0].SrcName, bucket.CreatedBy.Calls[0].Line)
		}
		result.WriteString(fmt.Sprintf("%d: %s%s\n", len(bucket.IDs), bucket.State, extra))

		// Print the stack lines.
		for _, line := range bucket.Stack.Calls {
			arg := line.Args
			result.WriteString(fmt.Sprintf(
				"    %-*s %-*s %s(%s)\n",
				pkgLen, line.Func.DirName, srcLen,
				fmt.Sprintf("%s:%d", line.SrcName, line.Line),
				line.Func.Name, &arg))
		}
		if bucket.Stack.Elided {
			result.WriteString("    (...)\n")
		}
	}
	return result.String(), nil
}
