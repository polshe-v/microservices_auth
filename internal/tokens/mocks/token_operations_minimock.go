// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

package mocks

import (
	"sync"
	mm_atomic "sync/atomic"
	"time"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"github.com/polshe-v/microservices_auth/internal/model"
)

// TokenOperationsMock implements tokens.TokenOperations
type TokenOperationsMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcGenerate          func(user model.User, secretKey []byte, duration time.Duration) (s1 string, err error)
	inspectFuncGenerate   func(user model.User, secretKey []byte, duration time.Duration)
	afterGenerateCounter  uint64
	beforeGenerateCounter uint64
	GenerateMock          mTokenOperationsMockGenerate

	funcVerify          func(tokenStr string, secretKey []byte) (up1 *model.UserClaims, err error)
	inspectFuncVerify   func(tokenStr string, secretKey []byte)
	afterVerifyCounter  uint64
	beforeVerifyCounter uint64
	VerifyMock          mTokenOperationsMockVerify
}

// NewTokenOperationsMock returns a mock for tokens.TokenOperations
func NewTokenOperationsMock(t minimock.Tester) *TokenOperationsMock {
	m := &TokenOperationsMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GenerateMock = mTokenOperationsMockGenerate{mock: m}
	m.GenerateMock.callArgs = []*TokenOperationsMockGenerateParams{}

	m.VerifyMock = mTokenOperationsMockVerify{mock: m}
	m.VerifyMock.callArgs = []*TokenOperationsMockVerifyParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mTokenOperationsMockGenerate struct {
	mock               *TokenOperationsMock
	defaultExpectation *TokenOperationsMockGenerateExpectation
	expectations       []*TokenOperationsMockGenerateExpectation

	callArgs []*TokenOperationsMockGenerateParams
	mutex    sync.RWMutex
}

// TokenOperationsMockGenerateExpectation specifies expectation struct of the TokenOperations.Generate
type TokenOperationsMockGenerateExpectation struct {
	mock    *TokenOperationsMock
	params  *TokenOperationsMockGenerateParams
	results *TokenOperationsMockGenerateResults
	Counter uint64
}

// TokenOperationsMockGenerateParams contains parameters of the TokenOperations.Generate
type TokenOperationsMockGenerateParams struct {
	user      model.User
	secretKey []byte
	duration  time.Duration
}

// TokenOperationsMockGenerateResults contains results of the TokenOperations.Generate
type TokenOperationsMockGenerateResults struct {
	s1  string
	err error
}

// Expect sets up expected params for TokenOperations.Generate
func (mmGenerate *mTokenOperationsMockGenerate) Expect(user model.User, secretKey []byte, duration time.Duration) *mTokenOperationsMockGenerate {
	if mmGenerate.mock.funcGenerate != nil {
		mmGenerate.mock.t.Fatalf("TokenOperationsMock.Generate mock is already set by Set")
	}

	if mmGenerate.defaultExpectation == nil {
		mmGenerate.defaultExpectation = &TokenOperationsMockGenerateExpectation{}
	}

	mmGenerate.defaultExpectation.params = &TokenOperationsMockGenerateParams{user, secretKey, duration}
	for _, e := range mmGenerate.expectations {
		if minimock.Equal(e.params, mmGenerate.defaultExpectation.params) {
			mmGenerate.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGenerate.defaultExpectation.params)
		}
	}

	return mmGenerate
}

// Inspect accepts an inspector function that has same arguments as the TokenOperations.Generate
func (mmGenerate *mTokenOperationsMockGenerate) Inspect(f func(user model.User, secretKey []byte, duration time.Duration)) *mTokenOperationsMockGenerate {
	if mmGenerate.mock.inspectFuncGenerate != nil {
		mmGenerate.mock.t.Fatalf("Inspect function is already set for TokenOperationsMock.Generate")
	}

	mmGenerate.mock.inspectFuncGenerate = f

	return mmGenerate
}

// Return sets up results that will be returned by TokenOperations.Generate
func (mmGenerate *mTokenOperationsMockGenerate) Return(s1 string, err error) *TokenOperationsMock {
	if mmGenerate.mock.funcGenerate != nil {
		mmGenerate.mock.t.Fatalf("TokenOperationsMock.Generate mock is already set by Set")
	}

	if mmGenerate.defaultExpectation == nil {
		mmGenerate.defaultExpectation = &TokenOperationsMockGenerateExpectation{mock: mmGenerate.mock}
	}
	mmGenerate.defaultExpectation.results = &TokenOperationsMockGenerateResults{s1, err}
	return mmGenerate.mock
}

