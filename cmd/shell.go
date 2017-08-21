// Copyright © 2017 Aidan Steele <aidan.steele@glassechidna.com.au>
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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/glassechidna/spotty/common"
	"fmt"
	"os"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "No args passed in. You must provide a shell command to execute upon notification of spot instance termination.")
			os.Exit(1)
		} else if len(args) > 1 {
			fmt.Fprintln(os.Stderr, "Multiple args passed in. This usually means that you didn't escape the shell command you want executed.")
			os.Exit(1)
		}

		common.DefaultActionsList().Add("shell", args[0])
	},
}

func init() {
	RootCmd.AddCommand(shellCmd)
}
