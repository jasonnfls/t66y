package main

import (
    "net"
    "net/http"
    "log"
    "fmt"
    "strings"
    "github.com/PuerkitoBio/goquery"
    "common"
    P "page"
    "time"
    "flag"
)

var untilDate string
var fromDate string
var debug bool
var page int

func init() {
    today := time.Now().Format("2006-01-02")
    flag.StringVar(&untilDate, "until", today, "Until")
    flag.StringVar(&fromDate, "from", today, "From")
    flag.BoolVar(&debug, "debug", false, "debug")
    flag.IntVar(&page, "page", 1, "starting page")
    flag.Parse()
    if fromDate < untilDate {
        log.Fatal("Wrong date range", fromDate, untilDate)
    }
}

var netClient = &http.Client{
    Timeout: time.Second * 60,
    Transport: &http.Transport{
        Dial: (&net.Dialer{
            Timeout: 60 * time.Second,
        }).Dial,
        TLSHandshakeTimeout: 60 * time.Second,
    },
}

type Post struct {
    tag string
    title string
    url string
    date string
}

func listingPage(page int) []Post{
    var url string
    if page == 1 {
        url = "http://www.t66y.com/thread0806.php?fid=8"
    } else {
        url = fmt.Sprintf("http://www.t66y.com/thread0806.php?fid=8&search=&page=%d", page)
    }
    request, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }
    request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:52.0) Gecko/20100101 Firefox/52.0")
    request.Header.Add("Accept-Encoding", "gzip")
    response, err := netClient.Do(request)
    if err != nil {
        log.Fatal(err)
    }

    doc, err := goquery.NewDocumentFromResponse(response)
    if err != nil {
        log.Fatal(err)
    }

    result := make([]Post, 0, 100)

    doc.Find("div#main div.t").Last().Find("tbody tr").Each(func(i int, sel *goquery.Selection){
        tal := sel.Find("td.tal")
        tag := strings.TrimSpace(common.DecodegbkStr(tal.Contents().Eq(0).Text()))
        if strings.HasPrefix(tag, "[") && strings.HasSuffix(tag, "]"){
            tag = strings.TrimSuffix(strings.TrimPrefix(tag, "["), "]")
            a := tal.Find("h3 a")
            title := common.DecodegbkStr(a.Text())
            url,_ := a.Attr("href")
            date := sel.Children().Eq(2).Find(".f10").Text()
            result = append(result, Post{
                tag: tag,
                title: title,
                url : url,
                date: date,
            })
        }

    })
    return result
}

var keywords = []string {
    "put", "your", "nauty", "keywords", "here", "utf8中文也可以"
}

func filterPost (post *Post) bool {
    if post.tag == "歐美" { // say you dont like 欧美 style
        return false
    }

    for i := range keywords {
        if strings.Contains(post.title, keywords[i]) {
            return true
        }
    }

    return false
}

func main(){
    limiter := time.Tick(time.Second * 6)
    outer:
    for ;page<100;page++ {
        if debug {
            fmt.Println("page:",page)
        }
        for _, post := range listingPage(page) {
            if debug {
                fmt.Println(post.date, post.title)
            }
            if fromDate < post.date {
                continue
            }
            if post.date < untilDate {
                break outer
            }
            if filterPost(&post) {
                if debug {
                    fmt.Println("selected: ", post)
                } else {
                    url := fmt.Sprintf("http://www.t66y.com/%s", post.url)
                    P.Crawl(url, post.title, post.date)
                }
            }
        }
        <-limiter
    }
}

