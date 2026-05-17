package cli

import (
	"context"
	"fmt"
	"os"
	"github.com/lvu-8/go-musicplayer/player"
	"strings"
	"time"

	"golang.org/x/term"
)

const (
	PROGRESS_BAR_SPACE = iota + 1
	BUTTONS_SPACE
	MESSAGE_SPACE
	ENDING_LINE

	PROGRESS_BAR_SIZE = 40

	QUIT     = 'q'
	PAUSED   = 'p'
	FORWARD  = 'x'
	BACKWARD = 'z'

	QUITE_MESSAGE   = "- Quitting... -"
	PAUSED_MESSAGE  = "- Paused -"
	PLAYING_MESSAGE = "- Playing -"

	HIDE_CURSOR = "\033[?25l"
	SHOW_CURSOR = "\033[?25h"
	CLEAR_LINE  = "\033[2K\r "
)

var BUTTONS_MESSAGE = fmt.Sprintf(
	"[%c] Quit, [%c] Pause/Resume [%c] Forward, [%c] Backward",
	QUIT, PAUSED, FORWARD, BACKWARD,
)

const tickerSize = 100 * time.Millisecond

type CLI struct {
	player   player.Player
	oldState *term.State
}

func New(player player.Player) *CLI {
	return &CLI{
		player:   player,
		oldState: nil,
	}
}

func (c *CLI) Init() error {
	var err error = nil

	c.oldState, err = term.MakeRaw(int(os.Stdin.Fd()))

	fmt.Print(HIDE_CURSOR)
	return err
}

func (c *CLI) Destroy() {
	fmt.Print(SHOW_CURSOR)

	fmt.Print(strings.Repeat("\n", ENDING_LINE-1) + "\r")

	if c.oldState != nil {
        term.Restore(int(os.Stdin.Fd()), c.oldState)
    }
}

func (c *CLI) EventLoop(ctx context.Context, cancel context.CancelFunc) {
	buffer := make([]byte, 1)

	for c.processEvent(buffer) {
		select {
		case <-ctx.Done():
			return
		default:
			os.Stdin.Read(buffer)
		}
	}

	cancel()
}

func (c *CLI) printIntoLine(message string, n int) {
	s := strings.Repeat("\n", n-1) + CLEAR_LINE +
		message +
		strings.Repeat("\033[A", n)

	fmt.Println(s)
}

func (c *CLI) processEvent(buffer []byte) bool {
	switch buffer[0] {
	case QUIT:
		c.printIntoLine(QUITE_MESSAGE, MESSAGE_SPACE)
		return false
	case PAUSED:
		if c.player.TogglePause() {
			c.printIntoLine(PAUSED_MESSAGE, MESSAGE_SPACE)
		} else {
			c.printIntoLine(PLAYING_MESSAGE, MESSAGE_SPACE)
		}
	case FORWARD:
		c.player.SkipInMillisecond(player.DEFAULT_SKIPPING)
	case BACKWARD:
		c.player.SkipInMillisecond(-player.DEFAULT_SKIPPING)
	}

	return true
}

func (c *CLI) Loop(ctx context.Context) {
	length := c.player.GetLengthInMilliseconds()
	ticker := time.NewTicker(tickerSize)
	defer ticker.Stop()

	c.printIntoLine(BUTTONS_MESSAGE, BUTTONS_SPACE)
	c.printIntoLine(PLAYING_MESSAGE, MESSAGE_SPACE)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			current := c.player.GetCurrentMillisecond()
			c.printProgressBar(float64(current), float64(length))
		}
	}
}

func (c *CLI) printProgressBar(position float64, length float64) {
	if length <= 0 {
		length = 1
	}

	played := int(position / length * PROGRESS_BAR_SIZE)
	inverse := PROGRESS_BAR_SIZE - played

	bar := "[" +
		strings.Repeat("#", played) + strings.Repeat(" ", inverse) +
		"]"

	status := fmt.Sprintf("%s %.1fs/%.1fs", bar, position*0.001, length*0.001)
	c.printIntoLine(status, PROGRESS_BAR_SPACE)
}
