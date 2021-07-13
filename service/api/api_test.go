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
				Email:     "sulam3010gmail.com",
				Password:  "",
				PushToken: "sadf1234fas",
			},
			wantErr: true,
		},
		{
			name:    "empty push token",
			fields:  fields{
				Email:     "sulam3010gmail.com",
				Password:  "password",
				PushToken: "",
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
