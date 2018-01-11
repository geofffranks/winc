// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"code.cloudfoundry.org/winc/network"
	"github.com/Microsoft/hcsshim"
)

type HCSClient struct {
	GetHNSNetworkByNameStub        func(string) (*hcsshim.HNSNetwork, error)
	getHNSNetworkByNameMutex       sync.RWMutex
	getHNSNetworkByNameArgsForCall []struct {
		arg1 string
	}
	getHNSNetworkByNameReturns struct {
		result1 *hcsshim.HNSNetwork
		result2 error
	}
	getHNSNetworkByNameReturnsOnCall map[int]struct {
		result1 *hcsshim.HNSNetwork
		result2 error
	}
	CreateNetworkStub        func(*hcsshim.HNSNetwork, func() (bool, error)) (*hcsshim.HNSNetwork, error)
	createNetworkMutex       sync.RWMutex
	createNetworkArgsForCall []struct {
		arg1 *hcsshim.HNSNetwork
		arg2 func() (bool, error)
	}
	createNetworkReturns struct {
		result1 *hcsshim.HNSNetwork
		result2 error
	}
	createNetworkReturnsOnCall map[int]struct {
		result1 *hcsshim.HNSNetwork
		result2 error
	}
	DeleteNetworkStub        func(*hcsshim.HNSNetwork) (*hcsshim.HNSNetwork, error)
	deleteNetworkMutex       sync.RWMutex
	deleteNetworkArgsForCall []struct {
		arg1 *hcsshim.HNSNetwork
	}
	deleteNetworkReturns struct {
		result1 *hcsshim.HNSNetwork
		result2 error
	}
	deleteNetworkReturnsOnCall map[int]struct {
		result1 *hcsshim.HNSNetwork
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *HCSClient) GetHNSNetworkByName(arg1 string) (*hcsshim.HNSNetwork, error) {
	fake.getHNSNetworkByNameMutex.Lock()
	ret, specificReturn := fake.getHNSNetworkByNameReturnsOnCall[len(fake.getHNSNetworkByNameArgsForCall)]
	fake.getHNSNetworkByNameArgsForCall = append(fake.getHNSNetworkByNameArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetHNSNetworkByName", []interface{}{arg1})
	fake.getHNSNetworkByNameMutex.Unlock()
	if fake.GetHNSNetworkByNameStub != nil {
		return fake.GetHNSNetworkByNameStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getHNSNetworkByNameReturns.result1, fake.getHNSNetworkByNameReturns.result2
}

func (fake *HCSClient) GetHNSNetworkByNameCallCount() int {
	fake.getHNSNetworkByNameMutex.RLock()
	defer fake.getHNSNetworkByNameMutex.RUnlock()
	return len(fake.getHNSNetworkByNameArgsForCall)
}

func (fake *HCSClient) GetHNSNetworkByNameArgsForCall(i int) string {
	fake.getHNSNetworkByNameMutex.RLock()
	defer fake.getHNSNetworkByNameMutex.RUnlock()
	return fake.getHNSNetworkByNameArgsForCall[i].arg1
}

func (fake *HCSClient) GetHNSNetworkByNameReturns(result1 *hcsshim.HNSNetwork, result2 error) {
	fake.GetHNSNetworkByNameStub = nil
	fake.getHNSNetworkByNameReturns = struct {
		result1 *hcsshim.HNSNetwork
		result2 error
	}{result1, result2}
}

func (fake *HCSClient) GetHNSNetworkByNameReturnsOnCall(i int, result1 *hcsshim.HNSNetwork, result2 error) {
	fake.GetHNSNetworkByNameStub = nil
	if fake.getHNSNetworkByNameReturnsOnCall == nil {
		fake.getHNSNetworkByNameReturnsOnCall = make(map[int]struct {
			result1 *hcsshim.HNSNetwork
			result2 error
		})
	}
	fake.getHNSNetworkByNameReturnsOnCall[i] = struct {
		result1 *hcsshim.HNSNetwork
		result2 error
	}{result1, result2}
}

func (fake *HCSClient) CreateNetwork(arg1 *hcsshim.HNSNetwork, arg2 func() (bool, error)) (*hcsshim.HNSNetwork, error) {
	fake.createNetworkMutex.Lock()
	ret, specificReturn := fake.createNetworkReturnsOnCall[len(fake.createNetworkArgsForCall)]
	fake.createNetworkArgsForCall = append(fake.createNetworkArgsForCall, struct {
		arg1 *hcsshim.HNSNetwork
		arg2 func() (bool, error)
	}{arg1, arg2})
	fake.recordInvocation("CreateNetwork", []interface{}{arg1, arg2})
	fake.createNetworkMutex.Unlock()
	if fake.CreateNetworkStub != nil {
		return fake.CreateNetworkStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.createNetworkReturns.result1, fake.createNetworkReturns.result2
}

func (fake *HCSClient) CreateNetworkCallCount() int {
	fake.createNetworkMutex.RLock()
	defer fake.createNetworkMutex.RUnlock()
	return len(fake.createNetworkArgsForCall)
}

func (fake *HCSClient) CreateNetworkArgsForCall(i int) (*hcsshim.HNSNetwork, func() (bool, error)) {
	fake.createNetworkMutex.RLock()
	defer fake.createNetworkMutex.RUnlock()
	return fake.createNetworkArgsForCall[i].arg1, fake.createNetworkArgsForCall[i].arg2
}

func (fake *HCSClient) CreateNetworkReturns(result1 *hcsshim.HNSNetwork, result2 error) {
	fake.CreateNetworkStub = nil
	fake.createNetworkReturns = struct {
		result1 *hcsshim.HNSNetwork
		result2 error
	}{result1, result2}
}

func (fake *HCSClient) CreateNetworkReturnsOnCall(i int, result1 *hcsshim.HNSNetwork, result2 error) {
	fake.CreateNetworkStub = nil
	if fake.createNetworkReturnsOnCall == nil {
		fake.createNetworkReturnsOnCall = make(map[int]struct {
			result1 *hcsshim.HNSNetwork
			result2 error
		})
	}
	fake.createNetworkReturnsOnCall[i] = struct {
		result1 *hcsshim.HNSNetwork
		result2 error
	}{result1, result2}
}

func (fake *HCSClient) DeleteNetwork(arg1 *hcsshim.HNSNetwork) (*hcsshim.HNSNetwork, error) {
	fake.deleteNetworkMutex.Lock()
	ret, specificReturn := fake.deleteNetworkReturnsOnCall[len(fake.deleteNetworkArgsForCall)]
	fake.deleteNetworkArgsForCall = append(fake.deleteNetworkArgsForCall, struct {
		arg1 *hcsshim.HNSNetwork
	}{arg1})
	fake.recordInvocation("DeleteNetwork", []interface{}{arg1})
	fake.deleteNetworkMutex.Unlock()
	if fake.DeleteNetworkStub != nil {
		return fake.DeleteNetworkStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.deleteNetworkReturns.result1, fake.deleteNetworkReturns.result2
}

func (fake *HCSClient) DeleteNetworkCallCount() int {
	fake.deleteNetworkMutex.RLock()
	defer fake.deleteNetworkMutex.RUnlock()
	return len(fake.deleteNetworkArgsForCall)
}

func (fake *HCSClient) DeleteNetworkArgsForCall(i int) *hcsshim.HNSNetwork {
	fake.deleteNetworkMutex.RLock()
	defer fake.deleteNetworkMutex.RUnlock()
	return fake.deleteNetworkArgsForCall[i].arg1
}

func (fake *HCSClient) DeleteNetworkReturns(result1 *hcsshim.HNSNetwork, result2 error) {
	fake.DeleteNetworkStub = nil
	fake.deleteNetworkReturns = struct {
		result1 *hcsshim.HNSNetwork
		result2 error
	}{result1, result2}
}

func (fake *HCSClient) DeleteNetworkReturnsOnCall(i int, result1 *hcsshim.HNSNetwork, result2 error) {
	fake.DeleteNetworkStub = nil
	if fake.deleteNetworkReturnsOnCall == nil {
		fake.deleteNetworkReturnsOnCall = make(map[int]struct {
			result1 *hcsshim.HNSNetwork
			result2 error
		})
	}
	fake.deleteNetworkReturnsOnCall[i] = struct {
		result1 *hcsshim.HNSNetwork
		result2 error
	}{result1, result2}
}

func (fake *HCSClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getHNSNetworkByNameMutex.RLock()
	defer fake.getHNSNetworkByNameMutex.RUnlock()
	fake.createNetworkMutex.RLock()
	defer fake.createNetworkMutex.RUnlock()
	fake.deleteNetworkMutex.RLock()
	defer fake.deleteNetworkMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *HCSClient) recordInvocation(key string, args []interface{}) {
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

var _ network.HCSClient = new(HCSClient)
