// mauflag - An extendable command-line argument parser for Golang
// Copyright (C) 2016 Maunium
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package flag

import (
	"fmt"
	"os"
	"strings"
)

// Command is the command that was executed
var Command string

// Args contains the arguments that weren't affiliated with any flag
var Args []string

// Parse the command line arguments into mauflag form
func Parse() error {
	var flag *Flag
	var key string
	var noMoreFlags = false
	Command = os.Args[0]
	os.Args = os.Args[1:]

	for _, arg := range os.Args {
		arg = strings.ToLower(arg)
		if noMoreFlags {
			Args = append(Args, arg)
		} else if arg == "--" {
			noMoreFlags = true
		} else if flag != nil {
			err := flag.setValue(arg)
			if err != nil {
				return fmt.Errorf("Flag %s was not a %s", key, flag.Value.Name())
			}
			flag = nil
		} else if arg[0] == '-' {
			var err error
			key, flag, err = flagStart(arg)
			if err != nil {
				return err
			}
		} else {
			Args = append(Args, arg)
		}
	}
	return nil
}

func flagStart(arg string) (string, *Flag, error) {
	key := arg[1:]

	var val string
	if strings.ContainsRune(key, '=') {
		val = key[strings.Index(arg, "=")+1:]
		key = key[:strings.Index(arg, "=")]
	}

	flag, key := getFlag(key)
	if flag == nil {
		return "", nil, fmt.Errorf("Unknown flag: %s", key)
	} else if len(val) > 0 {
		flag.setValue(val)
		return "", nil, nil
	} else {
		_, ok := flag.Value.(*boolValue)
		if ok {
			flag.setValue("true")
			return "", nil, nil
		}
	}
	return key, flag, nil
}

func getFlag(key string) (*Flag, string) {
	if key[0] == '-' {
		key = key[1:]
		for _, lflag := range flags {
			for _, lkey := range lflag.longKeys {
				if lkey == key {
					return lflag, lkey
				}
			}
		}
	} else {
		for _, lflag := range flags {
			for _, lkey := range lflag.shortKeys {
				if lkey == key {
					return lflag, lkey
				}
			}
		}
	}
	return nil, key
}
