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
	"fmt"

	"github.com/spf13/cobra"
	"os"
	"text/template"
	"path/filepath"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		system := ""
		if len(args) == 0 {
			system = guessInitSystem()
			fmt.Fprintf(os.Stderr, "Guessed init system to be %s\n", system)
		} else {
			system = args[0]
		}

		DoInit(system)
	},
}

var script = `
# This is designed to be eval()'d by a shell in an EC2 userdata launch script.
# If you're reading this, it means you should probably add eval $(spotty init systemd)
# to your userdata.

cp {{ .currentPath }} {{ .targetPath }}
mkdir -p /etc/spotty
echo '[]' > /etc/spotty/spotty.yml

{{ if .isSystemd }}

cat << 'EOF' | tee /etc/systemd/system/spotty.service
[Unit]
Description={{ .description }}

[Service]
ExecStart={{ .targetPath }} run

[Install]
WantedBy=multi-user.target
EOF

systemctl enable spotty
systemctl start spotty

{{ else if .isUpstart }}

cat << 'EOF' | tee /etc/init/spotty.conf
description "{{ .description }}"

start on runlevel [2345]
exec {{ .targetPath }} run
EOF

service spotty start

{{ end }}
`

func getOwnPath() string {
	// this feels bad, but i'd rather not bring in cgo
	var here = os.Args[0]
	here, err := filepath.Abs(here)
	if err != nil {
		panic(err)
	}
	return here
}

func DoInit(system string) {
	isSystemd := ""
	isUpstart := ""

	switch system {
	case "systemd":
		isSystemd = "true"
	case "upstart":
		isUpstart = "true"
	default:
		panic("unsupported init system")
	}

	vals := map[string]string{
		"currentPath": getOwnPath(),
		"targetPath": "/usr/bin/spotty",
		"isSystemd": isSystemd,
		"isUpstart": isUpstart,
		"description": "A service to poll for EC2 spot termination notices and trigger events when one is received.",
	}

	tmpl, err := template.New("script").Parse(script)
	if err != nil { panic(err) }
	tmpl.Execute(os.Stdout, vals)
}

func guessInitSystem() string {
	dirExists := func(path string) bool {
		_, err := os.Stat(path)
		return !(err != nil && os.IsNotExist(err))
	}

	if dirExists("/usr/lib/systemd") {
		return "systemd"
	} else if dirExists("/usr/share/upstart") {
		return "upstart"
	} else {
		return ""
	}
}

func init() {
	RootCmd.AddCommand(initCmd)
}
