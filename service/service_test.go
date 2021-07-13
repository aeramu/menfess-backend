package service

import (
	"context"
	"errors"
	"github.com/aeramu/menfess-backend/constants"
	"github.com/aeramu/menfess-backend/entity"
	"github.com/aeramu/menfess-backend/mocks"
	"github.com/aeramu/menfess-backend/service/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var (
	adapter        Adapter
	mockUserModule *mocks.UserModule
	mockPostModule *mocks.PostModule
	mockAuthModule *mocks.AuthModule
	mockLogModule  *mocks.LogModule
	mockNotificationModule *mocks.NotificationModule
)

func initTest()  {
	mockUserModule = new(mocks.UserModule)
	mockPostModule = new(mocks.PostModule)
	mockAuthModule = new(mocks.AuthModule)
	mockNotificationModule = new(mocks.NotificationModule)
	mockLogModule = new(mocks.LogModule)
	adapter = Adapter{
		UserModule:         mockUserModule,
		PostModule:         mockPostModule,
		AuthModule:         mockAuthModule,
		NotificationModule: mockNotificationModule,
		LogModule:          mockLogModule,
	}
}

func Test_service_Login(t *testing.T)  {
	var (
		ctx = context.Background()
		err = errors.New("some error")
		user = entity.User{
			ID:      "id-1",
			Account: entity.Account{
				Email:    "sulam3010@gmail.com",
				Password: "hashedPassword",
			},
			Profile: entity.Profile{},
		}
		pushToken = "asdf1234"
		req = api.LoginReq{
			Email:     user.Account.Email,
			Password:  user.Account.Password,
			PushToken: pushToken,
		}
		token = "hadslfkjwq1434"
	)
	type args struct {
		ctx context.Context
		req api.LoginReq
	}
	tests := []struct {
		name    string
		prepare func()
		args    args
		want    *api.LoginRes
		wantErr bool
	}{
		{
			name:    "invalid request",
			prepare: nil,
			args:    args{
				ctx: ctx,
				req: api.LoginReq{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "user not found",
			prepare: func() {
				mockUserModule.On("FindUserByEmail", mock.Anything, req.Email).
					Return(nil, constants.ErrUserNotFound)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error when find user from repo",
			prepare: func() {
				mockUserModule.On("FindUserByEmail", mock.Anything, req.Email).
					Return(nil, err)
				mockLogModule.On("Log", err, req, mock.Anything)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "wrong password",
			prepare: func() {
				mockUserModule.On("FindUserByEmail", mock.Anything, req.Email).
					Return(&user, nil)
				mockAuthModule.On("ComparePassword", mock.Anything, user.Account.Password, req.Password).
					Return(errors.New("wrong password"))
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error when add push token",
			prepare: func() {
				mockUserModule.On("FindUserByEmail", mock.Anything, req.Email).
					Return(&user, nil)
				mockAuthModule.On("ComparePassword", mock.Anything, user.Account.Password, req.Password).
					Return(nil)
				mockNotificationModule.On("AddPushToken", mock.Anything, user.ID, req.PushToken).
					Return(err)
				mockLogModule.On("Log", err, req, mock.Anything)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error when generate token",
			prepare: func() {
				mockUserModule.On("FindUserByEmail", mock.Anything, req.Email).
					Return(&user, nil)
				mockAuthModule.On("ComparePassword", mock.Anything, user.Account.Password, req.Password).
					Return(nil)
				mockNotificationModule.On("AddPushToken", mock.Anything, user.ID, req.PushToken).
					Return(nil)
				mockAuthModule.On("GenerateToken", mock.Anything, user).
					Return("", err)
				mockLogModule.On("Log", err, req, mock.Anything)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "success",
			prepare: func() {
				mockUserModule.On("FindUserByEmail", mock.Anything, req.Email).
					Return(&user, nil)
				mockAuthModule.On("ComparePassword", mock.Anything, user.Account.Password, req.Password).
					Return(nil)
				mockNotificationModule.On("AddPushToken", mock.Anything, user.ID, req.PushToken).
					Return(nil)
				mockAuthModule.On("GenerateToken", mock.Anything, user).
					Return(token, nil)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    &api.LoginRes{Token: token},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initTest()
			if tt.prepare != nil {
				tt.prepare()
			}
			s := &service{
				adapter: adapter,
			}
			got, err := s.Login(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_service_Register(t *testing.T)  {
	var (

	)
	type args struct {
		ctx context.Context
		req api.RegisterReq
	}
	tests := []struct {
		name    string
		prepare func()
		args    args
		want    *api.RegisterRes
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initTest()
			if tt.prepare != nil {
				tt.prepare()
			}
			s := &service{
				adapter: adapter,
			}
			got, err := s.Register(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_service_UpdateProfile(t *testing.T)  {
	var (

	)
	type args struct {
		ctx context.Context
		req api.UpdateProfileReq
	}
	tests := []struct {
		name    string
		prepare func()
		args    args
		want    *api.UpdateProfileRes
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initTest()
			if tt.prepare != nil {
				tt.prepare()
			}
			s := &service{
				adapter: adapter,
			}
			got, err := s.UpdateProfile(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_service_GetUser(t *testing.T)  {
	var (

	)
	type args struct {
		ctx context.Context
		req api.GetUserReq
	}
	tests := []struct {
		name    string
		prepare func()
		args    args
		want    *api.GetUserRes
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initTest()
			if tt.prepare != nil {
				tt.prepare()
			}
			s := &service{
				adapter: adapter,
			}
			got, err := s.GetUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_service_GetMenfessList(t *testing.T)  {
	var (

	)
	type args struct {
		ctx context.Context
		req api.GetMenfessListReq
	}
	tests := []struct {
		name    string
		prepare func()
		args    args
		want    *api.GetMenfessListRes
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initTest()
			if tt.prepare != nil {
				tt.prepare()
			}
			s := &service{
				adapter: adapter,
			}
			got, err := s.GetMenfessList(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_service_GetPost(t *testing.T)  {
	var (

	)
	type args struct {
		ctx context.Context
		req api.GetPostReq
	}
	tests := []struct {
		name    string
		prepare func()
		args    args
		want    *api.GetPostRes
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initTest()
			if tt.prepare != nil {
				tt.prepare()
			}
			s := &service{
				adapter: adapter,
			}
			got, err := s.GetPost(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_service_GetPostList(t *testing.T)  {
	var (

	)
	type args struct {
		ctx context.Context
		req api.GetPostListReq
	}
	tests := []struct {
		name    string
		prepare func()
		args    args
		want    *api.GetPostListRes
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initTest()
			if tt.prepare != nil {
				tt.prepare()
			}
			s := &service{
				adapter: adapter,
			}
			got, err := s.GetPostList(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_service_CreatePost(t *testing.T)  {
	var (

	)
	type args struct {
		ctx context.Context
		req api.CreatePostReq
	}
	tests := []struct {
		name    string
		prepare func()
		args    args
		want    *api.CreatePostRes
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initTest()
			if tt.prepare != nil {
				tt.prepare()
			}
			s := &service{
				adapter: adapter,
			}
			got, err := s.CreatePost(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_service_LikePost(t *testing.T)  {
	var (

	)
	type args struct {
		ctx context.Context
		req api.LikePostReq
	}
	tests := []struct {
		name    string
		prepare func()
		args    args
		want    *api.LikePostRes
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initTest()
			if tt.prepare != nil {
				tt.prepare()
			}
			s := &service{
				adapter: adapter,
			}
			got, err := s.LikePost(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
