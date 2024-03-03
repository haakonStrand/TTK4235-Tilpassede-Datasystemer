package music

import (
	"fmt"
	"os"
	"os/exec"
)

func PlayMusic(ch chan bool) {
	for {
		boolMusic := true
		<-ch
		boolMusic = !boolMusic
		switch boolMusic {
		case true:
			mp3FilePath := "music/elev_music.mp3"
			cmd := exec.Command("mpg123", "-vC", mp3FilePath)
			// Redirect standard output and error to the terminal
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// Run the command
			err := cmd.Run()
			if err != nil {
				panic(err)
			}
		case false:
			fmt.Println("in else")
			mp3FilePath := "music/elev_music_long.mp3"
			cmd := exec.Command("mpg123", "-vC", mp3FilePath)
			// Redirect standard output and error to the terminal
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// Run the command
			err := cmd.Run()
			if err != nil {
				panic(err)
			}
		}

	}

}
