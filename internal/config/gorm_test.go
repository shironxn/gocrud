package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestDB_Connection(t *testing.T) {
	type fields struct {
		cfg *Config
	}

	cfg, err := NewConfig()
	require.NoError(t, err)

	tests := []struct {
		name    string
		fields  fields
		want    interface{}
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				cfg: cfg,
			},
			want:    &gorm.DB{},
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				cfg: &Config{
					Database: struct {
						Host string
						Port string
						Name string
						User string
						Pass string
					}{},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DB{
				cfg: tt.fields.cfg,
			}
			got, err := d.Connection()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, tt.want, got)
			}
		})
	}
}