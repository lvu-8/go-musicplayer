package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"player/player"
	"strings"
	"time"

	"golang.org/x/term"
)

const (
	PROGRESS_BAR_SIZE = 40

	QUITE    = 'q'
	PAUSED   = 'p'
	FORWARD  = 'x'
	BACKWARD = 'z'

	PROGRESS_BAR_SPACE = 1
	BUTTONS_SPACE      = 2
	MESSAGE_SPACE      = 3
	ENDING_LINE        = 4

	SPACES          = "       "
	QUITE_MESSAGE   = "- Quitting... -" + SPACES
	PAUSED_MESSAGE  = "- Paused -" + SPACES
	PLAYING_MESSAGE = "- Playing -" + SPACES
)

var BUTTONS_MASSAGE = fmt.Sprintf(
	"[%c] Quit, [%c] Pause/Resume [%c] Forward, [%c] Backward",
	QUITE, PAUSED, FORWARD, BACKWARD,
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

func (c *CLI) Init() {
	var err error

	c.oldState, err = term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		log.Fatal(err)
	}
}

func (c *CLI) Destroy() {
	fmt.Print(strings.Repeat("\n", ENDING_LINE-1) + "\r")

	term.Restore(int(os.Stdin.Fd()), c.oldState)
}

func (c *CLI) EventLoop(ctx context.Context, cancel context.CancelFunc) {
	buffer := make([]byte, 1)

	for c.processEvent(buffer) {
		select {
		case <-ctx.Done():
			return
		default:
			_, err := os.Stdin.Read(buffer)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	cancel()
}

func (c *CLI) printIntoLine(message string, n int) {
	s := strings.Repeat("\n", n-1) + "\r " +
		message +
		strings.Repeat("\033[A", n)

	fmt.Println(s)
}

func (c *CLI) processEvent(buffer []byte) bool {
	switch buffer[0] {
	case QUITE:
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

	c.printIntoLine(BUTTONS_MASSAGE, BUTTONS_SPACE)
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
	played := int(position / length * PROGRESS_BAR_SIZE)
	inverse := PROGRESS_BAR_SIZE - played

	bar := "[" +
		strings.Repeat("#", played) + strings.Repeat(" ", inverse) +
		"]"

	s := fmt.Sprintf("%s %.1f/%.1f sec", bar, position*0.001, length*0.001)
	c.printIntoLine(s, PROGRESS_BAR_SPACE)
}
