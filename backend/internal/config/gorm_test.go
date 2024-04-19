package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestDB_Connection(t *testing.T) {
	type fields struct {
		config *Config
	}

	config, err := NewConfig()
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
				config: config,
			},
			want:    &gorm.DB{},
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				config: &Config{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DB{
				config: tt.fields.config,
			}
			got, err := d.Connection()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, tt.want, got)
				assert.NotNil(t, got)
			}
		})
	}
}
