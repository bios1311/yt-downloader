package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/kkdai/youtube/v2"
	"github.com/spf13/pflag"
)

func main() {
	playlis_url := pflag.String("urlplay", "", "Enter the url/link of playlist (DEFAULT:EMPTY)")
	video_url := pflag.String("vidurl", "", "Video url to be passed here(DEFALUT:EMPTY)")
	playlis := pflag.String("playlist", "", "Playlist id to passed(DEFALUT:EMPTY)")
	vid := pflag.String("video", "", "Video id to be passed here(DEFALUT:EMPTY)")
	pflag.Parse()

	os.Mkdir("./videos", 0777)
	os.Mkdir("./songs", 0777)
	if len(*playlis_url) > 0 || len(*video_url) > 0 || len(*video_url) > 0 || len(*vid) > 0 {
		if len(*playlis) != 0 {
			client, playlist := seekplaylist(*playlis)
			downpsimple(client, playlist)
		} else if len(*vid) != 0 {
			downloadvideo(*vid)
		} else if len(*playlis_url) != 0 {
			pid, err := ExtractPlayistId(*playlis_url)
			if err != nil {
				panic(err)
			}
			client, playlist := seekplaylist(pid)
			downpsimple(client, playlist)
		} else if len(*video_url) != 0 {
			vidid, err := youtube.ExtractVideoID(*video_url)
			if err != nil {
				panic(err)
			}
			downloadvideo(vidid)
		} else {
			fmt.Println("Provide playlist id or video id or playlist url or video url:\nExample for video https://www.youtube.com/watch?v=1gEhUnwc8GA\n 1gEhUnwc8GA is id\nhttps://www.youtube.com/list=shkdbhksdbck\nshkdbhksdbck is playlist id")
		}
	} else {
		var play string
		fmt.Printf("Enter playlist url  or video url :\n")
		fmt.Scan(&play)

		_, err := ExtractPlayistId(play)
		if err != nil {
			id, err := youtube.ExtractVideoID(play)
			if err != nil {
				panic(err)

			}
			downloadvideo(id)

		} else {
			client, playlist := seekplaylist(play)
			//downpsimple(client, playlist)
			download_parallel(client, playlist)
		}
	}

	//downloadplaylist(client, playlist)

}
func ExtractPlayistId(url string) (string, error) {
	urlslice := strings.Split(url, "list=")
	if len(urlslice) != 2 {
		return "", fmt.Errorf("wrong url, playlist url will contain word list or playlist in between")
	}
	return urlslice[len(urlslice)-1], nil
}

func convert(inputVideo string, outputMP3 string, wg *sync.WaitGroup) {
	defer wg.Done()
	inputVideo = fmt.Sprint("./videos/", inputVideo)
	outputMP3 = fmt.Sprint("./songs/", outputMP3)

	cmd := exec.Command("./ffmpeg.exe", "-i", inputVideo, "-vn", "-acodec", "libmp3lame", outputMP3)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error converting video:", err)
	} else {
		fmt.Println("MP3 conversion successful!")
	}
}
func videoname(unformatted_name string) string {
	name_arr := strings.Split(unformatted_name, "|")
	var formatted_name string
	for _, val := range name_arr {
		formatted_name += val
	}
	name_arr = strings.Split(formatted_name, "/")
	formatted_name = ""
	for _, val := range name_arr {
		formatted_name += val
	}
	return formatted_name

}
func downloadvideo(videoID string) {
	client := youtube.Client{}
	wg := &sync.WaitGroup{}
	video, err := client.GetVideo(videoID)
	if err != nil {
		panic(err)
	}
	formats := video.Formats.WithAudioChannels()

	stream, _, err := client.GetStream(video, &formats[2])
	if err != nil {
		panic(err)
	}
	defer stream.Close()
	inpname := fmt.Sprint(videoname(video.Title), ".mp4")
	outpname := fmt.Sprint(videoname(video.Title), ".mp3")
	file, err := os.Create("./videos/" + inpname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}
	wg.Add(1)
	go convert(inpname, outpname, wg)
	wg.Wait()
}

