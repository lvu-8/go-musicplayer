package player

import (
	"testing"

	"github.com/faiface/beep"
)

type mockStreamer struct {
	pos int
	len int
}

func (m *mockStreamer) Stream(samples [][2]float64) (n int, ok bool) { return 0, false }
func (m *mockStreamer) Err() error                                   { return nil }
func (m *mockStreamer) Len() int                                     { return m.len }
func (m *mockStreamer) Position() int                                { return m.pos }
func (m *mockStreamer) Seek(p int) error {
	m.pos = p
	return nil
}
func (m *mockStreamer) Close() error { return nil }

func TestSkipInMillisecond(t *testing.T) {
	format := beep.Format{SampleRate: 1000}
	streamer := &mockStreamer{pos: 500, len: 1000}
	player := New(streamer, format)

	player.SkipInMillisecond(200)
	if streamer.pos != 700 {
		t.Errorf("Expected position 700, got %d", streamer.pos)
	}

	player.SkipInMillisecond(-800)
	if streamer.pos != 0 {
		t.Errorf("Expected position 0, got %d", streamer.pos)
	}

	player.SkipInMillisecond(2000)
	if streamer.pos != 1000 {
		t.Errorf("Expected position 999, got %d", streamer.pos)
	}
}

func TestGetCurrentMillisecond(t *testing.T) {
	format := beep.Format{SampleRate: 1000}
	streamer := &mockStreamer{pos: 500, len: 1000}
	player := New(streamer, format)

	ms := player.GetCurrentMillisecond()
	if ms != 500 {
		t.Errorf("Expected 500 ms, got %d", ms)
	}
}

func TestGetLengthInMilliseconds(t *testing.T) {
	format := beep.Format{SampleRate: 1000}
	streamer := &mockStreamer{pos: 0, len: 1000}
	player := New(streamer, format)

	ms := player.GetLengthInMilliseconds()
	if ms != 1000 {
		t.Errorf("Expected 1000 ms, got %d", ms)
	}
}

func TestRestart(t *testing.T) {
	format := beep.Format{SampleRate: 1000}
	streamer := &mockStreamer{pos: 500, len: 1000}
	player := New(streamer, format)

	player.Restart()
	if streamer.pos != 0 {
		t.Errorf("Expected position 0 after restart, got %d", streamer.pos)
	}
}

func TestTogglePause(t *testing.T) {
	format := beep.Format{SampleRate: 1000}
	streamer := &mockStreamer{pos: 0, len: 1000}
	player := New(streamer, format)

	paused := player.TogglePause()
	if !paused {
		t.Error("Expected paused to be true")
	}
	paused = player.TogglePause()
	if paused {
		t.Error("Expected paused to be false")
	}
}
