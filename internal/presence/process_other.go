//go:build !unix && !windows

package presence

func probeDaemonProcess(_ int, _ daemonProcessFiles) string {
	return DaemonUnknown
}
