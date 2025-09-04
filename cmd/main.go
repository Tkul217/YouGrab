package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	services "yougrab/internal"
	"yougrab/pkg"
)

func main() {
	url := flag.String("url", "", "Link to YouTube video or playlist")
	output := flag.String("output", "downloads", "Folder for saving files")
	quality := flag.String("quality", "", "Quality: for video - '720p', '1080p'; for audio - '128kbps', '160kbps'")
	audioOnly := flag.Bool("audio-only", false, "Download only audio (no video)")
	videoOnly := flag.Bool("video-only", false, "Download only video (no audio)")
	playlist := flag.Bool("playlist", false, "Download the entire playlist")
	concurrent := flag.Int("concurrent", 3, "Number of simultaneous downloads for playlists")
	flag.Parse()

	if *url == "" {
		fmt.Printf("%s Error: reequired link (--url)\n", pkg.Red("✗"))
		flag.Usage()
		os.Exit(1)
	}

	// Check ffmpeg if needed
	if !*audioOnly && !*videoOnly {
		if _, err := exec.LookPath("ffmpeg"); err != nil {
			fmt.Printf("%s Error: ffmpeg not installed. Install ffmpeg for merge video and audio.\n", pkg.Red("✗"))
			os.Exit(1)
		}
	}

	// Create output directory
	if err := os.MkdirAll(*output, 0755); err != nil {
		fmt.Printf("%s Error creating folder %s: %v\n", pkg.Red("✗"), *output, err)
		os.Exit(1)
	}

	if *playlist {
		if err := services.DownloadPlaylist(*url, *output, *quality, *audioOnly, *videoOnly, *concurrent); err != nil {
			fmt.Printf("%s Error: %v\n", pkg.Red("✗"), err)
			os.Exit(1)
		}
	} else {
		if err := services.DownloadSingleVideo(*url, *output, *quality, *audioOnly, *videoOnly, *concurrent); err != nil {
			fmt.Printf("%s Error: %v\n", pkg.Red("✗"), err)
			os.Exit(1)
		}
	}
}
