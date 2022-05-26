package src

import (
	"github.com/anaskhan96/soup"
	"log"
	"strings"
	"time"
)

var hosts []string

type PageType string

const (
	NEW     PageType = "forum-561"
	ACG              = "forum-231"
	NOVEL            = "forum-383"
	WESTERN          = "forum-229"
	INDEX            = "-"
)

var SISPaths map[PageType]string

func init() {
	hosts = []string{
		// "https://sis001.com/",
		"http://154.84.6.38/",
		"http://162.252.9.11/",
		"http://154.84.5.249/",
		"http://154.84.5.211/",
		"http://162.252.9.2/",
		"http://68.168.16.150/",
		"http://68.168.16.151/",
		"http://68.168.16.153/",
		"http://68.168.16.154/",
		"http://23.225.255.95/",
		"http://23.225.255.96/",
		"https://pux.sisurl.com/",
		"http://23.225.172.96/",
	}

	SISPaths = map[PageType]string{
		INDEX:   "bbs",
		NEW:     "bbs/" + string(NEW),
		ACG:     "bbs/" + string(ACG),
		NOVEL:   "bbs/" + string(NOVEL),
		WESTERN: "bbs/" + string(WESTERN),
	}
}

func FindAvailableHost() string {
	ch := make(chan string, 1)
	for _, host := range hosts {
		bbs := host + SISPaths[INDEX]
		go func(h string) {
			resp, err := soup.Get(h)
			if err != nil {
				log.Printf("❌ 访问 Host 地址出错：%v\n", err)
				return
			}
			doc := soup.HTMLParse(resp)
			txt := doc.Find("title").Text()
			txt = strings.Replace(txt, "  ", " ", -1)
			txt = strings.TrimSpace(txt)
			if strings.Contains(txt, "SiS001! Board -") {
				ch <- h
			}
		}(bbs)
	}

	select {
	case host := <-ch:
		return host
	case <-time.After(3 * time.Second):
		log.Fatal("❌ 无可用的 Host 地址")
		return ""
	}
}
