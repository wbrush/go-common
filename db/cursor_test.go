package db

import "testing"

func TestEncodeIdToCursor(t *testing.T) {
	tests := []struct {
		name string
		id   int64
		want string
	}{
		{
			name: "encode cursor",
			id:   42,
			want: "NDI=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeIdToCursor(tt.id); got != tt.want {
				t.Errorf("EncodeIdToCursor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeCursorToId(t *testing.T) {
	tests := []struct {
		name    string
		cursor  string
		want    int64
		wantErr bool
	}{
		{
			name:    "decode good cursor case",
			cursor:  "NDI=",
			want:    42,
			wantErr: false,
		},
		{
			name:    "decode wrong cursor case",
			cursor:  "wrong",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeCursorToId(tt.cursor)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeCursorToId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecodeCursorToId() = %v, want %v", got, tt.want)
			}
		})
	}
}
