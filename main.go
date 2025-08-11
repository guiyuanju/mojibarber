package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unicode/utf8"

	"golang.org/x/term"
)

func main() {
	faces := []string{
		"( ´･ω･)",
		"(　´･ω)",
		"( 　´･)",
		"( 　 ´)",
		"(     )",
		"(`　　)",
		"(･`   )",
		"(ω･`　)",
		"(･ω･` )",
		"(´･ω･`)",
	}

	fd := int(os.Stdout.Fd())
	width, height, err := term.GetSize(fd)
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
		return
	}

	column := width / (utf8.RuneCountInString(faces[0]) + 1) // +1 for the tab character

	if len(faces) > height {
		faces = faces[:height]
	} else {
		i := 0
		j := len(faces)
		for j < height {
			faces = append(faces, faces[i])
			i++
			j++
		}
	}

	length := len(faces)

	var merged []string
	for i := range length {
		var cur string
		for j := range length {
			idx := (i + j) % length
			for range column {
				cur += faces[idx] + "\t"
			}
			cur += "\r\n"
		}
		merged = append(merged, cur)
	}

	hideCursor()
	defer showCursor()

	alternateScreenMode()
	defer normalScreenMode()

	sigChan := make(chan os.Signal, 1)
	stopChan := make(chan struct{})
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		_ = <-sigChan
		close(stopChan)
	}()

	var i int
	for {
		select {
		case <-stopChan:
			return
		default:
			fmt.Print(merged[i])
			i = (i + 1) % len(merged)
			time.Sleep(50 * time.Millisecond)
			clear()
		}
	}
}

func clear() {
	fmt.Print("\033[H\033[2J")
}

func hideCursor() {
	fmt.Print("\x1b[?25l")
}

func showCursor() {
	fmt.Print("\x1b[?25h")
}

func alternateScreenMode() {
	fmt.Print("\033[?1049h")
}

func normalScreenMode() {
	fmt.Print("\033[?1049l")
}
