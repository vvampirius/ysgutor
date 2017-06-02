package ysgutor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)

const VERSION float32 = 0.1

type PreStartHandlerFunc func(*Ysgutor) bool
type StartHandlerFunc func(*Ysgutor)
type ExitHandlerFunc func(*Ysgutor, error)
type KillHandlerFunc func(*Ysgutor) bool

type Ysgutor struct {
	Identifier     interface{}
	CommandPath    string
	CommandArgs    []string
	Context        context.Context
	Env            []string
	StdInReader    io.Reader
	StdOutWriter   io.Writer
	StdErrWriter   io.Writer
	OnPreStart     PreStartHandlerFunc
	OnStart        StartHandlerFunc
	OnExit         ExitHandlerFunc
	OnFail         ExitHandlerFunc
	KillHandler    KillHandlerFunc
	StartedAt      time.Time
	ExitedAt       time.Time
	Terminated     bool
	executeMutex   sync.Mutex
	terminateMutex sync.Mutex
	Cmd            *exec.Cmd
}

func (self *Ysgutor) Execute() {
	//TODO: wipe vars if reexecuted
	self.executeMutex.Lock()
	defer self.executeMutex.Unlock()
	self.Cmd = exec.Command(self.CommandPath, self.CommandArgs...)
	self.Cmd.Env = self.Env
	self.Cmd.Stdin = self.StdInReader
	self.Cmd.Stdout = self.StdOutWriter
	self.Cmd.Stderr = self.StdErrWriter
	startAllowed := true
	if self.OnPreStart != nil {
		startAllowed = self.OnPreStart(self)
	}
	if startAllowed {
		if startErr := self.Cmd.Start(); startErr == nil {
			self.StartedAt = time.Now()
			var cancel context.CancelFunc
			if self.Context != nil {
				ctx, _cncl := context.WithCancel(self.Context)
				go self.contextTerminator(ctx)
				cancel = _cncl
			}
			if self.OnStart != nil {
				self.OnStart(self)
			}
			waitErr := self.Cmd.Wait()
			self.ExitedAt = time.Now()
			if cancel != nil {
				cancel()
			}
			if self.OnExit != nil {
				self.OnExit(self, waitErr)
			}
		} else {
			if self.OnFail != nil {
				self.OnFail(self, startErr)
			}
		}
	}
}

func (self *Ysgutor) Terminate() {
	self.terminateMutex.Lock()
	defer self.terminateMutex.Unlock()
	if !self.Terminated {
		if self.KillHandler != nil {
			self.Terminated = self.KillHandler(self)
		} else {
			if self.Cmd != nil && self.Cmd.ProcessState == nil && self.Cmd.Process != nil {
				self.Terminated = true
			} else {
			}
		}
	}
}

func (self *Ysgutor) contextTerminator(ctx context.Context) {
	if ctx != nil {
		<-ctx.Done()
		self.Terminate()
	}
}

func New(command interface{}, shell interface{}) (Ysgutor, error) {
	commandPath := ""
	commandArgs := make([]string, 0)
	if commandArray, ok := command.([]string); ok {
		if len(commandArray) > 0 {
			commandPath = commandArray[0]
			commandArgs = commandArray[1:]
		} else {
			return Ysgutor{}, errors.New("No command to execute in the array!")
		}
	} else if commandString, ok := command.(string); ok {
		if len(commandString) > 0 {
			if shellPlaceholder, ok := shell.(string); ok {
				commandInShell := fmt.Sprintf(shellPlaceholder, commandString)
				commandInShellParsed := ParseCommandLine(commandInShell)
				commandPath = commandInShellParsed[0]
				commandArgs = commandInShellParsed[1:]
			} else {
				commandParsed := ParseCommandLine(commandString)
				commandPath = commandParsed[0]
				commandArgs = commandParsed[1:]
			}
			if c, err := exec.LookPath(commandPath); err == nil {
				commandPath = c
			} else {
				return Ysgutor{}, err
			}
		} else {
			return Ysgutor{}, errors.New("No command to execute in the string!")
		}
	} else {
		return Ysgutor{}, errors.New("Unknown command type!")
	}
	ysgutor := Ysgutor{
		CommandPath: commandPath,
		CommandArgs: commandArgs,
	}
	return ysgutor, nil
}