// Set uses given function f to mock the TokenOperations.Generate method
func (mmGenerate *mTokenOperationsMockGenerate) Set(f func(user model.User, secretKey []byte, duration time.Duration) (s1 string, err error)) *TokenOperationsMock {
	if mmGenerate.defaultExpectation != nil {
		mmGenerate.mock.t.Fatalf("Default expectation is already set for the TokenOperations.Generate method")
	}

	if len(mmGenerate.expectations) > 0 {
		mmGenerate.mock.t.Fatalf("Some expectations are already set for the TokenOperations.Generate method")
	}

	mmGenerate.mock.funcGenerate = f
	return mmGenerate.mock
}

// When sets expectation for the TokenOperations.Generate which will trigger the result defined by the following
// Then helper
func (mmGenerate *mTokenOperationsMockGenerate) When(user model.User, secretKey []byte, duration time.Duration) *TokenOperationsMockGenerateExpectation {
	if mmGenerate.mock.funcGenerate != nil {
		mmGenerate.mock.t.Fatalf("TokenOperationsMock.Generate mock is already set by Set")
	}

	expectation := &TokenOperationsMockGenerateExpectation{
		mock:   mmGenerate.mock,
		params: &TokenOperationsMockGenerateParams{user, secretKey, duration},
	}
	mmGenerate.expectations = append(mmGenerate.expectations, expectation)
	return expectation
}

// Then sets up TokenOperations.Generate return parameters for the expectation previously defined by the When method
func (e *TokenOperationsMockGenerateExpectation) Then(s1 string, err error) *TokenOperationsMock {
	e.results = &TokenOperationsMockGenerateResults{s1, err}
	return e.mock
}

// Generate implements tokens.TokenOperations
func (mmGenerate *TokenOperationsMock) Generate(user model.User, secretKey []byte, duration time.Duration) (s1 string, err error) {
	mm_atomic.AddUint64(&mmGenerate.beforeGenerateCounter, 1)
	defer mm_atomic.AddUint64(&mmGenerate.afterGenerateCounter, 1)

	if mmGenerate.inspectFuncGenerate != nil {
		mmGenerate.inspectFuncGenerate(user, secretKey, duration)
	}

	mm_params := TokenOperationsMockGenerateParams{user, secretKey, duration}

	// Record call args
	mmGenerate.GenerateMock.mutex.Lock()
	mmGenerate.GenerateMock.callArgs = append(mmGenerate.GenerateMock.callArgs, &mm_params)
	mmGenerate.GenerateMock.mutex.Unlock()

	for _, e := range mmGenerate.GenerateMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.s1, e.results.err
		}
	}

	if mmGenerate.GenerateMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGenerate.GenerateMock.defaultExpectation.Counter, 1)
		mm_want := mmGenerate.GenerateMock.defaultExpectation.params
		mm_got := TokenOperationsMockGenerateParams{user, secretKey, duration}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGenerate.t.Errorf("TokenOperationsMock.Generate got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGenerate.GenerateMock.defaultExpectation.results
		if mm_results == nil {
			mmGenerate.t.Fatal("No results are set for the TokenOperationsMock.Generate")
		}
		return (*mm_results).s1, (*mm_results).err
	}
	if mmGenerate.funcGenerate != nil {
		return mmGenerate.funcGenerate(user, secretKey, duration)
	}
	mmGenerate.t.Fatalf("Unexpected call to TokenOperationsMock.Generate. %v %v %v", user, secretKey, duration)
	return
}

// GenerateAfterCounter returns a count of finished TokenOperationsMock.Generate invocations
func (mmGenerate *TokenOperationsMock) GenerateAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGenerate.afterGenerateCounter)
}

// GenerateBeforeCounter returns a count of TokenOperationsMock.Generate invocations
func (mmGenerate *TokenOperationsMock) GenerateBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGenerate.beforeGenerateCounter)
}

// Calls returns a list of arguments used in each call to TokenOperationsMock.Generate.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGenerate *mTokenOperationsMockGenerate) Calls() []*TokenOperationsMockGenerateParams {
	mmGenerate.mutex.RLock()

	argCopy := make([]*TokenOperationsMockGenerateParams, len(mmGenerate.callArgs))
	copy(argCopy, mmGenerate.callArgs)

	mmGenerate.mutex.RUnlock()

	return argCopy
}

// MinimockGenerateDone returns true if the count of the Generate invocations corresponds
// the number of defined expectations
func (m *TokenOperationsMock) MinimockGenerateDone() bool {
	for _, e := range m.GenerateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GenerateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGenerateCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGenerate != nil && mm_atomic.LoadUint64(&m.afterGenerateCounter) < 1 {
		return false
	}
	return true
}