// seek playlist prints all videos info that are in the playlist
func seekplaylist(playlistid string) (*youtube.Client, *youtube.Playlist) {

	client := youtube.Client{}

	playlist, err := client.GetPlaylist(playlistid)
	if err != nil {
		panic(err)
	}

	/* ----- Enumerating playlist videos ----- */
	header := fmt.Sprintf("Playlist %s by %s", playlist.Title, playlist.Author)
	println(header)
	println(strings.Repeat("=", len(header)) + "\n")

	for k, v := range playlist.Videos {
		fmt.Printf("(%d) %s - '%s'\n", k+1, v.Author, v.Title)
	}

	return &client, playlist
}
func downpsimple(client *youtube.Client, playlist *youtube.Playlist) {
	wg := sync.WaitGroup{}
	for _, PlaylistEntry := range playlist.Videos {

		video, err := client.VideoFromPlaylistEntry(PlaylistEntry)
		if err != nil {
			panic(err)
		}
		// Now it's fully loaded.

		fmt.Printf("Downloading %s by '%s'!\n", video.Title, video.Author)
		formats := video.Formats.WithAudioChannels()
		stream, _, err := client.GetStream(video, &formats[2])
		if err != nil {
			panic(err)
		}
		inpname := fmt.Sprint(videoname(PlaylistEntry.Title), ".mp4")
		outpname := fmt.Sprint(videoname(PlaylistEntry.Title), ".mp3")
		file, err := os.Create("./videos/" + inpname)

		if err != nil {
			log.Print(err)
		}

		defer file.Close()
		_, err = io.Copy(file, stream)

		if err != nil {
			log.Print(err)
		}

		println("Downloaded :" + inpname)
		wg.Add(1)
		go convert(inpname, outpname, &wg)
	}
	wg.Wait()

}

func worker(client *youtube.Client, playlist *youtube.Playlist, video_det_obj <-chan *video_det, done chan struct{}) {
	wg := sync.WaitGroup{}
	for video_ := range video_det_obj {

		video := video_.Video
		PlaylistEntry := video_.PlaylistEntry
		fmt.Printf("Downloading %s by '%s'!\n", video.Title, video.Author)
		formats := video.Formats.WithAudioChannels()
		stream, _, err := client.GetStream(video, &formats[2])
		if err != nil {
			panic(err)
		}
		inpname := fmt.Sprint(videoname(PlaylistEntry.Title), ".mp4")
		outpname := fmt.Sprint(videoname(PlaylistEntry.Title), ".mp3")
		file, err := os.Create("./videos/" + inpname)

		if err != nil {
			log.Print(err)
		}

		_, err = io.Copy(file, stream)

		if err != nil {
			log.Print(err)
		}
		wg.Add(1)
		go convert(inpname, outpname, &wg)
		println("Downloaded :" + inpname)
	}
	wg.Wait()
	done <- struct{}{}

}

func download_parallel(client *youtube.Client, playlist *youtube.Playlist) {
	video_chan := make(chan *video_det, 10)
	done := make(chan struct{})
	//worker pool creation with 5 workers
	for i := 0; i < 5; i++ {
		go worker(client, playlist, video_chan, done)
	}

	for _, PlaylistEntry := range playlist.Videos {

		video, err := client.VideoFromPlaylistEntry(PlaylistEntry)
		if err != nil {
			panic(err)
		}
		// Now it's fully loaded.
		video_chan <- &video_det{
			Video:         video,
			PlaylistEntry: PlaylistEntry,
		}
	}
	i := 0
	for {
		select {
		case <-done:
			log.Printf("Download:%v finished", i)
			i++
		}
	}

}

type video_det struct {
	Video         *youtube.Video
	PlaylistEntry *youtube.PlaylistEntry
}
