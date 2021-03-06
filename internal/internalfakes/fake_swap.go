// Code generated by counterfeiter. DO NOT EDIT.
package internalfakes

import (
	"os/exec"
	"sync"

	procswap "github.com/billiford/procswap/internal"
)

type FakeSwap struct {
	CmdStub        func() *exec.Cmd
	cmdMutex       sync.RWMutex
	cmdArgsForCall []struct {
	}
	cmdReturns struct {
		result1 *exec.Cmd
	}
	cmdReturnsOnCall map[int]struct {
		result1 *exec.Cmd
	}
	KillStub        func() error
	killMutex       sync.RWMutex
	killArgsForCall []struct {
	}
	killReturns struct {
		result1 error
	}
	killReturnsOnCall map[int]struct {
		result1 error
	}
	PIDStub        func() int
	pIDMutex       sync.RWMutex
	pIDArgsForCall []struct {
	}
	pIDReturns struct {
		result1 int
	}
	pIDReturnsOnCall map[int]struct {
		result1 int
	}
	PathStub        func() string
	pathMutex       sync.RWMutex
	pathArgsForCall []struct {
	}
	pathReturns struct {
		result1 string
	}
	pathReturnsOnCall map[int]struct {
		result1 string
	}
	ShowOutputStub        func(bool)
	showOutputMutex       sync.RWMutex
	showOutputArgsForCall []struct {
		arg1 bool
	}
	StartStub        func() error
	startMutex       sync.RWMutex
	startArgsForCall []struct {
	}
	startReturns struct {
		result1 error
	}
	startReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeSwap) Cmd() *exec.Cmd {
	fake.cmdMutex.Lock()
	ret, specificReturn := fake.cmdReturnsOnCall[len(fake.cmdArgsForCall)]
	fake.cmdArgsForCall = append(fake.cmdArgsForCall, struct {
	}{})
	stub := fake.CmdStub
	fakeReturns := fake.cmdReturns
	fake.recordInvocation("Cmd", []interface{}{})
	fake.cmdMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeSwap) CmdCallCount() int {
	fake.cmdMutex.RLock()
	defer fake.cmdMutex.RUnlock()
	return len(fake.cmdArgsForCall)
}

func (fake *FakeSwap) CmdCalls(stub func() *exec.Cmd) {
	fake.cmdMutex.Lock()
	defer fake.cmdMutex.Unlock()
	fake.CmdStub = stub
}

func (fake *FakeSwap) CmdReturns(result1 *exec.Cmd) {
	fake.cmdMutex.Lock()
	defer fake.cmdMutex.Unlock()
	fake.CmdStub = nil
	fake.cmdReturns = struct {
		result1 *exec.Cmd
	}{result1}
}

func (fake *FakeSwap) CmdReturnsOnCall(i int, result1 *exec.Cmd) {
	fake.cmdMutex.Lock()
	defer fake.cmdMutex.Unlock()
	fake.CmdStub = nil
	if fake.cmdReturnsOnCall == nil {
		fake.cmdReturnsOnCall = make(map[int]struct {
			result1 *exec.Cmd
		})
	}
	fake.cmdReturnsOnCall[i] = struct {
		result1 *exec.Cmd
	}{result1}
}

func (fake *FakeSwap) Kill() error {
	fake.killMutex.Lock()
	ret, specificReturn := fake.killReturnsOnCall[len(fake.killArgsForCall)]
	fake.killArgsForCall = append(fake.killArgsForCall, struct {
	}{})
	stub := fake.KillStub
	fakeReturns := fake.killReturns
	fake.recordInvocation("Kill", []interface{}{})
	fake.killMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeSwap) KillCallCount() int {
	fake.killMutex.RLock()
	defer fake.killMutex.RUnlock()
	return len(fake.killArgsForCall)
}

func (fake *FakeSwap) KillCalls(stub func() error) {
	fake.killMutex.Lock()
	defer fake.killMutex.Unlock()
	fake.KillStub = stub
}

