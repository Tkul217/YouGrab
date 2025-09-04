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
		return fmt.Errorf("ошибка получения плейлиста %s: %v", url, err)
	}

	fmt.Printf("%s Плейлист: %s (%d видео)\n", pkg.Yellow("➤"), playlist.Title, len(playlist.Videos))

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
				fmt.Printf("%s Ошибка при скачивании видео %s: %v\n", pkg.Red("✗"), videoID, err)
			}
		}(video.ID)
	}

	wg.Wait()
	fmt.Printf("%s Плейлист успешно скачан\n", pkg.Green("✓"))
	return nil
}
