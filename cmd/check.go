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

	"github.com/snowplow/conntest/pkg"
	"github.com/spf13/cobra"
)

var dsn string
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "A brief description of your command",
	Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var tags = map[string]string{}

		tagsStr, terr := cmd.Flags().GetString("tags")
		if terr == nil && tagsStr != "" {
			tags = pkg.ParseTags(tagsStr)
		}

		dsn, err := pkg.DB(dsn)
		if err == nil {
			event := pkg.Check(*dsn, tags)
			res, _ := json.Marshal(event)
			fmt.Println(string(res))
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&dsn, "dsn", "d", "", "database DSN")
}
