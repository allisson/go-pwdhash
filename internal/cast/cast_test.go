package cast

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertStringToUint32(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint32
		wantErr bool
	}{
		{name: "ok", input: "12345", want: 12345},
		{name: "nonNumeric", input: "abc", wantErr: true},
		{name: "negative", input: "-1", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertStringToUint32(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestConvertStringToUint8(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint8
		wantErr bool
	}{
		{name: "ok", input: "255", want: 255},
		{name: "nonNumeric", input: "hello", wantErr: true},
		{name: "overflow", input: "512", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertStringToUint8(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestConvertIntToUint32(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		want    uint32
		wantErr bool
	}{
		{name: "zero", input: 0, want: 0},
		{name: "positive", input: 65535, want: 65535},
		{name: "negative", input: -1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertIntToUint32(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
