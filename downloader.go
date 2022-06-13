package main

import (
	"fmt"
	"github.com/grafov/m3u8"
	"github.com/urfave/cli"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

func absolutize(rawURL string, u *url.URL) (uri *url.URL, err error) {

	subURL := rawURL
	uri, err = u.Parse(subURL)
	if err != nil {
		return
	}

	if rawURL == u.String() {
		return
	}

	if !uri.IsAbs() { // relative URI
		if rawURL[0] == '/' { // from the root
			subURL = fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, rawURL)
		} else { // last element
			split := strings.Split(u.String(), "/")
			split[len(split)-1] = rawURL

			subURL = strings.Join(split, "/")
		}
	}

	subURL, err = url.QueryUnescape(subURL)
	if err != nil {
		return
	}

	uri, err = u.Parse(subURL)
	if err != nil {
		return
	}
	return
}

type Downloader struct {
	*url.URL
	client *http.Client
	strict bool
	ctx    *cli.Context

	outPath string
}

func NewDownloader(u string) (*Downloader, error) {
	theURL, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("cms18> %v", err.Error())
	}

	client := &http.Client{}
	return &Downloader{
		URL:    theURL,
		client: client,
		strict: true,
	}, nil
}

func (d *Downloader) getContent(u *url.URL) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("cms1> %v", err.Error())
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cms2> %v", err.Error())
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Received HTTP %v for %v\n", resp.StatusCode, u.String())
	}

	return resp.Body, err
}

func (d *Downloader) writePlayList() error {
	fileName := path.Base(u.Path)
	out, err := os.Create(OUT_PATH + fileName)
	if err != nil {
		log.Fatal("cms3> " + err.Error())
	}
	defer out.Close()

	_, err = mpl.Encode().WriteTo(out)
	if err != nil {
		log.Fatal("cms4> " + err.Error())
	}
}

func (d *Downloader) getPlayList(u *url.URL) error {
	content, err := d.getContent(u)
	if err != nil {
		return err
	}

	playlist, listType, err := m3u8.DecodeFrom(content, d.strict)
	if err != nil {
		return err
	}

	if listType != m3u8.MEDIA && listType != m3u8.MASTER {
		return fmt.Errorf("cms11> " + "Not a valid playlist")
	}

	switch listType {
	case m3u8.MASTER:
		masterpl := playlist.(*m3u8.MasterPlaylist)
		for k, variant := range masterpl.Variants {
			if variant != nil {
				msURL, err := absolutize(variant.URI, d.URL)
				if err != nil {
					return err
				}
				err = d.getPlayList(msURL)
				if err != nil {
					return err
				}
				fmt.Print("cms13> "+"Downloaded chunk list number ", k+1, "\n\n")
				//break
			}
		}
		writePlaylist(u, m3u8.Playlist(masterpl))
		log.Print("cms14> "+"Downloaded Master Playlist: ", path.Base(d.Path), "\n")
	case m3u8.MEDIA:
		break
	}
}
