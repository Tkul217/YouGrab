package internal

import (
	"fmt"

	"github.com/kkdai/youtube/v2"
)

// FormatOption for interactive quality selection
type FormatOption struct {
	Label    string
	Format   youtube.Format
	Size     int64
	IsMerged bool // Requires merging
}

func FormatSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	}
	return fmt.Sprintf("%.2f GB", float64(size)/(1024*1024*1024))
}
