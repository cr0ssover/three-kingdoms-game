package db

import (
	"testing"

	"github.com/cr0ssover/three-kingdoms-game/server/config"
)

func TestInitDB(t *testing.T) {
	tests := []struct {
		name string
		key  map[string]string
	}{
		{
			name: "idle<conn",
			key: map[string]string{
				"max_idle": "2",
				"max_conn": "10",
			},
		},
		{
			name: "idle=conn",
			key: map[string]string{
				"max_idle": "1",
				"max_conn": "1",
			},
		},
		{
			name: "idle>conn",
			key: map[string]string{
				"max_idle": "10",
				"max_conn": "1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.key {
				config.File.Section("mysql").DeleteKey(tt.key[k])
				config.File.Section("mysql").NewKey(k, v)
			}
			InitDB()
		})
	}
}
