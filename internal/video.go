package internal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"yougrab/pkg"

	"github.com/kkdai/youtube/v2"
	"github.com/schollz/progressbar/v3"
)

func DownloadStream(client youtube.Client, video *youtube.Video, format *youtube.Format, outputPath, filename string) error {
	fmt.Printf("%s Starting download: %s\n", pkg.Yellow("➤"), filename)
	filePath := filepath.Join(outputPath, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	defer file.Close()

	reader, size, err := client.GetStream(video, format)
	if err != nil {
		return fmt.Errorf("error fetching stream: %v", err)
	}
	defer reader.Close()

	bar := progressbar.NewOptions64(
		size,
		progressbar.OptionSetDescription(fmt.Sprintf("Downloading %s", filename)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
	)

	teeReader := io.TeeReader(reader, bar)
	_, err = io.Copy(file, teeReader)
	if err != nil {
		return fmt.Errorf("download error: %v", err)
	}
	fmt.Printf("%s Download completed: %s\n", pkg.Green("✓"), filePath)
	return nil
}

func MergeVideoAudio(videoPath, audioPath, outputPath string) error {
	fmt.Printf("%s Merging video and audio into %s\n", pkg.Yellow("➤"), outputPath)
	// Use -shortest to avoid duration mismatches
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-i", audioPath, "-c:v", "copy", "-c:a", "copy", "-map", "0:v:0", "-map", "1:a:0", "-shortest", outputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("merge error: %v", err)
	}
	fmt.Printf("%s Merge completed\n", pkg.Green("✓"))
	return nil
}

func DownloadSingleVideo(url, outputPath, quality string, audioOnly, videoOnly bool, concurrent int) error {
	client := youtube.Client{}
	video, err := client.GetVideo(url)
	if err != nil {
		return fmt.Errorf("error fetching video %s: %v", url, err)
	}

	fmt.Printf("%s Downloading: %s\n", pkg.Yellow("➤"), video.Title)

	// Clean filename
	safeTitle := strings.ReplaceAll(video.Title, "/", "_")
	safeTitle = strings.ReplaceAll(safeTitle, "\\", "_")

	// Get stream sizes and debug all formats
	fmt.Printf("%s All available formats for video:\n", pkg.Yellow("➤"))
	for _, format := range video.Formats {
		fmt.Printf("  Format: %s, QualityLabel: %s, AudioChannels: %d, Bitrate: %d, FPS: %d\n",
			format.MimeType, format.QualityLabel, format.AudioChannels, format.Bitrate, format.FPS)
	}

	// Get stream sizes
	sizeMap := make(map[string]int64)
	for _, format := range video.Formats {
		_, size, err := client.GetStream(video, &format)
		if err == nil {
			sizeMap[format.QualityLabel] = size
		} else {
			fmt.Printf("%s Error getting size for format %s: %v\n", pkg.Red("✗"), format.QualityLabel, err)
		}
	}

	var stream youtube.Format
	var needsMerge bool
	var filename string

	if audioOnly {
		streams := video.Formats.WithAudioChannels()
		if len(streams) == 0 {
			return fmt.Errorf("audio streams not found")
		}
		if quality == "" {
			stream, _, err = SelectQuality(streams, true, sizeMap)
			if err != nil {
				return err
			}
		} else {
			for _, format := range streams {
				if format.QualityLabel == quality {
					stream = format
					break
				}
			}
			if stream.QualityLabel == "" {
				fmt.Printf("%s Audio with quality %s not found, selecting the best\n", pkg.Yellow("!"), quality)
				for _, format := range streams {
					if stream.QualityLabel == "" || format.Bitrate > stream.Bitrate {
						stream = format
					}
				}
			}
		}
		filename = safeTitle + ".mp3"
	} else if videoOnly {
		streams := video.Formats.Type("video")
		var filtered []youtube.Format
		for _, f := range streams {
			if f.AudioChannels == 0 {
				filtered = append(filtered, f)
			}
		}
		if len(filtered) == 0 {
			return fmt.Errorf("video streams without audio not found")
		}
		if quality == "" {
			stream, _, err = SelectQuality(filtered, false, sizeMap)
			if err != nil {
				return err
			}
		} else {
			for _, format := range filtered {
				if format.QualityLabel == quality {
					stream = format
					break
				}
			}
			if stream.QualityLabel == "" {
				fmt.Printf("%s Video with quality %s not found, selecting the best\n", pkg.Yellow("!"), quality)
				for _, format := range filtered {
					if stream.QualityLabel == "" || format.Height > stream.Height {
						stream = format
					}
				}
			}
		}
		filename = safeTitle + ".mp4"
	} else {
		// Use all video streams, prefer progressive if available
		streams := video.Formats.Type("video")
		if len(streams) == 0 {
			return fmt.Errorf("video streams not found")
		}

		fmt.Printf("%s Found %d video streams\n", pkg.Yellow("➤"), len(streams))

		if quality == "" {
			stream, needsMerge, err = SelectQuality(streams, false, sizeMap)
			if err != nil {
				return err
			}
		} else {
			for _, format := range streams {
				if format.QualityLabel == quality {
					stream = format
					needsMerge = format.AudioChannels == 0
					break
				}
			}
			if stream.QualityLabel == "" {
				fmt.Printf("%s Видео с качеством %s не найдено, выбираем лучшее\n", pkg.Yellow("!"), quality)
				for _, format := range streams {
					if stream.QualityLabel == "" || format.Height > stream.Height {
						stream = format
						needsMerge = format.AudioChannels == 0
					}
				}
			}
		}
		filename = safeTitle + ".mp4"
	}

	fmt.Printf("%s Selected format: %s (needsMerge: %v)\n", pkg.Yellow("➤"), stream.QualityLabel, needsMerge)

	if needsMerge {
		// Download video
		videoFile := filepath.Join(outputPath, safeTitle+"_video_temp.mp4")
		if err := DownloadStream(client, video, &stream, outputPath, filepath.Base(videoFile)); err != nil {
			return err
		}

		// Download audio
		var audioStream youtube.Format
		audioStreams := video.Formats.WithAudioChannels()
		fmt.Printf("%s Available audio streams:\n", pkg.Yellow("➤"))
		for _, format := range audioStreams {
			fmt.Printf("  Audio: %s, bitrate: %d kbps\n", format.MimeType, format.Bitrate/1000)
		}
		if len(audioStreams) == 0 {
			return fmt.Errorf("audio streams not found for merging")
		}
		for _, format := range audioStreams {
			if audioStream.QualityLabel == "" || format.Bitrate > audioStream.Bitrate {
				audioStream = format
			}
		}
		fmt.Printf("%s Selected audio format: %s\n", pkg.Yellow("➤"), audioStream.QualityLabel)
		audioFile := filepath.Join(outputPath, safeTitle+"_audio_temp.mp3")
		if err := DownloadStream(client, video, &audioStream, outputPath, filepath.Base(audioFile)); err != nil {
			return err
		}

		// Merge with ffmpeg
		finalFile := filepath.Join(outputPath, filename)
		if err := MergeVideoAudio(videoFile, audioFile, finalFile); err != nil {
			return fmt.Errorf("merge error: %v", err)
		}

		// Remove temporary files
		if err := os.Remove(videoFile); err != nil {
			fmt.Printf("%s Error removing temporary file %s: %v\n", pkg.Red("✗"), videoFile, err)
		}
		if err := os.Remove(audioFile); err != nil {
			fmt.Printf("%s Error removing temporary file %s: %v\n", pkg.Red("✗"), audioFile, err)
		}
		fmt.Printf("%s Video with audio downloaded and merged: %s\n", pkg.Green("✓"), finalFile)
	} else {
		// Simple download
		if err := DownloadStream(client, video, &stream, outputPath, filename); err != nil {
			return err
		}
		fmt.Printf("%s Успешно скачано: %s\n", pkg.Green("✓"), filepath.Join(outputPath, filename))
	}

	return nil
}