// MinimockGenerateInspect logs each unmet expectation
func (m *TokenOperationsMock) MinimockGenerateInspect() {
	for _, e := range m.GenerateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to TokenOperationsMock.Generate with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GenerateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGenerateCounter) < 1 {
		if m.GenerateMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to TokenOperationsMock.Generate")
		} else {
			m.t.Errorf("Expected call to TokenOperationsMock.Generate with params: %#v", *m.GenerateMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGenerate != nil && mm_atomic.LoadUint64(&m.afterGenerateCounter) < 1 {
		m.t.Error("Expected call to TokenOperationsMock.Generate")
	}
}

type mTokenOperationsMockVerify struct {
	mock               *TokenOperationsMock
	defaultExpectation *TokenOperationsMockVerifyExpectation
	expectations       []*TokenOperationsMockVerifyExpectation

	callArgs []*TokenOperationsMockVerifyParams
	mutex    sync.RWMutex
}

// TokenOperationsMockVerifyExpectation specifies expectation struct of the TokenOperations.Verify
type TokenOperationsMockVerifyExpectation struct {
	mock    *TokenOperationsMock
	params  *TokenOperationsMockVerifyParams
	results *TokenOperationsMockVerifyResults
	Counter uint64
}

// TokenOperationsMockVerifyParams contains parameters of the TokenOperations.Verify
type TokenOperationsMockVerifyParams struct {
	tokenStr  string
	secretKey []byte
}

// TokenOperationsMockVerifyResults contains results of the TokenOperations.Verify
type TokenOperationsMockVerifyResults struct {
	up1 *model.UserClaims
	err error
}

// Expect sets up expected params for TokenOperations.Verify
func (mmVerify *mTokenOperationsMockVerify) Expect(tokenStr string, secretKey []byte) *mTokenOperationsMockVerify {
	if mmVerify.mock.funcVerify != nil {
		mmVerify.mock.t.Fatalf("TokenOperationsMock.Verify mock is already set by Set")
	}

	if mmVerify.defaultExpectation == nil {
		mmVerify.defaultExpectation = &TokenOperationsMockVerifyExpectation{}
	}

	mmVerify.defaultExpectation.params = &TokenOperationsMockVerifyParams{tokenStr, secretKey}
	for _, e := range mmVerify.expectations {
		if minimock.Equal(e.params, mmVerify.defaultExpectation.params) {
			mmVerify.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmVerify.defaultExpectation.params)
		}
	}

	return mmVerify
}

// Inspect accepts an inspector function that has same arguments as the TokenOperations.Verify
func (mmVerify *mTokenOperationsMockVerify) Inspect(f func(tokenStr string, secretKey []byte)) *mTokenOperationsMockVerify {
	if mmVerify.mock.inspectFuncVerify != nil {
		mmVerify.mock.t.Fatalf("Inspect function is already set for TokenOperationsMock.Verify")
	}

	mmVerify.mock.inspectFuncVerify = f

	return mmVerify
}

// Return sets up results that will be returned by TokenOperations.Verify
func (mmVerify *mTokenOperationsMockVerify) Return(up1 *model.UserClaims, err error) *TokenOperationsMock {
	if mmVerify.mock.funcVerify != nil {
		mmVerify.mock.t.Fatalf("TokenOperationsMock.Verify mock is already set by Set")
	}

	if mmVerify.defaultExpectation == nil {
		mmVerify.defaultExpectation = &TokenOperationsMockVerifyExpectation{mock: mmVerify.mock}
	}
	mmVerify.defaultExpectation.results = &TokenOperationsMockVerifyResults{up1, err}
	return mmVerify.mock
}

// Set uses given function f to mock the TokenOperations.Verify method
func (mmVerify *mTokenOperationsMockVerify) Set(f func(tokenStr string, secretKey []byte) (up1 *model.UserClaims, err error)) *TokenOperationsMock {
	if mmVerify.defaultExpectation != nil {
		mmVerify.mock.t.Fatalf("Default expectation is already set for the TokenOperations.Verify method")
	}

	if len(mmVerify.expectations) > 0 {
		mmVerify.mock.t.Fatalf("Some expectations are already set for the TokenOperations.Verify method")
	}

	mmVerify.mock.funcVerify = f
	return mmVerify.mock
}

