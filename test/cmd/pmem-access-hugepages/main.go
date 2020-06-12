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
	"syscall"
	"unsafe"
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
	const fname = "/mnt/data2"
	var stat syscall.Stat_t
	const size = 2*1024*1024 + 4
	//const size = 16*1024*1024 + 4
	runPid := os.Getpid()
	//fmt.Printf("runPid: %d\n", runPid)
	// Create the file
	map_file, err := os.Create(fname)
	handleError(err)
	_, err = map_file.Seek(int64(size-1), 0)
	handleError(err)
	_, err = map_file.Write([]byte(" "))
	handleError(err)
	// Get inode number of the file
	err = syscall.Stat(fname, &stat)
	handleError(err)
	//str_inode := fmt.Sprintf("0x%x", stat.Ino)

	mmap, err := syscall.Mmap(int(map_file.Fd()), 0, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	//mmap, err := syscall.Mmap(int(map_file.Fd()), 0, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE)
	handleError(err)
	mapped := (*[size]byte)(unsafe.Pointer(&mmap[0]))
	//str_mapped := fmt.Sprintf("%p", mapped)
	for i := 1; i < size; i++ {
		mapped[i] = byte(runPid)
	}

	//fmt.Println(*mapped)
	err = syscall.Munmap(mmap)
	handleError(err)
	err = map_file.Close()
	handleError(err)
	fmt.Printf("0x%x%p", stat.Ino, mapped)
}
