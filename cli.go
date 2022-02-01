package main

import (
	"errors"
	"fmt"
	"github.com/bogem/id3v2"
	"github.com/winterssy/ghttp"
	"github.com/winterssy/glog"
	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/mxget/pkg/utils"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
)

func ConcurrentDownload(ctx context.Context, client api.Provider, savePath string, s *api.Song) (string, error) {
	songInfo := fmt.Sprintf("%s - %s", s.Artist, s.Name)
	if s.ListenURL == "" {
		glog.Errorf("Download [%s] failed: song unavailable", songInfo)
		return "", errors.New("song unavailable")
	}

	filePath := filepath.Join(savePath, utils.TrimInvalidFilePathChars(songInfo))
	glog.Infof("Start download [%s]", songInfo)
	mp3FilePath := filePath + ".mp3"

	req, _ := ghttp.NewRequest(ghttp.MethodGet, s.ListenURL)
	req.SetContext(ctx)
	resp, err := client.SendRequest(req)
	if err == nil {
		err = resp.SaveFile(mp3FilePath, 0664)
	}
	if err != nil {
		glog.Errorf("Download [%s] failed: %v", songInfo, err)
		_ = os.Remove(mp3FilePath)
		return "", errors.New("Download failed")
	}
	glog.Infof("Download [%s] complete", songInfo)
	glog.Infof("Update music metadata: [%s]", songInfo)
	writeTag(ctx, client, mp3FilePath, s)
	return mp3FilePath, nil
}

func writeTag(ctx context.Context, client api.Provider, filePath string, song *api.Song) {
	tag, err := id3v2.Open(filePath, id3v2.Options{Parse: true})
	if err != nil {
		return
	}
	defer tag.Close()

	tag.SetDefaultEncoding(id3v2.EncodingUTF8)
	tag.SetTitle(song.Name)
	tag.SetArtist(song.Artist)
	tag.SetAlbum(song.Album)

	if song.Lyric != "" {
		uslt := id3v2.UnsynchronisedLyricsFrame{
			Encoding:          id3v2.EncodingUTF8,
			Language:          "eng",
			ContentDescriptor: song.Name,
			Lyrics:            song.Lyric,
		}
		tag.AddUnsynchronisedLyricsFrame(uslt)
	}

	if song.PicURL != "" {
		req, _ := ghttp.NewRequest(ghttp.MethodGet, song.PicURL)
		req.SetContext(ctx)
		resp, err := client.SendRequest(req)
		if err == nil {
			pic, err := resp.Content()
			if err == nil {
				picFrame := id3v2.PictureFrame{
					Encoding:    id3v2.EncodingUTF8,
					MimeType:    "image/jpeg",
					PictureType: id3v2.PTFrontCover,
					Description: "Front cover",
					Picture:     pic,
				}
				tag.AddAttachedPicture(picFrame)
			}
		}
	}

	_ = tag.Save()
}
