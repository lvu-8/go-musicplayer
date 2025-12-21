package player

import "context"

const DEFAULT_SKIPPING = 3000 // in milliseconds

type Player interface {
	GetCurrentMillisecond() int
	GetLengthInMilliseconds() int
	TogglePause() bool
	SkipInMillisecond(milliseconds int)
	Init() error
	Play(ctx context.Context, cancel context.CancelFunc)
}
