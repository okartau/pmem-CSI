/*
Copyright 2017 The Kubernetes Authors.
Copyright 2018 Intel Coporation.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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
	var num_pagefaults = 0
	var num_fallbacks = 0
	if len(os.Args) != 3 {
		fmt.Printf("Need 2 args: mapped_addr inode\n")
		os.Exit(2)
	}
	//fmt.Printf("len of args: %d\n", len(os.Args))
	str_mapped := os.Args[1]
	str_inode := os.Args[2]
	dat, err := ioutil.ReadFile("/tmp/hugepagedat")
	handleError(err)

	tracedStr := string(dat)
	fmt.Println(tracedStr)
	lines := strings.Split(tracedStr, "\n")
	for _, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) != 22 {
			continue
		}
		traced_inode := fields[8]
		traced_vm_start := fields[14]
		pagingtype := fields[21]
		/*pidstr := fields[0]
		  tracedPid := -1
		  _, err := fmt.Sscanf(pidstr, "<...>-%d", &tracedPid)
		  handleError(err)*/
		/*fmt.Printf("inode:%s vm_start:%s vm_end:%s ptype:%s\n",
		  pidstr, tracedPid, traced_inode, traced_vm_start, pagingtype)*/
		if str_mapped == traced_vm_start {
			//fmt.Printf("vm_start matches\n")
		} else {
			fmt.Printf("vm_start mismatch\n")
			continue
		}
		if str_inode == traced_inode {
			//fmt.Printf("inode matches\n")
		} else {
			fmt.Printf("inode mismatch\n")
			continue
		}
		/* Skip PID check, as some entries show off-by-few PID for some reason
		   if runPid == tracedPid {
		           fmt.Printf("pid matches\n")
		   } else {
		           continue
		   }*/
		//fmt.Printf("All checks match, this is our mapping\n")
		if pagingtype == "NOPAGE" {
			num_pagefaults += 1
		} else if pagingtype == "FALLBACK" {
			num_fallbacks += 1
		}

	}

	fmt.Printf("Page faults:%d Fallbacks:%d\n", num_pagefaults, num_fallbacks)
	if num_pagefaults > 0 {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
