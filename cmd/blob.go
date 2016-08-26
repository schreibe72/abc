// Copyright © 2016 Manfred Schreiber <software@manfredschreiber.de>
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
	"errors"

	"github.com/spf13/cobra"
)

// blobCmd represents the blob command
var blobCmd = &cobra.Command{
	Use:   "blob",
	Short: "All Blob manipulating commands",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if key == "" {
			return errors.New("Key should not be empty!")
		}
		if account == "" {
			return errors.New("Account shoud not be empty!")
		}
		if container == "" {
			return errors.New("Container shoud not be empty!")
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(blobCmd)
}
