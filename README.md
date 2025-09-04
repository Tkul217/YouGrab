# YouGrab

**YouGrab** is a fast, user-friendly command-line tool written in Go for downloading videos and audio from YouTube. Featuring a sleek CLI with progress bars, color-coded output, and interactive quality selection, it supports playlists and high-resolution video merging with `ffmpeg`. Ideal for developers and users seeking a lightweight, customizable YouTube downloader.

## Features
- **Download Videos and Audio**: Download full videos (with audio), audio-only (MP3), or video-only streams.
- **Playlist Support**: Download entire YouTube playlists with configurable concurrent downloads.
- **Quality Selection**: Choose video resolution (e.g., 720p, 1080p) or audio bitrate (e.g., 128kbps, 160kbps) via flags or an interactive menu.
- **High-Quality Merging**: Merge high-resolution video and audio streams using `ffmpeg` for 1080p+ content.
- **Beautiful CLI**: Includes progress bars, color-coded output (green for success, red for errors, yellow for progress), and an interactive quality selection menu.
- **Concurrent Downloads**: Optimize playlist downloading with adjustable concurrency limits (default: 3).
- **Cross-Platform**: Compatible with Windows, macOS, and Linux.

## Installation
1. Ensure [Go](https://golang.org/dl/) (1.16+) and [ffmpeg](https://ffmpeg.org/download.html) are installed:
    - **Ubuntu**: `sudo apt install ffmpeg`
    - **macOS**: `brew install ffmpeg`
    - **Windows**: Download from [ffmpeg.org](https://ffmpeg.org/download.html) and add to PATH.
2. Clone the repository:
   ```bash
   git clone https://github.com/<your-username>/yougrab.git
   cd yougrab
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Build the tool:
   ```bash
   go build -o yougrab ./cmd
   ```

## Usage
```bash
# Download a video with audio
./yougrab -url "https://www.youtube.com/watch?v=VIDEO_ID" -quality 1080p

# Download audio only
./yougrab -url "https://www.youtube.com/watch?v=VIDEO_ID" -audio-only

# Download video only (no audio)
./yougrab -url "https://www.youtube.com/watch?v=VIDEO_ID" -video-only -quality 1080p

# Download a playlist
./yougrab -url "https://www.youtube.com/playlist?list=PLAYLIST_ID" -playlist -concurrent 5

# Interactive quality selection (omit -quality)
./yougrab -url "https://www.youtube.com/watch?v=VIDEO_ID"

# Specify output folder
./yougrab -url "https://www.youtube.com/watch?v=VIDEO_ID" -output ./my_downloads
```

See all options:
```bash
./yougrab --help
```

## Example Output
```
➤ Downloading: Sample Video Title
Select quality:
  * Video: video/mp4; codecs="avc1.640028", res: 1080p, fps: 30, ~150.5 MB (need a merge)
    Video: video/mp4; codecs="avc1.4d401f", res: 720p, fps: 30, ~80.2 MB
    Video: video/mp4; codecs="avc1.4d401e", res: 360p, fps: 30, ~40.1 MB
✓ Selected: Video: video/mp4; codecs="avc1.640028", res: 1080p, fps: 30, ~150.5 MB (need a merge)
Downloading... Sample_Video_Title_video_temp.mp4 [===========>        ] 150.5MB/150.5MB
Downloading... Sample_Video_Title_audio_temp.mp3 [===========>        ] 10.2MB/10.2MB
✓ Video with audio downloaded and merged: downloads/Sample_Video_Title.mp4
```

## Requirements
- Go 1.16 or higher
- `ffmpeg` for merging high-quality video and audio streams
- Internet connection for downloading

## Notes
- **Subtitles**: Not currently supported due to limitations in the `kkdai/youtube` library. For subtitle support, consider using `yt-dlp` or request an integration by opening an issue.
- **YouTube API**: The tool relies on `kkdai/youtube`, which may require updates for YouTube API changes. Run `go get -u github.com/kkdai/youtube/v2` to stay current.
- **High-Quality Videos**: Videos in 1080p+ often require separate video and audio streams, which are merged using `ffmpeg`. Ensure `ffmpeg` is installed for this feature.
- **Concurrency**: Use the `-concurrent` flag to adjust the number of simultaneous downloads for playlists. Higher values may trigger YouTube rate limits.

## License
This project is licensed under the [MIT License](LICENSE).

## Acknowledgments
- Built with [kkdai/youtube](https://github.com/kkdai/youtube) for YouTube API interactions.
- Uses [schollz/progressbar](https://github.com/schollz/progressbar) for progress bars, [fatih/color](https://github.com/fatih/color) for colored output, and [manfoldco/promptui](https://github.com/manfoldco/promptui) for interactive prompts.