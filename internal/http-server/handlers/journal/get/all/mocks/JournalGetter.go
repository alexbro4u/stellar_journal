// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	stellar_journal_models "stellar_journal/internal/models/stellar_journal_models"

	mock "github.com/stretchr/testify/mock"
)

// JournalGetter is an autogenerated mock type for the JournalGetter type
type JournalGetter struct {
	mock.Mock
}

// GetJournal provides a mock function with given fields:
func (_m *JournalGetter) GetJournal() (*[]stellar_journal_models.APOD, error) {
	ret := _m.Called()

	var r0 *[]stellar_journal_models.APOD
	var r1 error
	if rf, ok := ret.Get(0).(func() (*[]stellar_journal_models.APOD, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *[]stellar_journal_models.APOD); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]stellar_journal_models.APOD)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewJournalGetter interface {
	mock.TestingT
	Cleanup(func())
}

// NewJournalGetter creates a new instance of JournalGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewJournalGetter(t mockConstructorTestingTNewJournalGetter) *JournalGetter {
	mock := &JournalGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
