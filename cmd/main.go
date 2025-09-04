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
	url := flag.String("url", "", "Ссылка на видео или плейлист YouTube")
	output := flag.String("output", "downloads", "Папка для сохранения файлов")
	quality := flag.String("quality", "", "Качество: для видео - '720p', '1080p'; для аудио - '128kbps', '160kbps'")
	audioOnly := flag.Bool("audio-only", false, "Скачивать только аудио")
	videoOnly := flag.Bool("video-only", false, "Скачивать только видео (без аудио)")
	playlist := flag.Bool("playlist", false, "Скачивать весь плейлист")
	concurrent := flag.Int("concurrent", 3, "Количество одновременных загрузок для плейлистов")
	flag.Parse()

	if *url == "" {
		fmt.Printf("%s Ошибка: требуется ссылка (--url)\n", pkg.Red("✗"))
		flag.Usage()
		os.Exit(1)
	}

	// Check ffmpeg if needed
	if !*audioOnly && !*videoOnly {
		if _, err := exec.LookPath("ffmpeg"); err != nil {
			fmt.Printf("%s Ошибка: ffmpeg не установлен. Установите ffmpeg для слияния видео и аудио.\n", pkg.Red("✗"))
			os.Exit(1)
		}
	}

	// Create output directory
	if err := os.MkdirAll(*output, 0755); err != nil {
		fmt.Printf("%s Ошибка создания папки %s: %v\n", pkg.Red("✗"), *output, err)
		os.Exit(1)
	}

	if *playlist {
		if err := services.DownloadPlaylist(*url, *output, *quality, *audioOnly, *videoOnly, *concurrent); err != nil {
			fmt.Printf("%s Ошибка: %v\n", pkg.Red("✗"), err)
			os.Exit(1)
		}
	} else {
		if err := services.DownloadSingleVideo(*url, *output, *quality, *audioOnly, *videoOnly, *concurrent); err != nil {
			fmt.Printf("%s Ошибка: %v\n", pkg.Red("✗"), err)
			os.Exit(1)
		}
	}
}
