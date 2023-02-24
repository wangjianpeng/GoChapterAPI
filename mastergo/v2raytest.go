package main

import (
    "fmt"
    "time"
    "net"
)

func main() {
    urls := []string{"trojan://example.com:443", "vmess://example.com:8888", "ss://example.com:8888"}
    for _, url := range urls {
        start := time.Now()
        conn, err := net.DialTimeout("tcp", url, time.Second*5)
        if err != nil {
            fmt.Println(url, "connection error:", err)
            continue
        }
        defer conn.Close()

        fmt.Println(url, "connection time:", time.Since(start))
    }
}


