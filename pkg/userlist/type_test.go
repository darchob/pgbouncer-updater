package userlist

import (
	"bytes"
	"testing"
)

func TestUser_WriteMany(t *testing.T) {
	data := []byte{}
	buf := bytes.NewBuffer(data)

	list := map[string]string{
		"postgres": "postgres",
	}

	type args struct {
		users interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				users: list,
			},
			wantErr: false,
		},
		{
			args: args{
				users: []*User{
					{
						UserName: "postgres",
						Md5:      "postgres",
					},
				},
			},
			wantErr: false,
		},
		{
			args: args{
				users: []User{
					{
						UserName: "postgres",
						Md5:      "postgres",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				file: buf,
			}

			if err := u.WriteMany(tt.args.users); (err != nil) != tt.wantErr {
				t.Errorf("User.WriteMany() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

	}
}
