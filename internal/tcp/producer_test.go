package tcp

import "testing"

func TestHello(t *testing.T) {
	got := Broadcast()
	want := "Hello, world"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
