//go:build !status
// +build !status

/*  Copyright (c) 2022 Avesha, Inc. All rights reserved.
 *
 *  SPDX-License-Identifier: Apache-2.0
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package server

import (
	"testing"

	"github.com/kubeslice/netops/logger"
)

func TestRunCommand(t *testing.T) {
	testCases := []struct {
		testCase  string
		cmdString string
		expected  string
		errString string
	}{
		{
			"Length validation",
			"",
			"",
			"No command defined : ",
		},
		{
			"command Execution Test",
			`echo "hello"`,
			"hello\n",
			"",
		},
	}
	logger.GlobalLogger = logger.NewLogger("ERROR")
	for _, tt := range testCases {
		returned, err := runCommand(tt.cmdString)
		if err != nil {
			t.Log(err.Error())
			if err.Error() != tt.errString {
				t.Error("Expected :", tt.errString, " but got ", err)
			}
		}
		if returned != tt.expected {
			t.Error("Expected ", tt.expected, " but got ", returned)
		}
	}
}
