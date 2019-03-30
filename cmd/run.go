// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/zwzn/jmd/run"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <jmd file>",
	Short: "runs a jmd file and outputs to stdout",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inFile := args[0]
		outFile, _ := cmd.Flags().GetString("out")
		return runFile(inFile, outFile)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("out", "o", "./out.md", "the file to output")
}

func runFile(inFile, outFile string) error {
	out, err := run.RunFile(inFile)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outFile, []byte(out), 0644)
	if err != nil {
		return err
	}
	return nil
}
