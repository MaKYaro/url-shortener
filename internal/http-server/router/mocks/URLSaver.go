// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	domain "github.com/MaKYaro/url-shortener/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// URLSaver is an autogenerated mock type for the URLSaver type
type URLSaver struct {
	mock.Mock
}

// SaveURL provides a mock function with given fields: url
func (_m *URLSaver) SaveURL(url string) (*domain.Alias, error) {
	ret := _m.Called(url)

	if len(ret) == 0 {
		panic("no return value specified for SaveURL")
	}

	var r0 *domain.Alias
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.Alias, error)); ok {
		return rf(url)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.Alias); ok {
		r0 = rf(url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Alias)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewURLSaver creates a new instance of URLSaver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewURLSaver(t interface {
	mock.TestingT
	Cleanup(func())
}) *URLSaver {
	mock := &URLSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
