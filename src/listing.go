package main

import (
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

func init() {
    today := time.Now().Format("2006-01-02")
    flag.StringVar(&untilDate, "until", today, "Until")
    flag.StringVar(&fromDate, "from", today, "From")
    flag.BoolVar(&debug, "debug", false, "debug")
    flag.Parse()
    if fromDate < untilDate {
        log.Fatal("Wrong date range", fromDate, untilDate)
    }
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
    doc, err := goquery.NewDocument(url)
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
    outer:
    for page:=1;;page++ {
        if debug {
            fmt.Println("page:",page)
        }
        for _, post := range listingPage(page) {
            if fromDate < post.date {
                continue
            }
            if post.date < untilDate {
                break outer
            }
            if filterPost(&post) {
                if debug {
                    fmt.Println(post)
                } else {
                    url := fmt.Sprintf("http://www.t66y.com/%s", post.url)
                    P.Crawl(url, post.title, post.date)
                }
            }
        }
    }
}

