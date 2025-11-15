package main

import (
	"fmt"
	"log"
	"os"
	"player/app"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
)

func openMusic(path string) *os.File {
	f, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	return f
}

func decodeMusic(f *os.File) (streamer beep.StreamSeekCloser, format beep.Format) {
	streamer, format, err := mp3.Decode(f)

	if err != nil {
		log.Fatal(err)
	}

	return streamer, format
}

func playFile(path string) {
	f := openMusic(path)
	defer f.Close()

	streamer, format := decodeMusic(f)

	defer streamer.Close()

	a := app.NewMusicCLIApp(format, streamer)
	a.Init()
	defer a.Destroy()

	a.Run()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: musicplayer <path-to-mp3-file>")
		return
	}

	if !fileExists(os.Args[1]) {
		fmt.Printf("File not found: %s\n", os.Args[1])
		return
	}

	playFile(os.Args[1])
}
