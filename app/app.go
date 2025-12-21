package app

import (
	"context"
	"player/player"
	"player/ui"
	"player/ui/cli"

	"github.com/faiface/beep"
)

type MusicApp struct {
	player player.Player
	cli    ui.Ui
}

func NewMusicCLIApp(format beep.Format, streamer beep.StreamSeekCloser) *MusicApp {
	p := player.New(streamer, format)

	return &MusicApp{
		player: p,
		cli:    cli.New(p),
	}
}

func (app *MusicApp) Init() error {
	if err := app.player.Init(); err != nil {
		return err
	}

	if err := app.cli.Init(); err != nil {
		return err
	}

	return nil
}

// wg sync.WaitGroup

func (app *MusicApp) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go app.cli.EventLoop(ctx, cancel)
	go app.cli.Loop(ctx)
	go app.player.Play(ctx, cancel)

	<-ctx.Done()
}

func (app *MusicApp) Destroy() {
	app.cli.Destroy()
}
