package neoroute

import (
	"reflect"
	"testing"
)

func TestNewSession(t *testing.T) {
	type args struct {
		id   string
		data string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "create string session",
			args: args{id: "session-123", data: "hello_world"},
			want: "hello_world",
		},
		{
			name: "create string session without id and data",
			args: args{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSession(tt.args.id, tt.args.data)
			if got.id != tt.args.id {
				t.Errorf("NewSession() id = %v, want %v", got.id, tt.args.id)
			}
			if got.sessionData != tt.want {
				t.Errorf("NewSession() data = %v, want %v", got.sessionData, tt.want)
			}
		})
	}
}

func TestSession_Id(t *testing.T) {
	tests := []struct {
		name string
		s    *Session[string]
		want string
	}{
		{
			name: "get session id",
			s:    NewSession("id-abc", "data"),
			want: "id-abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Id(); got != tt.want {
				t.Errorf("Session.Id() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_Data(t *testing.T) {
	tests := []struct {
		name string
		s    *Session[string]
		want string
	}{
		{
			name: "get session data",
			s:    NewSession("id", "my-string-data"),
			want: "my-string-data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Data(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Session.Data() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_SetData(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		s    *Session[string]
		args args
		want string
	}{
		{
			name: "update data completely",
			s:    NewSession("id", "old-data"),
			args: args{data: "new-data"},
			want: "new-data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.SetData(tt.args.data)
			if got := tt.s.Data(); got != tt.want {
				t.Errorf("Session.SetData() resulted in data = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_UpdateData(t *testing.T) {
	type args struct {
		updateFunc func(data *string)
	}
	tests := []struct {
		name string
		s    *Session[string]
		args args
		want string
	}{
		{
			name: "mutate data via function closure",
			s:    NewSession("id", "hello"),
			args: args{
				updateFunc: func(data *string) {
					*data += " world"
				},
			},
			want: "hello world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.UpdateData(tt.args.updateFunc)
			if got := tt.s.Data(); got != tt.want {
				t.Errorf("Session.UpdateData() resulted in data = %v, want %v", got, tt.want)
			}
		})
	}
}
