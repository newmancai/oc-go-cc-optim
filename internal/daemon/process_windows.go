// oc-go-cc — Anthropic-to-OpenAI proxy for OpenCode Go + Claude Code
//
// Copyright (C) 2026  Samuel Tuyizere
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

//go:build windows

package daemon

import (
	"fmt"
	"os"
	"syscall"
)

const windowsSynchronize = 0x00100000

// IsProcessRunning checks if a process with the given PID is running.
func IsProcessRunning(pid int) bool {
	handle, err := syscall.OpenProcess(windowsSynchronize, false, uint32(pid))
	if err != nil {
		return false
	}
	defer func() { _ = syscall.CloseHandle(handle) }()

	event, err := syscall.WaitForSingleObject(handle, 0)
	return err == nil && event == syscall.WAIT_TIMEOUT
}

// StopProcess terminates a process on Windows.
// Unlike the Unix implementation which sends SIGTERM for graceful shutdown,
// this uses process.Kill() (TerminateProcess) which immediately terminates the
// process without cleanup. In-flight requests are dropped and deferred functions
// do not run. A future improvement could use a named pipe or event for graceful
// shutdown.
func StopProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("cannot find process: %w", err)
	}

	if err := process.Kill(); err != nil {
		return fmt.Errorf("cannot terminate process: %w", err)
	}

	_, _ = process.Wait()
	return nil
}
