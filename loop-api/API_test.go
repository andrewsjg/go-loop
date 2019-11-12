package loopapi

import "testing"

func TestConnect(t *testing.T) {
	want := true
	if got := Connect(); got != want {
		t.Errorf("Connect() = %t, want %t", got, want)
	}
}
