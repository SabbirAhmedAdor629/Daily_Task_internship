package main

import (
    "fmt"
    "math/rand"
    "strings"
    "time"
)







const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
    sb := strings.Builder{}
    sb.Grow(n)
    for i := 0; i < n; i++ {
        sb.WriteByte(charset[rand.Intn(len(charset))])
    }
    return sb.String()
}

func execute() {
    rand.Seed(time.Now().UnixNano())
    a := randomString(20)
    fmt.Println(a)
}

func main(){
    execute()
}
