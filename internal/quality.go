package internal

import (
	"fmt"
	"sort"
	"yougrab/pkg"

	"github.com/kkdai/youtube/v2"
	"github.com/manifoldco/promptui"
)

func SelectQuality(formats []youtube.Format, isAudio bool, sizeMap map[string]int64) (youtube.Format, bool, error) {
	fmt.Printf("%s Доступные форматы:\n", pkg.Yellow("➤"))
	for _, format := range formats {
		size := sizeMap[format.QualityLabel]
		if isAudio {
			fmt.Printf("  Аудио: %s, bitrate: %d kbps, ~%s\n", format.MimeType, format.Bitrate/1000, FormatSize(size))
		} else {
			fmt.Printf("  Видео: %s, res: %s, fps: %d, audio: %v, ~%s\n", format.MimeType, format.QualityLabel, format.FPS, format.AudioChannels > 0, FormatSize(size))
		}
	}

	var options []FormatOption
	for _, format := range formats {
		var label string
		size := sizeMap[format.QualityLabel]
		isMerged := format.AudioChannels == 0 && !isAudio
		if isAudio {
			label = fmt.Sprintf("Аудио: %s, bitrate: %d kbps, ~%s", format.MimeType, format.Bitrate/1000, FormatSize(size))
		} else {
			label = fmt.Sprintf("Видео: %s, res: %s, fps: %d, ~%s%s", format.MimeType, format.QualityLabel, format.FPS, FormatSize(size), pkg.Blue(" (требуется слияние)"))
		}
		options = append(options, FormatOption{Label: label, Format: format, Size: size, IsMerged: isMerged})
	}

	if len(options) == 0 {
		return youtube.Format{}, false, fmt.Errorf("нет доступных форматов")
	}

	// Sort: audio by bitrate, video by resolution
	if isAudio {
		sort.Slice(options, func(i, j int) bool {
			return options[i].Format.Bitrate > options[j].Format.Bitrate
		})
	} else {
		sort.Slice(options, func(i, j int) bool {
			return options[i].Format.Height > options[j].Format.Height
		})
	}

	prompt := promptui.Select{
		Label: "Выберите качество",
		Items: options,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}:",
			Active:   fmt.Sprintf("%s {{ .Label | underline }}", pkg.Green("*")),
			Inactive: "  {{ .Label }}",
			Selected: pkg.Green("Выбрано: {{ .Label }}"),
		},
	}

	i, _, err := prompt.Run()
	if err != nil {
		return youtube.Format{}, false, fmt.Errorf("ошибка выбора: %v", err)
	}
	return options[i].Format, options[i].IsMerged, nil
}
