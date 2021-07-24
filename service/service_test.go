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
		ctx = context.Background()
		err = errors.New("some error")
		user = entity.User{
			ID:      "id-1",
			Account: entity.Account{
				Email:    "sulam3010@gmail.com",
				Password: "password",
			},
			Profile: entity.Profile{},
		}
		hashedPassword = "hashedPassword"
		pushToken = "asdf1234"
		req = api.RegisterReq{
			Email:     user.Account.Email,
			Password:  user.Account.Password,
			PushToken: pushToken,
		}
		token = "hadslfkjwq1434"
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
		{
			name:    "invalid request",
			prepare: nil,
			args:    args{
				ctx: ctx,
				req: api.RegisterReq{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error when find user from db",
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
			name:    "email already registered",
			prepare: func() {
				mockUserModule.On("FindUserByEmail", mock.Anything, req.Email).
					Return(&user, nil)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "email already registered",
			prepare: func() {
				mockUserModule.On("FindUserByEmail", mock.Anything, req.Email).
					Return(&user, nil)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error when hash password",
			prepare: func() {
				mockUserModule.On("FindUserByEmail", mock.Anything, req.Email).
					Return(nil, constants.ErrUserNotFound)
				mockAuthModule.On("HashPassword", mock.Anything, req.Password).
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
			name:    "error when insert user",
			prepare: func() {
				mockUserModule.On("FindUserByEmail", mock.Anything, req.Email).
					Return(nil, constants.ErrUserNotFound)
				mockAuthModule.On("HashPassword", mock.Anything, req.Password).
					Return(hashedPassword, nil)
				mockUserModule.On("InsertUser", mock.Anything, mock.Anything).
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
			name:    "error when generate token",
			prepare: func() {
				mockUserModule.On("FindUserByEmail", mock.Anything, req.Email).
					Return(nil, constants.ErrUserNotFound)
				mockAuthModule.On("HashPassword", mock.Anything, req.Password).
					Return(hashedPassword, nil)
				mockUserModule.On("InsertUser", mock.Anything, mock.MatchedBy(func(u entity.User) bool{
					assert.Equal(t, req.Email, u.Account.Email)
					assert.Equal(t, hashedPassword, u.Account.Password)
					return true
				})).Return(user.ID, nil)
				user.Account.Password = hashedPassword
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
					Return(nil, constants.ErrUserNotFound)
				mockAuthModule.On("HashPassword", mock.Anything, req.Password).
					Return(hashedPassword, nil)
				mockUserModule.On("InsertUser", mock.Anything, mock.MatchedBy(func(u entity.User) bool{
					assert.Equal(t, req.Email, u.Account.Email)
					assert.Equal(t, hashedPassword, u.Account.Password)
					return true
				})).Return(user.ID, nil)
				user.Account.Password = hashedPassword
				mockAuthModule.On("GenerateToken", mock.Anything, user).
					Return(token, nil)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    &api.RegisterRes{
				Token: token,
			},
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
		ctx = context.Background()
		req = api.UpdateProfileReq{
			ID:     "id",
			Name:   "John",
			Avatar: "avatar",
			Bio:    "test",
		}
		err = errors.New("some error")
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
		{
			name:    "invalid request",
			prepare: nil,
			args:    args{
				ctx: ctx,
				req: api.UpdateProfileReq{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error when get user from db",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.ID).
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
			name:    "error when save profile",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.ID).
					Return(&entity.User{ID: req.ID}, nil)
				mockUserModule.On("SaveProfile", mock.Anything, mock.Anything).
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
			name:    "success",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.ID).
					Return(&entity.User{ID: req.ID}, nil)
				mockUserModule.On("SaveProfile", mock.Anything, mock.MatchedBy(func(u entity.User) bool{
					assert.Equal(t, req.ID, u.ID)
					assert.Equal(t, req.Name, u.Profile.Name)
					assert.Equal(t, req.Avatar, u.Profile.Avatar)
					assert.Equal(t, req.Bio, u.Profile.Bio)
					return true
				})).Return(nil)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    &api.UpdateProfileRes{Message: "success"},
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
		ctx = context.Background()
		req = api.GetUserReq{
			ID: "id",
		}
		err = errors.New("some error")
		user = entity.User{
			ID:      "id",
			Profile: entity.Profile{
				Name: "john",
				Avatar: "avatar",
				Bio: "",
			},
		}
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
		{
			name:    "invalid request",
			prepare: nil,
			args:    args{
				ctx: ctx,
				req: api.GetUserReq{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error when get user",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.ID).
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
			name:    "success",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.ID).
					Return(&user, nil)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    &api.GetUserRes{User: user},
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
		ctx = context.Background()
		req = api.GetMenfessListReq{}
		err = errors.New("some error")
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
		{
			name:    "error when get menfess list",
			prepare: func() {
				mockUserModule.On("FindMenfessList", mock.Anything).
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
			name:    "success",
			prepare: func() {
				mockUserModule.On("FindMenfessList", mock.Anything).
					Return([]entity.User{}, nil)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    &api.GetMenfessListRes{MenfessList: []entity.User{}},
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
		ctx = context.Background()
		req = api.GetPostReq{
			ID:     "id",
			UserID: "user-id",
		}
		err = errors.New("some error")
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
		{
			name:    "invalid request",
			prepare: nil,
			args:    args{
				ctx: ctx,
				req: api.GetPostReq{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error when get post",
			prepare: func() {
				mockPostModule.On("FindPostByID", mock.Anything, req.ID, req.UserID).
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
			name:    "error post not found",
			prepare: func() {
				mockPostModule.On("FindPostByID", mock.Anything, req.ID, req.UserID).
					Return(nil, constants.ErrPostNotFound)
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
				mockPostModule.On("FindPostByID", mock.Anything, req.ID, req.UserID).
					Return(&entity.Post{}, nil)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    &api.GetPostRes{Post: entity.Post{}},
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
		ctx = context.Background()
		err = errors.New("some error")
		req = api.GetPostListReq{
			ParentID:   "asdf",
			AuthorIDs:  []string{"123", "456"},
			UserID:     "user-id",
			Pagination: api.PaginationReq{
				First: 10,
				After: "1234",
			},
		}
		paginationRes = api.PaginationRes{
			EndCursor:   "some cursor",
			HasNextPage: false,
		}
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
		{
			name:    "invalid request",
			prepare: nil,
			args:    args{
				ctx: ctx,
				req: api.GetPostListReq{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error when get post list",
			prepare: func() {
				mockPostModule.On("FindPostListByParentIDAndAuthorIDs",
					mock.Anything,
					req.ParentID,
					req.AuthorIDs,
					req.UserID,
					req.Pagination,
				).Return(nil, nil, err)
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
				mockPostModule.On("FindPostListByParentIDAndAuthorIDs",
					mock.Anything,
					req.ParentID,
					req.AuthorIDs,
					req.UserID,
					req.Pagination,
				).Return([]entity.Post{}, &paginationRes, nil)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    &api.GetPostListRes{
				PostList:   []entity.Post{},
				Pagination: paginationRes,
			},
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
		ctx = context.Background()
		err = errors.New("some error")
		req = api.CreatePostReq{
			Body:     "body",
			UserID:   "user-id",
			AuthorID: "author-id",
			ParentID: "parent-id",
		}
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
		{
			name:    "invalid request",
			prepare: nil,
			args:    args{
				ctx: ctx,
				req: api.CreatePostReq{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error when get user from db",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.UserID).
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
			name:    "error when user not found",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.UserID).
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
			name:    "error when save post",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.UserID).
					Return(&entity.User{}, nil)
				mockPostModule.On("InsertPost", mock.Anything, mock.Anything).
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
			name:    "error when get post",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.UserID).
					Return(&entity.User{}, nil)
				mockPostModule.On("InsertPost", mock.Anything, mock.Anything).
					Return("post-id", nil)
				mockPostModule.On("FindPostByID", mock.Anything, req.ParentID, "").
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
			name:    "error when send notification",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.UserID).
					Return(&entity.User{}, nil)
				mockPostModule.On("InsertPost", mock.Anything, mock.Anything).
					Return("post-id", nil)
				mockPostModule.On("FindPostByID", mock.Anything, req.ParentID, "").
					Return(&entity.Post{}, nil)
				mockNotificationModule.On("SendCommentNotification", mock.Anything, entity.User{}, entity.Post{}).
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
			name:    "success",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.UserID).
					Return(&entity.User{}, nil)
				mockPostModule.On("InsertPost", mock.Anything, mock.MatchedBy(func(p entity.Post) bool {
					assert.Equal(t, req.Body, p.Body)
					assert.Equal(t, req.AuthorID, p.Author.ID)
					assert.Equal(t, req.ParentID, p.Parent.ID)
					assert.Equal(t, req.UserID, p.User.ID)
					return true
				})).Return("post-id", nil)
				mockPostModule.On("FindPostByID", mock.Anything, req.ParentID, "").
					Return(&entity.Post{}, nil)
				mockNotificationModule.On("SendCommentNotification", mock.Anything, entity.User{}, entity.Post{}).
					Return(nil)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    &api.CreatePostRes{Message: "success"},
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
		ctx = context.Background()
		err = errors.New("some error")
		req = api.LikePostReq{
			PostID: "post-id",
			UserID: "user-id",
		}
		post = entity.Post{
			ID:      "post-id",
			IsLiked: false,
			User:    entity.User{ID: "user-post-id"},
		}
		user = entity.User{
			ID:      "id",
			Profile: entity.Profile{
				Name: "John",
			},
		}
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
		{
			name:    "invalid request",
			prepare: nil,
			args:    args{
				ctx: ctx,
				req: api.LikePostReq{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error when get user",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.UserID).
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
			name:    "error user not found",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.UserID).
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
			name:    "error when get post",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.UserID).
					Return(&entity.User{}, nil)
				mockPostModule.On("FindPostByID", mock.Anything, req.PostID, req.UserID).
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
			name:    "error post not found",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.UserID).
					Return(&entity.User{}, nil)
				mockPostModule.On("FindPostByID", mock.Anything, req.PostID, req.UserID).
					Return(nil, constants.ErrPostNotFound)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error when like post",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.UserID).
					Return(&entity.User{}, nil)
				mockPostModule.On("FindPostByID", mock.Anything, req.PostID, req.UserID).
					Return(&post, nil)
				mockPostModule.On("LikePost", mock.Anything, req.PostID, req.UserID).
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
			name:    "error when unlike post",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.UserID).
					Return(&entity.User{}, nil)
				mockPostModule.On("FindPostByID", mock.Anything, req.PostID, req.UserID).
					Return(&entity.Post{IsLiked: true}, nil)
				mockPostModule.On("UnlikePost", mock.Anything, req.PostID, req.UserID).
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
			name:    "error when send notification",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.UserID).
					Return(&user, nil)
				mockPostModule.On("FindPostByID", mock.Anything, req.PostID, req.UserID).
					Return(&post, nil)
				mockPostModule.On("LikePost", mock.Anything, req.PostID, req.UserID).
					Return(nil)
				mockNotificationModule.On("SendLikeNotification", mock.Anything,
					user,
					post,
				).Return(err)
				mockLogModule.On("Log", err, req, mock.Anything)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    &api.LikePostRes{Message: "success"},
			wantErr: false,
		},
		{
			name:    "success like post",
			prepare: func() {
				mockUserModule.On("FindUserByID", mock.Anything, req.UserID).
					Return(&user, nil)
				mockPostModule.On("FindPostByID", mock.Anything, req.PostID, req.UserID).
					Return(&post, nil)
				mockPostModule.On("LikePost", mock.Anything, req.PostID, req.UserID).
					Return(nil)
				mockNotificationModule.On("SendLikeNotification", mock.Anything,
					user,
					post,
				).Return(nil)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    &api.LikePostRes{Message: "success"},
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

func Test_service_Logout(t *testing.T)  {
	var (
		ctx = context.Background()
		err = errors.New("some error")
		req = api.LogoutReq{
			UserID:    "user-id",
			PushToken: "token",
		}
	)
	type args struct {
		ctx context.Context
		req api.LogoutReq
	}
	tests := []struct {
		name    string
		prepare func()
		args    args
		want    *api.LogoutRes
		wantErr bool
	}{
		{
			name:    "invalid argument",
			prepare: nil,
			args:    args{
				ctx: ctx,
				req : api.LogoutReq{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error when remove push token",
			prepare: func() {
				mockNotificationModule.On("RemovePushToken", mock.Anything, req.UserID, req.PushToken).
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
			name:    "success",
			prepare: func() {
				mockNotificationModule.On("RemovePushToken", mock.Anything, req.UserID, req.PushToken).
					Return(nil)
			},
			args:    args{
				ctx: ctx,
				req: req,
			},
			want:    &api.LogoutRes{Message: "Success"},
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
			got, err := s.Logout(tt.args.ctx, tt.args.req)
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
