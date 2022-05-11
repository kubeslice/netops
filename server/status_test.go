//go:build !status
// +build !status

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
