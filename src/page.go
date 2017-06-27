package main

import (
    "log"
    "os"
    "page"
)


func main(){
    if len(os.Args) < 2 {
        log.Fatal("Needs url")
    }
    page.Crawl(os.Args[1], "", "")
}

