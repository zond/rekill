package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func rekill(pid int) error {
	cmd := []string{"pgrep", "-P", fmt.Sprint(pid)}
	output, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		if !strings.Contains(err.Error(), "exit status 1") {
			return fmt.Errorf("Error running %#v: %v", cmd, err)
		}
	}
	for _, child := range strings.Split(string(output), "\n") {
		if child = strings.TrimSpace(child); child != "" {
			childPid := 0
			if childPid, err = strconv.Atoi(strings.TrimSpace(child)); err != nil {
				return fmt.Errorf("Bad child pid %#v: %v", child, err)
			}
			if err = rekill(childPid); err != nil {
				return err
			}
		}
	}
	cmd = []string{"kill", "-9", fmt.Sprint(pid)}
	if err = exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
		if !strings.Contains(err.Error(), "exit status 1") {
			return fmt.Errorf("Error running %#v: %v", cmd, err)
		}
	}
	fmt.Printf("Killed %#v\n", pid)
	return nil
}

func main() {
	rootpid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Usage: rekill [pid] (%v)", err)
		return
	}
	if err = rekill(rootpid); err != nil {
		fmt.Println(err)
	}
}
