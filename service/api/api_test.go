package api

import "testing"

func TestLoginReq_Validate(t *testing.T) {
	type fields struct {
		Email     string
		Password  string
		PushToken string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "success",
			fields:  fields{
				Email:     "sulam3010@gmail.com",
				Password:  "password",
				PushToken: "sadf1234fas",
			},
			wantErr: false,
		},
		{
			name:    "empty email",
			fields:  fields{
				Email:     "",
				Password:  "password",
				PushToken: "sadf1234fas",
			},
			wantErr: true,
		},
		{
			name:    "invalid email",
			fields:  fields{
				Email:     "sulam3010gmail.com",
				Password:  "password",
				PushToken: "sadf1234fas",
			},
			wantErr: true,
		},
		{
			name:    "empty password",
			fields:  fields{
				Email:     "sulam3010@gmail.com",
				Password:  "",
				PushToken: "sadf1234fas",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := LoginReq{
				Email:     tt.fields.Email,
				Password:  tt.fields.Password,
				PushToken: tt.fields.PushToken,
			}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterReq_Validate(t *testing.T) {
	type fields struct {
		Email     string
		Password  string
		PushToken string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "success",
			fields:  fields{
				PushToken: "sadf1234fas",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := RegisterReq{
				PushToken: tt.fields.PushToken,
			}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateProfileReq_Validate(t *testing.T) {
	type fields struct {
		ID     string
		Name   string
		Avatar string
		Bio    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty id",
			fields:  fields{
				ID:     "",
				Name:   "",
				Avatar: "",
				Bio:    "",
			},
			wantErr: true,
		},
		{
			name:    "empty name",
			fields:  fields{
				ID:     "id",
				Name:   "",
				Avatar: "avatar",
				Bio:    "",
			},
			wantErr: true,
		},
		{
			name:    "name only space",
			fields:  fields{
				ID:     "id",
				Name:   " ",
				Avatar: "avatar",
				Bio:    "",
			},
			wantErr: true,
		},
		{
			name:    "empty avatar",
			fields:  fields{
				ID:     "id",
				Name:   "john",
				Avatar: "",
				Bio:    "",
			},
			wantErr: true,
		},
		{
			name:    "avatar only space",
			fields:  fields{
				ID:     "id",
				Name:   "john",
				Avatar: " ",
				Bio:    "",
			},
			wantErr: true,
		},
		{
			name:    "success",
			fields:  fields{
				ID:     "id",
				Name:   "john",
				Avatar: "avatar",
				Bio:    "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := UpdateProfileReq{
				ID:     tt.fields.ID,
				Name:   tt.fields.Name,
				Avatar: tt.fields.Avatar,
				Bio:    tt.fields.Bio,
			}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserReq_Validate(t *testing.T) {
	type fields struct {
		ID string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty id",
			fields:  fields{},
			wantErr: true,
		},
		{
			name:    "success",
			fields:  fields{
				ID: "id",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := GetUserReq{
				ID: tt.fields.ID,
			}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetMenfessListReq_Validate(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "success",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := GetMenfessListReq{}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetPostReq_Validate(t *testing.T) {
	type fields struct {
		ID     string
		UserID string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty id",
			fields:  fields{
				ID:     "",
				UserID: "",
			},
			wantErr: true,
		},
		{
			name:    "empty user id",
			fields:  fields{
				ID:     "id",
				UserID: "",
			},
			wantErr: true,
		},
		{
			name:    "success",
			fields:  fields{
				ID:     "id",
				UserID: "user-id",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := GetPostReq{
				ID:     tt.fields.ID,
				UserID: tt.fields.UserID,
			}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetPostListReq_Validate(t *testing.T) {
	type fields struct {
		ParentID   string
		AuthorIDs  []string
		UserID     string
		Pagination PaginationReq
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty user id",
			fields:  fields{
				ParentID:   "",
				AuthorIDs:  nil,
				UserID:     "",
				Pagination: PaginationReq{},
			},
			wantErr: true,
		},
		{
			name:    "success",
			fields:  fields{
				ParentID:   "",
				AuthorIDs:  nil,
				UserID:     "id",
				Pagination: PaginationReq{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &GetPostListReq{
				ParentID:   tt.fields.ParentID,
				UserID:     tt.fields.UserID,
				Pagination: tt.fields.Pagination,
			}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreatePostReq_Validate(t *testing.T) {
	type fields struct {
		Body     string
		UserID   string
		AuthorID string
		ParentID string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty body",
			fields:  fields{
				Body:     "",
				UserID:   "user-id",
				AuthorID: "author-id",
				ParentID: "parent-id",
			},
			wantErr: true,
		},
		{
			name:    "empty user id",
			fields:  fields{
				Body:     "body",
				UserID:   "",
				AuthorID: "author-id",
				ParentID: "parent-id",
			},
			wantErr: true,
		},
		{
			name:    "success",
			fields:  fields{
				Body:     "body",
				UserID:   "user-id",
				AuthorID: "author-id",
				ParentID: "parent-id",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreatePostReq{
				Body:     tt.fields.Body,
				UserID:   tt.fields.UserID,
				AuthorID: tt.fields.AuthorID,
				ParentID: tt.fields.ParentID,
			}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLikePostReq_Validate(t *testing.T) {
	type fields struct {
		PostID string
		UserID string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty post id",
			fields:  fields{
				PostID: "",
				UserID: "user-id",
			},
			wantErr: true,
		},
		{
			name:    "empty user id",
			fields:  fields{
				PostID: "post-id",
				UserID: "",
			},
			wantErr: true,
		},
		{
			name:    "empty post id",
			fields:  fields{
				PostID: "post-id",
				UserID: "user-id",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := LikePostReq{
				PostID: tt.fields.PostID,
				UserID: tt.fields.UserID,
			}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogoutReq_Validate(t *testing.T) {
	type fields struct {
		UserID    string
		PushToken string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty user id",
			fields:  fields{
				UserID:    "",
				PushToken: "token",
			},
			wantErr: true,
		},
		{
			name:    "empty token",
			fields:  fields{
				UserID:    "user-id",
				PushToken: "",
			},
			wantErr: true,
		},
		{
			name:    "success",
			fields:  fields{
				UserID:    "user-id",
				PushToken: "token",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := LogoutReq{
				UserID:    tt.fields.UserID,
				PushToken: tt.fields.PushToken,
			}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFeedReq_Validate(t *testing.T) {
	type fields struct {
		UserID     string
		Type       string
		Pagination PaginationReq
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "success",
			fields:  fields{
				UserID:     "id",
				Type:       "all",
				Pagination: PaginationReq{},
			},
			wantErr: false,
		},
		{
			name:    "invalid type",
			fields:  fields{
				UserID:     "id",
				Type:       "allType",
				Pagination: PaginationReq{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := FeedReq{
				UserID:     tt.fields.UserID,
				Type:       tt.fields.Type,
				Pagination: tt.fields.Pagination,
			}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}