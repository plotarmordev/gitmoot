//go:build unix

package presence

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func probeDaemonProcess(pid int, files daemonProcessFiles) string {
	state := probeProcess(pid)
	if state != DaemonRunning {
		return state
	}
	if processLooksLikeDaemon(pid, files.MetaFile) {
		return DaemonRunning
	}
	return DaemonStopped
}

func probeProcess(pid int) string {
	err := syscall.Kill(pid, 0)
	if err == nil || errors.Is(err, syscall.EPERM) {
		return DaemonRunning
	}
	if errors.Is(err, syscall.ESRCH) {
		return DaemonStopped
	}
	return DaemonUnknown
}

func processLooksLikeDaemon(pid int, metaFile string) bool {
	contents, err := os.ReadFile(metaFile)
	if err != nil {
		return false
	}
	var meta daemonMeta
	if err := json.Unmarshal(contents, &meta); err != nil {
		return false
	}
	if meta.PID != pid {
		return false
	}
	if processCmdlineLooksLikeDaemon(pid, meta) {
		return true
	}
	return processPSLooksLikeDaemon(pid, meta)
}

func processCmdlineLooksLikeDaemon(pid int, meta daemonMeta) bool {
	cmdline, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "cmdline"))
	if err != nil {
		return false
	}
	parts := strings.Split(strings.TrimRight(string(cmdline), "\x00"), "\x00")
	return daemonProcessArgsMatch(parts, meta)
}

func processPSLooksLikeDaemon(pid int, meta daemonMeta) bool {
	if hasWhitespace(meta.Executable) {
		return false
	}
	for _, arg := range meta.Args {
		if hasWhitespace(arg) {
			return false
		}
	}
	result, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "command=").Output()
	if err != nil {
		return false
	}
	command := strings.TrimSpace(string(result))
	if command == "" {
		return false
	}
	return daemonProcessArgsMatch(strings.Fields(command), meta)
}

func daemonProcessArgsMatch(argv []string, meta daemonMeta) bool {
	if len(argv) != len(meta.Args)+1 {
		return false
	}
	if meta.Executable != "" && argv[0] != meta.Executable {
		return false
	}
	return equalStringSlices(argv[1:], meta.Args)
}

func equalStringSlices(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

func hasWhitespace(value string) bool {
	return strings.ContainsAny(value, " \t\r\n")
}
