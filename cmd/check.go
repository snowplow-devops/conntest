/*
 * Copyright (c) 2022 Snowplow Analytics Ltd. All rights reserved.
 *
 * This program is licensed to you under the Apache License Version 2.0,
 * and you may not use this file except in compliance with the Apache License Version 2.0.
 * You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the Apache License Version 2.0 is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/snowplow/conntest/pkg"
	"github.com/spf13/cobra"
)

var dsn string
var tags tagsVar
var retryTimes uint
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "A brief description of your command",
	Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		pkg.RegisterDatabricks()
		dsn, err := pkg.DB(dsn)
		if err == nil {
			event := pkg.Check(*dsn, tags, retryTimes)
			res, _ := json.Marshal(event)
			fmt.Println(string(res))
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

type tagsVar map[string]string

func (t *tagsVar) String() string {
	var str []string
	for key, value := range *t {
		str = append(str, fmt.Sprintf("%s=%v", key, value))
	}
	return strings.Join(str, ";")
}

func (t *tagsVar) Set(value string) error {
	splits := strings.Split(strings.Trim(value, ";"), ";")
	tags := map[string]string{}

	for _, split := range splits {
		split := strings.Split(split, "=")
		tags[split[0]] = split[1]
	}

	*t = tagsVar(tags)
	return nil
}

func (t *tagsVar) Type() string {
	return "tagsVar"
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&dsn, "dsn", "d", "", "database DSN")
	checkCmd.Flags().UintVarP(&retryTimes, "retry-times", "r", 1, "number of times to retry using exponential time")
	checkCmd.PersistentFlags().VarP(&tags, "tags", "", "optional tags")
}
