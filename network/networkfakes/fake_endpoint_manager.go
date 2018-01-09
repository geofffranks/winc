// Code generated by counterfeiter. DO NOT EDIT.
package networkfakes

import (
	"sync"

	"code.cloudfoundry.org/winc/network"
	"code.cloudfoundry.org/winc/network/netrules"
	"github.com/Microsoft/hcsshim"
)

type FakeEndpointManager struct {
	CreateStub        func() (hcsshim.HNSEndpoint, error)
	createMutex       sync.RWMutex
	createArgsForCall []struct{}
	createReturns     struct {
		result1 hcsshim.HNSEndpoint
		result2 error
	}
	createReturnsOnCall map[int]struct {
		result1 hcsshim.HNSEndpoint
		result2 error
	}
	DeleteStub        func() error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct{}
	deleteReturns     struct {
		result1 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 error
	}
	ApplyMappingsStub        func(hcsshim.HNSEndpoint, []netrules.PortMapping) (hcsshim.HNSEndpoint, error)
	applyMappingsMutex       sync.RWMutex
	applyMappingsArgsForCall []struct {
		arg1 hcsshim.HNSEndpoint
		arg2 []netrules.PortMapping
	}
	applyMappingsReturns struct {
		result1 hcsshim.HNSEndpoint
		result2 error
	}
	applyMappingsReturnsOnCall map[int]struct {
		result1 hcsshim.HNSEndpoint
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeEndpointManager) Create() (hcsshim.HNSEndpoint, error) {
	fake.createMutex.Lock()
	ret, specificReturn := fake.createReturnsOnCall[len(fake.createArgsForCall)]
	fake.createArgsForCall = append(fake.createArgsForCall, struct{}{})
	fake.recordInvocation("Create", []interface{}{})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.createReturns.result1, fake.createReturns.result2
}

func (fake *FakeEndpointManager) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeEndpointManager) CreateReturns(result1 hcsshim.HNSEndpoint, result2 error) {
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 hcsshim.HNSEndpoint
		result2 error
	}{result1, result2}
}

func (fake *FakeEndpointManager) CreateReturnsOnCall(i int, result1 hcsshim.HNSEndpoint, result2 error) {
	fake.CreateStub = nil
	if fake.createReturnsOnCall == nil {
		fake.createReturnsOnCall = make(map[int]struct {
			result1 hcsshim.HNSEndpoint
			result2 error
		})
	}
	fake.createReturnsOnCall[i] = struct {
		result1 hcsshim.HNSEndpoint
		result2 error
	}{result1, result2}
}

func (fake *FakeEndpointManager) Delete() error {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct{}{})
	fake.recordInvocation("Delete", []interface{}{})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.deleteReturns.result1
}

func (fake *FakeEndpointManager) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeEndpointManager) DeleteReturns(result1 error) {
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeEndpointManager) DeleteReturnsOnCall(i int, result1 error) {
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeEndpointManager) ApplyMappings(arg1 hcsshim.HNSEndpoint, arg2 []netrules.PortMapping) (hcsshim.HNSEndpoint, error) {
	var arg2Copy []netrules.PortMapping
	if arg2 != nil {
		arg2Copy = make([]netrules.PortMapping, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.applyMappingsMutex.Lock()
	ret, specificReturn := fake.applyMappingsReturnsOnCall[len(fake.applyMappingsArgsForCall)]
	fake.applyMappingsArgsForCall = append(fake.applyMappingsArgsForCall, struct {
		arg1 hcsshim.HNSEndpoint
		arg2 []netrules.PortMapping
	}{arg1, arg2Copy})
	fake.recordInvocation("ApplyMappings", []interface{}{arg1, arg2Copy})
	fake.applyMappingsMutex.Unlock()
	if fake.ApplyMappingsStub != nil {
		return fake.ApplyMappingsStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.applyMappingsReturns.result1, fake.applyMappingsReturns.result2
}

func (fake *FakeEndpointManager) ApplyMappingsCallCount() int {
	fake.applyMappingsMutex.RLock()
	defer fake.applyMappingsMutex.RUnlock()
	return len(fake.applyMappingsArgsForCall)
}

func (fake *FakeEndpointManager) ApplyMappingsArgsForCall(i int) (hcsshim.HNSEndpoint, []netrules.PortMapping) {
	fake.applyMappingsMutex.RLock()
	defer fake.applyMappingsMutex.RUnlock()
	return fake.applyMappingsArgsForCall[i].arg1, fake.applyMappingsArgsForCall[i].arg2
}

func (fake *FakeEndpointManager) ApplyMappingsReturns(result1 hcsshim.HNSEndpoint, result2 error) {
	fake.ApplyMappingsStub = nil
	fake.applyMappingsReturns = struct {
		result1 hcsshim.HNSEndpoint
		result2 error
	}{result1, result2}
}

func (fake *FakeEndpointManager) ApplyMappingsReturnsOnCall(i int, result1 hcsshim.HNSEndpoint, result2 error) {
	fake.ApplyMappingsStub = nil
	if fake.applyMappingsReturnsOnCall == nil {
		fake.applyMappingsReturnsOnCall = make(map[int]struct {
			result1 hcsshim.HNSEndpoint
			result2 error
		})
	}
	fake.applyMappingsReturnsOnCall[i] = struct {
		result1 hcsshim.HNSEndpoint
		result2 error
	}{result1, result2}
}

func (fake *FakeEndpointManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.applyMappingsMutex.RLock()
	defer fake.applyMappingsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeEndpointManager) recordInvocation(key string, args []interface{}) {
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

var _ network.EndpointManager = new(FakeEndpointManager)
