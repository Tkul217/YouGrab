package internal

import (
	"fmt"
	"sync"
	"yougrab/pkg"

	"github.com/kkdai/youtube/v2"
)

func DownloadPlaylist(url, outputPath, quality string, audioOnly, videoOnly bool, concurrent int) error {
	client := youtube.Client{}
	playlist, err := client.GetPlaylist(url)
	if err != nil {
		return fmt.Errorf("Error fetching playlist %s: %v", url, err)
	}

	fmt.Printf("%s Playlist: %s (%d video)\n", pkg.Yellow("➤"), playlist.Title, len(playlist.Videos))

	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrent) // Limit concurrent downloads

	for _, video := range playlist.Videos {
		wg.Add(1)
		go func(videoID string) {
			defer wg.Done()
			sem <- struct{}{} // Acquire semaphore
			defer func() { <-sem }()

			videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
			if err := DownloadSingleVideo(videoURL, outputPath, quality, audioOnly, videoOnly, concurrent); err != nil {
				fmt.Printf("%s Error downloading video %s: %v\n", pkg.Red("✗"), videoID, err)
			}
		}(video.ID)
	}

	wg.Wait()
	fmt.Printf("%s Playlist downloaded successfully\n", pkg.Green("✓"))
	return nil
}
