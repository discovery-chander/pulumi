// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import gorm "github.com/jinzhu/gorm"
import mock "github.com/stretchr/testify/mock"

// Migration is an autogenerated mock type for the Migration type
type Migration struct {
	mock.Mock
}

// UpdateTables provides a mock function with given fields: _a0, _a1
func (_m *Migration) UpdateTables(_a0 []interface{}, _a1 *gorm.DB) {
	_m.Called(_a0, _a1)
}
