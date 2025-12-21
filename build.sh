go build -o builds/linux/musicplayer ./cmd/musicplayer
GOOS=windows GOARCH=amd64 go build -o builds/windows/musicplayer.exe ./cmd/musicplayer