func (fake *FakeSwap) KillReturns(result1 error) {
	fake.killMutex.Lock()
	defer fake.killMutex.Unlock()
	fake.KillStub = nil
	fake.killReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeSwap) KillReturnsOnCall(i int, result1 error) {
	fake.killMutex.Lock()
	defer fake.killMutex.Unlock()
	fake.KillStub = nil
	if fake.killReturnsOnCall == nil {
		fake.killReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.killReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeSwap) PID() int {
	fake.pIDMutex.Lock()
	ret, specificReturn := fake.pIDReturnsOnCall[len(fake.pIDArgsForCall)]
	fake.pIDArgsForCall = append(fake.pIDArgsForCall, struct {
	}{})
	stub := fake.PIDStub
	fakeReturns := fake.pIDReturns
	fake.recordInvocation("PID", []interface{}{})
	fake.pIDMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeSwap) PIDCallCount() int {
	fake.pIDMutex.RLock()
	defer fake.pIDMutex.RUnlock()
	return len(fake.pIDArgsForCall)
}

func (fake *FakeSwap) PIDCalls(stub func() int) {
	fake.pIDMutex.Lock()
	defer fake.pIDMutex.Unlock()
	fake.PIDStub = stub
}

func (fake *FakeSwap) PIDReturns(result1 int) {
	fake.pIDMutex.Lock()
	defer fake.pIDMutex.Unlock()
	fake.PIDStub = nil
	fake.pIDReturns = struct {
		result1 int
	}{result1}
}

func (fake *FakeSwap) PIDReturnsOnCall(i int, result1 int) {
	fake.pIDMutex.Lock()
	defer fake.pIDMutex.Unlock()
	fake.PIDStub = nil
	if fake.pIDReturnsOnCall == nil {
		fake.pIDReturnsOnCall = make(map[int]struct {
			result1 int
		})
	}
	fake.pIDReturnsOnCall[i] = struct {
		result1 int
	}{result1}
}

func (fake *FakeSwap) Path() string {
	fake.pathMutex.Lock()
	ret, specificReturn := fake.pathReturnsOnCall[len(fake.pathArgsForCall)]
	fake.pathArgsForCall = append(fake.pathArgsForCall, struct {
	}{})
	stub := fake.PathStub
	fakeReturns := fake.pathReturns
	fake.recordInvocation("Path", []interface{}{})
	fake.pathMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeSwap) PathCallCount() int {
	fake.pathMutex.RLock()
	defer fake.pathMutex.RUnlock()
	return len(fake.pathArgsForCall)
}

func (fake *FakeSwap) PathCalls(stub func() string) {
	fake.pathMutex.Lock()
	defer fake.pathMutex.Unlock()
	fake.PathStub = stub
}

func (fake *FakeSwap) PathReturns(result1 string) {
	fake.pathMutex.Lock()
	defer fake.pathMutex.Unlock()
	fake.PathStub = nil
	fake.pathReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeSwap) PathReturnsOnCall(i int, result1 string) {
	fake.pathMutex.Lock()
	defer fake.pathMutex.Unlock()
	fake.PathStub = nil
	if fake.pathReturnsOnCall == nil {
		fake.pathReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.pathReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeSwap) ShowOutput(arg1 bool) {
	fake.showOutputMutex.Lock()
	fake.showOutputArgsForCall = append(fake.showOutputArgsForCall, struct {
		arg1 bool
	}{arg1})
	stub := fake.ShowOutputStub
	fake.recordInvocation("ShowOutput", []interface{}{arg1})
	fake.showOutputMutex.Unlock()
	if stub != nil {
		fake.ShowOutputStub(arg1)
	}
}

func (fake *FakeSwap) ShowOutputCallCount() int {
	fake.showOutputMutex.RLock()
	defer fake.showOutputMutex.RUnlock()
	return len(fake.showOutputArgsForCall)
}

func (fake *FakeSwap) ShowOutputCalls(stub func(bool)) {
	fake.showOutputMutex.Lock()
	defer fake.showOutputMutex.Unlock()
	fake.ShowOutputStub = stub
}

func (fake *FakeSwap) ShowOutputArgsForCall(i int) bool {
	fake.showOutputMutex.RLock()
	defer fake.showOutputMutex.RUnlock()
	argsForCall := fake.showOutputArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeSwap) Start() error {
	fake.startMutex.Lock()
	ret, specificReturn := fake.startReturnsOnCall[len(fake.startArgsForCall)]
	fake.startArgsForCall = append(fake.startArgsForCall, struct {
	}{})
	stub := fake.StartStub
	fakeReturns := fake.startReturns
	fake.recordInvocation("Start", []interface{}{})
	fake.startMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeSwap) StartCallCount() int {
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	return len(fake.startArgsForCall)
}

func (fake *FakeSwap) StartCalls(stub func() error) {
	fake.startMutex.Lock()
	defer fake.startMutex.Unlock()
	fake.StartStub = stub
}

func (fake *FakeSwap) StartReturns(result1 error) {
	fake.startMutex.Lock()
	defer fake.startMutex.Unlock()
	fake.StartStub = nil
	fake.startReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeSwap) StartReturnsOnCall(i int, result1 error) {
	fake.startMutex.Lock()
	defer fake.startMutex.Unlock()
	fake.StartStub = nil
	if fake.startReturnsOnCall == nil {
		fake.startReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.startReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeSwap) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.cmdMutex.RLock()
	defer fake.cmdMutex.RUnlock()
	fake.killMutex.RLock()
	defer fake.killMutex.RUnlock()
	fake.pIDMutex.RLock()
	defer fake.pIDMutex.RUnlock()
	fake.pathMutex.RLock()
	defer fake.pathMutex.RUnlock()
	fake.showOutputMutex.RLock()
	defer fake.showOutputMutex.RUnlock()
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeSwap) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ procswap.Swap = new(FakeSwap)
