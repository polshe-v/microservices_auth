// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

package mocks

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"github.com/polshe-v/microservices_auth/internal/model"
)

// AccessRepositoryMock implements repository.AccessRepository
type AccessRepositoryMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcGetRoleEndpoints          func(ctx context.Context) (epa1 []*model.EndpointPermissions, err error)
	inspectFuncGetRoleEndpoints   func(ctx context.Context)
	afterGetRoleEndpointsCounter  uint64
	beforeGetRoleEndpointsCounter uint64
	GetRoleEndpointsMock          mAccessRepositoryMockGetRoleEndpoints
}

// NewAccessRepositoryMock returns a mock for repository.AccessRepository
func NewAccessRepositoryMock(t minimock.Tester) *AccessRepositoryMock {
	m := &AccessRepositoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetRoleEndpointsMock = mAccessRepositoryMockGetRoleEndpoints{mock: m}
	m.GetRoleEndpointsMock.callArgs = []*AccessRepositoryMockGetRoleEndpointsParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mAccessRepositoryMockGetRoleEndpoints struct {
	mock               *AccessRepositoryMock
	defaultExpectation *AccessRepositoryMockGetRoleEndpointsExpectation
	expectations       []*AccessRepositoryMockGetRoleEndpointsExpectation

	callArgs []*AccessRepositoryMockGetRoleEndpointsParams
	mutex    sync.RWMutex
}

// AccessRepositoryMockGetRoleEndpointsExpectation specifies expectation struct of the AccessRepository.GetRoleEndpoints
type AccessRepositoryMockGetRoleEndpointsExpectation struct {
	mock    *AccessRepositoryMock
	params  *AccessRepositoryMockGetRoleEndpointsParams
	results *AccessRepositoryMockGetRoleEndpointsResults
	Counter uint64
}

// AccessRepositoryMockGetRoleEndpointsParams contains parameters of the AccessRepository.GetRoleEndpoints
type AccessRepositoryMockGetRoleEndpointsParams struct {
	ctx context.Context
}

// AccessRepositoryMockGetRoleEndpointsResults contains results of the AccessRepository.GetRoleEndpoints
type AccessRepositoryMockGetRoleEndpointsResults struct {
	epa1 []*model.EndpointPermissions
	err  error
}

// Expect sets up expected params for AccessRepository.GetRoleEndpoints
func (mmGetRoleEndpoints *mAccessRepositoryMockGetRoleEndpoints) Expect(ctx context.Context) *mAccessRepositoryMockGetRoleEndpoints {
	if mmGetRoleEndpoints.mock.funcGetRoleEndpoints != nil {
		mmGetRoleEndpoints.mock.t.Fatalf("AccessRepositoryMock.GetRoleEndpoints mock is already set by Set")
	}

	if mmGetRoleEndpoints.defaultExpectation == nil {
		mmGetRoleEndpoints.defaultExpectation = &AccessRepositoryMockGetRoleEndpointsExpectation{}
	}

	mmGetRoleEndpoints.defaultExpectation.params = &AccessRepositoryMockGetRoleEndpointsParams{ctx}
	for _, e := range mmGetRoleEndpoints.expectations {
		if minimock.Equal(e.params, mmGetRoleEndpoints.defaultExpectation.params) {
			mmGetRoleEndpoints.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetRoleEndpoints.defaultExpectation.params)
		}
	}

	return mmGetRoleEndpoints
}

// Inspect accepts an inspector function that has same arguments as the AccessRepository.GetRoleEndpoints
func (mmGetRoleEndpoints *mAccessRepositoryMockGetRoleEndpoints) Inspect(f func(ctx context.Context)) *mAccessRepositoryMockGetRoleEndpoints {
	if mmGetRoleEndpoints.mock.inspectFuncGetRoleEndpoints != nil {
		mmGetRoleEndpoints.mock.t.Fatalf("Inspect function is already set for AccessRepositoryMock.GetRoleEndpoints")
	}

	mmGetRoleEndpoints.mock.inspectFuncGetRoleEndpoints = f

	return mmGetRoleEndpoints
}

