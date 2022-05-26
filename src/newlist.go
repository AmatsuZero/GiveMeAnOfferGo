package src

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"log"
	"path"
	"strconv"
	"strings"
)

type ThreadInfo struct {
	href, tag string
}

type PageParser interface {
	MaxPageSelector() string
	CurrentPageURL() string
	FindMaxPage() int
}

func (t *ThreadInfo) String() string {
	return fmt.Sprintf("【标签】%v 【链接】%v", t.tag, t.href)
}

type ThreadPage struct {
	currentPage, maxPage  int
	pageType, host, title string
}

func (p *ThreadPage) GetAllThreadsOnCurrentPage() {
	url := p.CurrentPageURL()
	resp, err := soup.Get(url)
	log.Printf("🔗 即将打开%v第%v页\n", p.title, p.currentPage)
	if err != nil {
		log.Fatal(err)
	}
	doc := soup.HTMLParse(resp)
	if p.maxPage == 0 {
		p.maxPage = p.FindMaxPage(doc)
	}
	doc = doc.Find("#wrapper > div:nth-child(1) > div.mainbox.threadlist > form")
	elms := doc.FindAll("tbody[id]")
	elms = FilterElements(elms, func(root soup.Root) bool {
		id := root.Attrs()["id"]
		return strings.HasPrefix(id, "normalthread_")
	})
	return
}

func FilterElements(elms []soup.Root, filter func(root soup.Root) bool) []soup.Root {
	tmp := elms[:0]
	for _, v := range elms {
		if filter != nil && filter(v) {
			tmp = append(tmp, v)
		}
	}
	return tmp
}

func (p *ThreadPage) CurrentPageURL() string {
	path := fmt.Sprintf("%v-%v.html", p.pageType, p.currentPage)
	return strings.Join([]string{p.host, path}, "/")
}

func (p *ThreadPage) MaxPageSelector() string {
	return "#wrapper > div:nth-child(1) > div:nth-child(9) > div > a.last"
}

func (p *ThreadPage) FindMaxPage(doc soup.Root) (page int) {
	link := doc.Find(p.MaxPageSelector()).Attrs()["href"]
	if len(link) > 0 {
		idx := strings.IndexRune(link, '/')
		link = link[idx:]
		ext := path.Ext(link)
		idx = strings.LastIndex(link, ext)
		link = link[idx:]
		page, _ = strconv.Atoi(link)
		return
	}
	return
}
