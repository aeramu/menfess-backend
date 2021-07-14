// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	api "github.com/aeramu/menfess-backend/service/api"

	entity "github.com/aeramu/menfess-backend/entity"

	mock "github.com/stretchr/testify/mock"
)

// PostModule is an autogenerated mock type for the PostModule type
type PostModule struct {
	mock.Mock
}

// FindPostByID provides a mock function with given fields: ctx, id, userID
func (_m *PostModule) FindPostByID(ctx context.Context, id string, userID string) (*entity.Post, error) {
	ret := _m.Called(ctx, id, userID)

	var r0 *entity.Post
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *entity.Post); ok {
		r0 = rf(ctx, id, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Post)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, id, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindPostListByParentIDAndAuthorIDs provides a mock function with given fields: ctx, parentID, authorIDs, userID, pagination
func (_m *PostModule) FindPostListByParentIDAndAuthorIDs(ctx context.Context, parentID string, authorIDs []string, userID string, pagination api.PaginationReq) ([]entity.Post, *api.PaginationRes, error) {
	ret := _m.Called(ctx, parentID, authorIDs, userID, pagination)

	var r0 []entity.Post
	if rf, ok := ret.Get(0).(func(context.Context, string, []string, string, api.PaginationReq) []entity.Post); ok {
		r0 = rf(ctx, parentID, authorIDs, userID, pagination)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Post)
		}
	}

	var r1 *api.PaginationRes
	if rf, ok := ret.Get(1).(func(context.Context, string, []string, string, api.PaginationReq) *api.PaginationRes); ok {
		r1 = rf(ctx, parentID, authorIDs, userID, pagination)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*api.PaginationRes)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, []string, string, api.PaginationReq) error); ok {
		r2 = rf(ctx, parentID, authorIDs, userID, pagination)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
