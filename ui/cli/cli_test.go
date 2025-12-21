package cli

import (
	"context"
	"testing"
)

type mockPlayer struct {
	paused bool
}

func (m *mockPlayer) GetCurrentMillisecond() int    { return 5000 }
func (m *mockPlayer) GetLengthInMilliseconds() int  { return 10000 }
func (m *mockPlayer) SkipInMillisecond(seconds int) {}
func (m *mockPlayer) TogglePause() bool {
	m.paused = !m.paused
	return m.paused
}
func (m *mockPlayer) Init() error {
	return nil
}
func (m *mockPlayer) Play(ctx context.Context, cancel context.CancelFunc) {}

func TestProcessEventQuit(t *testing.T) {
	c := New(&mockPlayer{})
	buffer := []byte{'q'}
	result := c.processEvent(buffer)
	if result {
		t.Error("processEvent should return false on quit")
	}
}

func TestProcessEventPausePlay(t *testing.T) {
	m := &mockPlayer{}
	c := New(m)
	buffer := []byte{'p'}
	result := c.processEvent(buffer)
	if !result {
		t.Error("processEvent should return true on pause")
	}
	if !m.paused {
		t.Error("Player should be paused after 'p'")
	}
	// Toggle again
	result = c.processEvent(buffer)
	if !result {
		t.Error("processEvent should return true on play")
	}
	if m.paused {
		t.Error("Player should be playing after second 'p'")
	}
}

func TestProcessEventForwardBackward(t *testing.T) {
	m := &mockPlayer{}
	c := New(m)
	buffer := []byte{'z'}
	result := c.processEvent(buffer)
	if !result {
		t.Error("processEvent should return true on forward")
	}
	buffer = []byte{'x'}
	result = c.processEvent(buffer)
	if !result {
		t.Error("processEvent should return true on backward")
	}
}
func TestPrintProgressBar(t *testing.T) {
	c := New(&mockPlayer{})
	c.printProgressBar(5000, 10000)
}

func TestTogglePause(t *testing.T) {
	m := &mockPlayer{}
	c := New(m)
	c.player.TogglePause()
	if !m.paused {
		t.Error("Pause should be toggled to true")
	}
	c.player.TogglePause()
	if m.paused {
		t.Error("Pause should be toggled to false")
	}
}
