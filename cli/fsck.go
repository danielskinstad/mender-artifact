// Copyright 2023 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

//go:build !darwin

package cli

import (
	"os/exec"
	"syscall"

	"github.com/pkg/errors"

	"github.com/mendersoftware/mender-artifact/utils"
)

// From the fsck man page:
// The exit code returned by fsck is the sum of the following conditions:
//
//	0      No errors
//	1      Filesystem errors corrected
//	2      System should be rebooted
//	4      Filesystem errors left uncorrected
//	8      Operational error
//	16     Usage or syntax error
//	32     Checking canceled by user request
//	128    Shared-library error
func runFsck(image, fstype string) error {
	bin, err := utils.GetBinaryPath("fsck." + fstype)
	if err != nil {
		return errors.Wrap(err, "fsck command not found")
	}
	cmd := exec.Command(bin, "-a", image)
	if err := cmd.Run(); err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			if ws.ExitStatus() == 0 || ws.ExitStatus() == 1 {
				return nil
			}
			if ws.ExitStatus() == 8 {
				return errFsTypeUnsupported
			}
			return errors.Wrap(err, "fsck error")
		}
		return errors.New("fsck returned unparsed error")
	}
	return nil
}