// When sets expectation for the TokenOperations.Verify which will trigger the result defined by the following
// Then helper
func (mmVerify *mTokenOperationsMockVerify) When(tokenStr string, secretKey []byte) *TokenOperationsMockVerifyExpectation {
	if mmVerify.mock.funcVerify != nil {
		mmVerify.mock.t.Fatalf("TokenOperationsMock.Verify mock is already set by Set")
	}

	expectation := &TokenOperationsMockVerifyExpectation{
		mock:   mmVerify.mock,
		params: &TokenOperationsMockVerifyParams{tokenStr, secretKey},
	}
	mmVerify.expectations = append(mmVerify.expectations, expectation)
	return expectation
}

// Then sets up TokenOperations.Verify return parameters for the expectation previously defined by the When method
func (e *TokenOperationsMockVerifyExpectation) Then(up1 *model.UserClaims, err error) *TokenOperationsMock {
	e.results = &TokenOperationsMockVerifyResults{up1, err}
	return e.mock
}

// Verify implements tokens.TokenOperations
func (mmVerify *TokenOperationsMock) Verify(tokenStr string, secretKey []byte) (up1 *model.UserClaims, err error) {
	mm_atomic.AddUint64(&mmVerify.beforeVerifyCounter, 1)
	defer mm_atomic.AddUint64(&mmVerify.afterVerifyCounter, 1)

	if mmVerify.inspectFuncVerify != nil {
		mmVerify.inspectFuncVerify(tokenStr, secretKey)
	}

	mm_params := TokenOperationsMockVerifyParams{tokenStr, secretKey}

	// Record call args
	mmVerify.VerifyMock.mutex.Lock()
	mmVerify.VerifyMock.callArgs = append(mmVerify.VerifyMock.callArgs, &mm_params)
	mmVerify.VerifyMock.mutex.Unlock()

	for _, e := range mmVerify.VerifyMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.up1, e.results.err
		}
	}

	if mmVerify.VerifyMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmVerify.VerifyMock.defaultExpectation.Counter, 1)
		mm_want := mmVerify.VerifyMock.defaultExpectation.params
		mm_got := TokenOperationsMockVerifyParams{tokenStr, secretKey}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmVerify.t.Errorf("TokenOperationsMock.Verify got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmVerify.VerifyMock.defaultExpectation.results
		if mm_results == nil {
			mmVerify.t.Fatal("No results are set for the TokenOperationsMock.Verify")
		}
		return (*mm_results).up1, (*mm_results).err
	}
	if mmVerify.funcVerify != nil {
		return mmVerify.funcVerify(tokenStr, secretKey)
	}
	mmVerify.t.Fatalf("Unexpected call to TokenOperationsMock.Verify. %v %v", tokenStr, secretKey)
	return
}

// VerifyAfterCounter returns a count of finished TokenOperationsMock.Verify invocations
func (mmVerify *TokenOperationsMock) VerifyAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmVerify.afterVerifyCounter)
}

// VerifyBeforeCounter returns a count of TokenOperationsMock.Verify invocations
func (mmVerify *TokenOperationsMock) VerifyBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmVerify.beforeVerifyCounter)
}

// Calls returns a list of arguments used in each call to TokenOperationsMock.Verify.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmVerify *mTokenOperationsMockVerify) Calls() []*TokenOperationsMockVerifyParams {
	mmVerify.mutex.RLock()

	argCopy := make([]*TokenOperationsMockVerifyParams, len(mmVerify.callArgs))
	copy(argCopy, mmVerify.callArgs)

	mmVerify.mutex.RUnlock()

	return argCopy
}

// MinimockVerifyDone returns true if the count of the Verify invocations corresponds
// the number of defined expectations
func (m *TokenOperationsMock) MinimockVerifyDone() bool {
	for _, e := range m.VerifyMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.VerifyMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterVerifyCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcVerify != nil && mm_atomic.LoadUint64(&m.afterVerifyCounter) < 1 {
		return false
	}
	return true
}

// MinimockVerifyInspect logs each unmet expectation
func (m *TokenOperationsMock) MinimockVerifyInspect() {
	for _, e := range m.VerifyMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to TokenOperationsMock.Verify with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.VerifyMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterVerifyCounter) < 1 {
		if m.VerifyMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to TokenOperationsMock.Verify")
		} else {
			m.t.Errorf("Expected call to TokenOperationsMock.Verify with params: %#v", *m.VerifyMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcVerify != nil && mm_atomic.LoadUint64(&m.afterVerifyCounter) < 1 {
		m.t.Error("Expected call to TokenOperationsMock.Verify")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *TokenOperationsMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockGenerateInspect()

			m.MinimockVerifyInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *TokenOperationsMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *TokenOperationsMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGenerateDone() &&
		m.MinimockVerifyDone()
}
