package api_key

import (
	"os"
	"reflect"
	"testing"
)

func TestGetApiKeys_1(t *testing.T) {
	os.Setenv(ENV_REQUEST_API_KEYS, "{ \"api_keys\": [ \"test-1\" ] }")

	tests := []struct {
		name string
		want RequestApiKeys
	}{
		{
			name: "#1: ENV_REQUEST_API_KEYS set",
			want: RequestApiKeys{
				ApiKeys: []string{"test-1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetApiKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetApiKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetApiKeys_2(t *testing.T) {
	os.Setenv(ENV_REQUEST_API_KEYS, "{ \"api_keys\": [ \"\" ] }")

	tests := []struct {
		name string
		want RequestApiKeys
	}{
		{
			name: "#2: ENV_REQUEST_API_KEYS not set",
			want: RequestApiKeys{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetApiKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetApiKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetApiKeys_3(t *testing.T) {
	os.Setenv(ENV_REQUEST_API_KEYS, "{ \"api_keys\": [ ] }")
	tests := []struct {
		name string
		want RequestApiKeys
	}{
		{
			name: "#3: ENV_REQUEST_API_KEYS not set",
			want: RequestApiKeys{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetApiKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetApiKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetApiKeys_4(t *testing.T) {
	tests := []struct {
		name string
		want RequestApiKeys
	}{
		{
			name: "#4: ENV_REQUEST_API_KEYS not set",
			want: RequestApiKeys{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetApiKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetApiKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
