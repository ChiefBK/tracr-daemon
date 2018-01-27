package util

import (
	"os/exec"
	log "github.com/inconshreveable/log15"
	"errors"
	"syscall"
)

func ExecCommand(command *exec.Cmd) error {
	startErr := command.Start()
	if startErr != nil {
		log.Error("Error starting command", "module", "main", "error", startErr)
		return errors.New("error starting command")
	}

	waitErr := command.Wait()
	if waitErr != nil {
		log.Error("Error waiting for command", "module", "main", "error", waitErr)
		return errors.New("error waiting for command")
	}

	return nil
}

func ExecCommandWithCode(command *exec.Cmd) int {
	if err := command.Start(); err != nil {
		log.Error("Error starting command (with code)", "module", "main", "error", err)
		return -1
	}

	if err := command.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return int(status)
			}
		} else {
			log.Error("Error waiting for command (with code)", "module", "main", "error", err)
			return -1
		}
	}

	return 0
}
