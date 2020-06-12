/*
Copyright 2017 The Kubernetes Authors.
Copyright 2018 Intel Coporation.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runCommand(cmd string, args ...string) (string, error) {
	output, err := exec.Command(cmd, args...).CombinedOutput()
	handleError(err)
	strOutput := string(output)
	return strOutput, err
}

func main() {
	/* DEBUG: show mounts
	out, _ := runCommand("/usr/bin/mount")
	fmt.Printf("mount_cmd_out:\n%s\n", out)*/
	// Start trace watch
	runCommand("echo", "1", ">", "/sys/kernel/debug/tracing/events/fs_dax/dax_pmd_fault_done/enable")
	runCommand("echo", "1", ">", "/sys/kernel/debug/tracing/tracing_on")
	watchcmd := &exec.Cmd{
		Path:   "/usr/bin/cat",
		Args:   []string{"/usr/bin/cat", "/sys/kernel/debug/tracing/trace_pipe"},
		Stdout: &b,
		//Stderr: os.Stdout,
	}
	err := watchcmd.Start()
	handleError(err)
	// Cmd remains running, will be terminated from outside.
}
