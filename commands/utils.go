package commands

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
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
	for _, host := range hosts {
		if checkIsHostAvailable(host) {
			return host
		}
	}
	log.Fatalf("❌ 无可用的 Host")
	return ""
}

func checkIsHostAvailable(host string) bool {
	bbs := host + SISPaths[INDEX]
	res, err := http.Get(bbs)
	if err != nil {
		log.Printf("❌ 访问 Host 地址出错：%v\n", err)
		return false
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s\n", res.StatusCode, res.Status)
		return false
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
		return false
	}
	txt := doc.Find("title").Text()
	txt = strings.Replace(txt, "  ", " ", -1)
	txt = strings.TrimSpace(txt)
	if strings.Contains(txt, "SiS001! Board -") {
		return true
	}
	return false
}

func createClient() *resty.Client {
	// 参考：https://www.loginradius.com/blog/engineering/tune-the-go-http-client-for-high-performance/
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	client := resty.New()
	client.SetDebug(false)
	client.SetTransport(t)
	client.SetRetryCount(3)
	client.SetTimeout(10 * time.Second)

	return client
}
