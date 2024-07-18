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

// SaveURL provides a mock function with given fields: alias
func (_m *URLSaver) SaveURL(alias *domain.Alias) error {
	ret := _m.Called(alias)

	if len(ret) == 0 {
		panic("no return value specified for SaveURL")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Alias) error); ok {
		r0 = rf(alias)
	} else {
		r0 = ret.Error(0)
	}

	return r0
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