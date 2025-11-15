package ui

import "context"

type Ui interface {
	Destroy()
	Init()
	EventLoop(ctx context.Context, cancel context.CancelFunc)
	Loop(ctx context.Context)
}
