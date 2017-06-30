package page

import (
    "net/http"
    "io/ioutil"
    "strings"
    "os"
    "fmt"
    "net"
    "log"
    "time"
    "sync"
    "github.com/PuerkitoBio/goquery"
    "common"
)

const root = "/volume1/photo/t66y/"

var wg sync.WaitGroup

var dirMutex sync.Mutex
var sem = make(chan int, 4)

var netClient = &http.Client{
    Timeout: time.Second * 60,
    Transport: &http.Transport{
        Dial: (&net.Dialer{
            Timeout: 60 * time.Second,
        }).Dial,
        TLSHandshakeTimeout: 60 * time.Second,
    },
}

func Crawl(url, title, date string) {
    doc, err := goquery.NewDocument(url)
    if err != nil {
        log.Fatal(err)
    }

    if title == "" {
        title = doc.Find("div#main div.t2 h4").Text()

        titleUTF8Bytes, err := common.Decodegbk([]byte(title))
        if err != nil {
            log.Fatal(err)
        }
        title = string(titleUTF8Bytes)
    }

    if date == "" {
        dateStr, err := doc.Find("div.tipad").Html()
        if err != nil {
            log.Fatal(err)
        }
        pos := strings.Index(dateStr, "Posted:")
        if pos <0 {
            log.Fatal("Could not find date")
        }
        pos = pos + 7
        date = dateStr[pos:pos+10]
    }


    parent := doc.Find("div.do_not_catch").First()
    parent.Find("input").Each(processInputNode(date, title))
    wg.Wait()
}

func processInputNode(date, title string) func(int, *goquery.Selection) {
    return func(i int, selection *goquery.Selection) {
        url, exists := selection.Attr("src")
        if exists {
            wg.Add(1)
            go download(date, title, url, i)
        }
    }
}

func mkdir(base string) {
    dirMutex.Lock()
    fi, err := os.Stat(base)
    if err != nil && os.IsNotExist(err) {
        log.Println("creating directory", base)
        err := os.Mkdir(base, 0775)
        if err != nil {
            log.Fatal(err)
        }
    } else {
        if !fi.IsDir() {
            log.Fatal(base, "is not a directory")
        }
    }
    dirMutex.Unlock()
}

func download(date, title, url string, i int) {
    defer wg.Done()
    tokens := strings.Split(url, "/")
    base := root + date + "/" + title

    mkdir(root + date)
    mkdir(base)

    filename := fmt.Sprintf("%s/%02d_%s", base, i, tokens[len(tokens)-1])

    if _, err := os.Stat(filename); os.IsNotExist(err) {
        sem <- 1
        log.Println("Downloading",url)
        defer func(){
            <-sem
        }()

        request, err := http.NewRequest("GET", url, nil)
        if err != nil {
            log.Fatal(err)
        }
        request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:52.0) Gecko/20100101 Firefox/52.0")
        response, err := netClient.Do(request)
        if err != nil {
            log.Println(err)
            return
        }

        defer response.Body.Close()
        content, err := ioutil.ReadAll(response.Body)

        if err != nil {
            log.Println(err)
            return
        }

        ioutil.WriteFile(filename, content, 0644)
    } else {
        log.Println(url, "already Downloaded")
    }
}

