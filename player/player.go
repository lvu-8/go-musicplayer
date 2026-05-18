package player

import (
	"time"

	"context"

	beep "github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

func clamp(x, min, max int) int {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

type AudioPlayer struct {
	streamer beep.StreamSeekCloser
	format   beep.Format
	ctrl     *beep.Ctrl
}

func New(streamer beep.StreamSeekCloser, format beep.Format) *AudioPlayer {
	return &AudioPlayer{
		streamer: streamer,
		format:   format,
		ctrl: &beep.Ctrl{Streamer: streamer, Paused: false},
	}
}

func (p *AudioPlayer) Init() error {
	return speaker.Init(
		p.format.SampleRate,
		p.format.SampleRate.N(time.Second/10),
	)
}

func (p *AudioPlayer) Play(ctx context.Context, cancel context.CancelFunc) {
	speaker.Play(
		beep.Seq(
			p.ctrl,
			beep.Callback(func() {
				cancel()
			}),
		),
	)
}

func (p *AudioPlayer) TogglePause() bool {
	speaker.Lock()
	p.ctrl.Paused = !p.ctrl.Paused
	speaker.Unlock()
	
	return p.ctrl.Paused
}

func (p *AudioPlayer) Restart() {
	speaker.Lock()
	p.streamer.Seek(0)
	speaker.Unlock()
}

func (p *AudioPlayer) SkipInMillisecond(milliseconds int) {
	speaker.Lock()
	defer speaker.Unlock()

	delta := p.format.SampleRate.N(time.Duration(milliseconds) * time.Millisecond)
	nextPosition := p.streamer.Position() + delta

	nextPosition = clamp(nextPosition, 0, p.streamer.Len())

	p.streamer.Seek(nextPosition)
}

func (p *AudioPlayer) GetCurrentMillisecond() int {
    speaker.Lock()
    defer speaker.Unlock()

    return int(p.format.SampleRate.D(p.streamer.Position()) / time.Millisecond)
}

func (p *AudioPlayer) GetLengthInMilliseconds() int {
    return int(p.format.SampleRate.D(p.streamer.Len()) / time.Millisecond)
}

func (p *AudioPlayer) Close() {
	speaker.Clear()
	p.streamer.Close()
}
