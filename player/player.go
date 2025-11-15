package player

import (
	"time"

	"context"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
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
	paused   bool
}

func New(streamer beep.StreamSeekCloser, format beep.Format) *AudioPlayer {
	return &AudioPlayer{
		streamer: streamer,
		format:   format,
		paused:   false,
	}
}

func (p *AudioPlayer) Play(ctx context.Context, cancel context.CancelFunc) {
	speaker.Play(
		beep.Seq(
			p.streamer,
			beep.Callback(func() {
				cancel()
			}),
		),
	)
}

func (p *AudioPlayer) TogglePause() bool {
	p.paused = !p.paused

	if p.paused {
		speaker.Lock()
	} else {
		speaker.Unlock()
	}

	return p.paused
}

func (p *AudioPlayer) Close() {
	p.streamer.Close()
}

func (p *AudioPlayer) Init() {
	speaker.Init(
		p.format.SampleRate,
		p.format.SampleRate.N(time.Second/10),
	)
}

func (p *AudioPlayer) Restart() {
	p.streamer.Seek(0)
}

func (p *AudioPlayer) SkipInMillisecond(milliseconds int) {
	if !p.paused {
		speaker.Lock()
		defer speaker.Unlock()
	}

	nextPosition := p.streamer.Position() + p.format.SampleRate.N(
		time.Duration(milliseconds)*time.Millisecond)

	nextPosition = clamp(nextPosition, 0, p.streamer.Len())

	p.streamer.Seek(nextPosition)
}

func (p *AudioPlayer) GetCurrentMillisecond() int {
	return int(float64(p.streamer.Position()) * 1000 / float64(p.format.SampleRate))
}

func (p *AudioPlayer) GetLengthInMilliseconds() int {
	return int(float64(p.streamer.Len()) * 1000 / float64(p.format.SampleRate))
}
