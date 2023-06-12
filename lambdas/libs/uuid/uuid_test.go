package uuid

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestValidateUUID(t *testing.T) {
	t1, _ := uuid.Parse("c64e899f-8ad5-4952-9353-f3b4d07aa6ca")

	type args struct {
		src string
	}
	tests := []struct {
		name    string
		args    args
		want    *uuid.UUID
		wantErr bool
	}{
		{
			name: "Test UUID #1: c64e899f-8ad5-4952-9353-f3b4d07aa6ca",
			args: args{
				"c64e899f-8ad5-4952-9353-f3b4d07aa6ca",
			},
			want:    &t1,
			wantErr: false,
		},
		{
			name: "Test UUID #2: Trim leading space",
			args: args{
				src: " c64e899f-8ad5-4952-9353-f3b4d07aa6ca",
			},
			want:    &t1,
			wantErr: false,
		},
		{
			name: "Test UUID #3: Trim trailing space",
			args: args{
				src: "c64e899f-8ad5-4952-9353-f3b4d07aa6ca ",
			},
			want:    &t1,
			wantErr: false,
		},
		{
			name: "Test UUID #4: Trim leading and trailing spaces",
			args: args{
				src: "  c64e899f-8ad5-4952-9353-f3b4d07aa6ca  ",
			},
			want:    &t1,
			wantErr: false,
		},
		{
			name: "Test UUID: #5: Convert to lowercase",
			args: args{
				src: "C64E899F-8AD5-4952-9353-F3B4D07AA6CA",
			},
			want:    &t1,
			wantErr: false,
		},
		{
			name: "Test UUID: #6: Invalid UUID",
			args: args{
				src: "ABCDC64E899F-8AD5-4952-9353-F3B4D07AA6CA",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test UUID: #7: Invalid UUID",
			args: args{
				src: "ABCDC64E899F8AD5-4952-9353-F3B4D07AA6CA",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateUUID(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