// Return sets up results that will be returned by AccessRepository.GetRoleEndpoints
func (mmGetRoleEndpoints *mAccessRepositoryMockGetRoleEndpoints) Return(epa1 []*model.EndpointPermissions, err error) *AccessRepositoryMock {
	if mmGetRoleEndpoints.mock.funcGetRoleEndpoints != nil {
		mmGetRoleEndpoints.mock.t.Fatalf("AccessRepositoryMock.GetRoleEndpoints mock is already set by Set")
	}

	if mmGetRoleEndpoints.defaultExpectation == nil {
		mmGetRoleEndpoints.defaultExpectation = &AccessRepositoryMockGetRoleEndpointsExpectation{mock: mmGetRoleEndpoints.mock}
	}
	mmGetRoleEndpoints.defaultExpectation.results = &AccessRepositoryMockGetRoleEndpointsResults{epa1, err}
	return mmGetRoleEndpoints.mock
}

// Set uses given function f to mock the AccessRepository.GetRoleEndpoints method
func (mmGetRoleEndpoints *mAccessRepositoryMockGetRoleEndpoints) Set(f func(ctx context.Context) (epa1 []*model.EndpointPermissions, err error)) *AccessRepositoryMock {
	if mmGetRoleEndpoints.defaultExpectation != nil {
		mmGetRoleEndpoints.mock.t.Fatalf("Default expectation is already set for the AccessRepository.GetRoleEndpoints method")
	}

	if len(mmGetRoleEndpoints.expectations) > 0 {
		mmGetRoleEndpoints.mock.t.Fatalf("Some expectations are already set for the AccessRepository.GetRoleEndpoints method")
	}

	mmGetRoleEndpoints.mock.funcGetRoleEndpoints = f
	return mmGetRoleEndpoints.mock
}

// When sets expectation for the AccessRepository.GetRoleEndpoints which will trigger the result defined by the following
// Then helper
func (mmGetRoleEndpoints *mAccessRepositoryMockGetRoleEndpoints) When(ctx context.Context) *AccessRepositoryMockGetRoleEndpointsExpectation {
	if mmGetRoleEndpoints.mock.funcGetRoleEndpoints != nil {
		mmGetRoleEndpoints.mock.t.Fatalf("AccessRepositoryMock.GetRoleEndpoints mock is already set by Set")
	}

	expectation := &AccessRepositoryMockGetRoleEndpointsExpectation{
		mock:   mmGetRoleEndpoints.mock,
		params: &AccessRepositoryMockGetRoleEndpointsParams{ctx},
	}
	mmGetRoleEndpoints.expectations = append(mmGetRoleEndpoints.expectations, expectation)
	return expectation
}

// Then sets up AccessRepository.GetRoleEndpoints return parameters for the expectation previously defined by the When method
func (e *AccessRepositoryMockGetRoleEndpointsExpectation) Then(epa1 []*model.EndpointPermissions, err error) *AccessRepositoryMock {
	e.results = &AccessRepositoryMockGetRoleEndpointsResults{epa1, err}
	return e.mock
}

