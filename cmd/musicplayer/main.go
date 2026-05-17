package main

import (
	"fmt"
	"os"
	"github.com/lvu-8/go-musicplayer/app"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
)

func openAndDecode(path string) (beep.StreamSeekCloser, beep.Format, *os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, beep.Format{}, nil, fmt.Errorf("failed opening file: %v", err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		f.Close()
		return nil, beep.Format{}, nil, fmt.Errorf("failed decoding mp3: %v", err)
	}

	return streamer, format, f, nil
}

func playFile(path string) error {
	streamer, format, file, err := openAndDecode(path)
	if err != nil {
		return err
	}

	defer file.Close()
	defer streamer.Close()

	a := app.NewMusicCLIApp(format, streamer)
	if err := a.Init(); err != nil {
		return err
	}
	defer a.Destroy()

	a.Run()
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: musicplayer <path-to-mp3-file>")
		os.Exit(1)
	}

	if err := playFile(os.Args[1]); err != nil {
		fmt.Printf("Error playing file: %v\n", err)
		os.Exit(1)
	}
}
