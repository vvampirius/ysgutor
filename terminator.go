package ysgutor

import (
	"fmt"
	"github.com/mitchellh/go-ps"
	"syscall"
	"time"
	"os"
	"context"
)

// q will be syscall.Signal to send signal to pid or time.Duration to wait
func TerminateByPid(pid int, recursive bool, timeout context.Context, q ...interface{}) bool {
	childs := make(map[int]bool)
	if recursive {
		appendChilds(pid, childs)
	}
	exited := false
	pidWaitContext, pidWaitCancel := context.WithCancel(context.Background())
	go PidWait(pid, &exited, timeout, pidWaitCancel)
	for _, v := range q {
		breakCycle := false
		if s, ok := v.(syscall.Signal); ok {
			//fmt.Printf("Sending signal %d to pid %d\n", s, pid)
			syscall.Kill(pid, s)
		} else if wait, ok := v.(time.Duration); ok {
			select {
				case <-pidWaitContext.Done():
				case <-time.After(wait):
				case <-timeout.Done():
					breakCycle = true
			}
		} else {
			fmt.Print("Error value: ")
			fmt.Println(v)
		}
		if exited || breakCycle {
			break
		}
	}
	if !exited {
		select {
			case <-pidWaitContext.Done():
			case <-timeout.Done():
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

// PidWait check if pid exited, sets exited to true (if exited) and call pidWaitCancel().
//
// It use os.FindProcess(pid).Signal(syscall.Signal(0)) to check if process alive. Because with syscall.Wait4() can't
// check subprocesses of subprocesses.
func PidWait(pid int, exited *bool, timeout context.Context, pidWaitCancel context.CancelFunc) {
	cycle := true
	for cycle {
		select {
		default:
			if p, err := os.FindProcess(pid); err == nil {
				if errr := p.Signal(syscall.Signal(0)); errr != nil {
					cycle = false
					*exited = true
					pidWaitCancel()
				}
				//fmt.Printf("process.Signal on pid %d returned: %v\n", pid, errr)
			} else {
				fmt.Println(err.Error())
				cycle = false
				*exited = true
				pidWaitCancel()
			}
		case <-timeout.Done():
			cycle = false
		}
	}
}