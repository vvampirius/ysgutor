package ysgutor

import (
	"fmt"
	"github.com/mitchellh/go-ps"
	"syscall"
	"time"
	"os"
)

func TerminateByPid(pid int, recursive bool, timeout time.Duration, q ...interface{}) bool {
	timeoutChan := time.After(timeout)
	childs := make(map[int]bool)
	if recursive {
		appendChilds(pid, childs)
	}
	exited := false
	exitChan := make(chan bool, 1)
	go PidWait(pid, &exited, exitChan)
	for _, v := range q {
		if s, ok := v.(syscall.Signal); ok {
			//fmt.Printf("Sending signal %d to pid %d\n", s, pid)
			syscall.Kill(pid, s)
		} else if wait, ok := v.(time.Duration); ok {
			select {
				case <-exitChan:
				case <-time.After(wait):
			}
		} else {
			fmt.Print("Error value: ")
			fmt.Println(v)
		}
		if exited {
			break
		}
	}
	if !exited {
		select {
			case <-exitChan:
			case <-timeoutChan:
		}
	}
	if recursive {
		appendChilds(pid, childs)
		for p, _ := range childs {
			TerminateByPid(p, recursive, timeout, q...)
		}
	}
	return exited
}

func appendChilds(parentPid int, childs map[int]bool) {
	if pp, err := ps.Processes(); err==nil {
		for _, p := range pp {
			ppid := p.PPid()
			if parentPid == ppid {
				childs[p.Pid()] = false
			}
		}
	}
}

func PidWait(pid int, exited *bool, exitChan chan<- bool) {
	for {
		if p, err := os.FindProcess(pid); err == nil {
			if errr := p.Signal(syscall.Signal(0)); errr!=nil {
				break
			}
			//fmt.Printf("process.Signal on pid %d returned: %v\n", pid, errr)
		} else {
			fmt.Println(err.Error())
			break
		}
	}
	*exited = true
	exitChan <- true
}