// GetRoleEndpoints implements repository.AccessRepository
func (mmGetRoleEndpoints *AccessRepositoryMock) GetRoleEndpoints(ctx context.Context) (epa1 []*model.EndpointPermissions, err error) {
	mm_atomic.AddUint64(&mmGetRoleEndpoints.beforeGetRoleEndpointsCounter, 1)
	defer mm_atomic.AddUint64(&mmGetRoleEndpoints.afterGetRoleEndpointsCounter, 1)

	if mmGetRoleEndpoints.inspectFuncGetRoleEndpoints != nil {
		mmGetRoleEndpoints.inspectFuncGetRoleEndpoints(ctx)
	}

	mm_params := AccessRepositoryMockGetRoleEndpointsParams{ctx}

	// Record call args
	mmGetRoleEndpoints.GetRoleEndpointsMock.mutex.Lock()
	mmGetRoleEndpoints.GetRoleEndpointsMock.callArgs = append(mmGetRoleEndpoints.GetRoleEndpointsMock.callArgs, &mm_params)
	mmGetRoleEndpoints.GetRoleEndpointsMock.mutex.Unlock()

	for _, e := range mmGetRoleEndpoints.GetRoleEndpointsMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.epa1, e.results.err
		}
	}

	if mmGetRoleEndpoints.GetRoleEndpointsMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetRoleEndpoints.GetRoleEndpointsMock.defaultExpectation.Counter, 1)
		mm_want := mmGetRoleEndpoints.GetRoleEndpointsMock.defaultExpectation.params
		mm_got := AccessRepositoryMockGetRoleEndpointsParams{ctx}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetRoleEndpoints.t.Errorf("AccessRepositoryMock.GetRoleEndpoints got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetRoleEndpoints.GetRoleEndpointsMock.defaultExpectation.results
		if mm_results == nil {
			mmGetRoleEndpoints.t.Fatal("No results are set for the AccessRepositoryMock.GetRoleEndpoints")
		}
		return (*mm_results).epa1, (*mm_results).err
	}
	if mmGetRoleEndpoints.funcGetRoleEndpoints != nil {
		return mmGetRoleEndpoints.funcGetRoleEndpoints(ctx)
	}
	mmGetRoleEndpoints.t.Fatalf("Unexpected call to AccessRepositoryMock.GetRoleEndpoints. %v", ctx)
	return
}

// GetRoleEndpointsAfterCounter returns a count of finished AccessRepositoryMock.GetRoleEndpoints invocations
func (mmGetRoleEndpoints *AccessRepositoryMock) GetRoleEndpointsAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetRoleEndpoints.afterGetRoleEndpointsCounter)
}

// GetRoleEndpointsBeforeCounter returns a count of AccessRepositoryMock.GetRoleEndpoints invocations
func (mmGetRoleEndpoints *AccessRepositoryMock) GetRoleEndpointsBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetRoleEndpoints.beforeGetRoleEndpointsCounter)
}

// Calls returns a list of arguments used in each call to AccessRepositoryMock.GetRoleEndpoints.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetRoleEndpoints *mAccessRepositoryMockGetRoleEndpoints) Calls() []*AccessRepositoryMockGetRoleEndpointsParams {
	mmGetRoleEndpoints.mutex.RLock()

	argCopy := make([]*AccessRepositoryMockGetRoleEndpointsParams, len(mmGetRoleEndpoints.callArgs))
	copy(argCopy, mmGetRoleEndpoints.callArgs)

	mmGetRoleEndpoints.mutex.RUnlock()

	return argCopy
}

// MinimockGetRoleEndpointsDone returns true if the count of the GetRoleEndpoints invocations corresponds
// the number of defined expectations
func (m *AccessRepositoryMock) MinimockGetRoleEndpointsDone() bool {
	for _, e := range m.GetRoleEndpointsMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetRoleEndpointsMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetRoleEndpointsCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetRoleEndpoints != nil && mm_atomic.LoadUint64(&m.afterGetRoleEndpointsCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetRoleEndpointsInspect logs each unmet expectation
func (m *AccessRepositoryMock) MinimockGetRoleEndpointsInspect() {
	for _, e := range m.GetRoleEndpointsMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to AccessRepositoryMock.GetRoleEndpoints with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetRoleEndpointsMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetRoleEndpointsCounter) < 1 {
		if m.GetRoleEndpointsMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to AccessRepositoryMock.GetRoleEndpoints")
		} else {
			m.t.Errorf("Expected call to AccessRepositoryMock.GetRoleEndpoints with params: %#v", *m.GetRoleEndpointsMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetRoleEndpoints != nil && mm_atomic.LoadUint64(&m.afterGetRoleEndpointsCounter) < 1 {
		m.t.Error("Expected call to AccessRepositoryMock.GetRoleEndpoints")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *AccessRepositoryMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockGetRoleEndpointsInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *AccessRepositoryMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *AccessRepositoryMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetRoleEndpointsDone()
}
