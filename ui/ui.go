package ui

import "context"

type Ui interface {
	Destroy()
	Init() error
	EventLoop(ctx context.Context, cancel context.CancelFunc)
	Loop(ctx context.Context)
}
