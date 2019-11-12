package loopapi

import (
	"testing"
)

func TestLoopEnergy_Connect(t *testing.T) {
	type fields struct {
		Config config
	}

	var cfg fields
	cfg.Config.LoopPort = 443
	cfg.Config.LoopServer = "https://www.your-loop.com"

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// Test cases.
		{"Simple Connect", cfg, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loopEng := &LoopEnergy{
				Config: tt.fields.Config,
			}
			if got := loopEng.Connect(); got != tt.want {
				t.Errorf("LoopEnergy.Connect() = %v, want %v", got, tt.want)
			}
		})
	}
}
