package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "io/ioutil"
    "strconv"
    "net/http"
    "net/url"
)
// Multiple definition of flag containing int
// souce: https://lawlessguy.wordpress.com/2013/07/23/filling-a-slice-using-command-line-flags-in-go-golang/
type intslice []int

func (i *intslice) String() string {
    return fmt.Sprintf("%d", *i)
}

func (i *intslice) Set(value string) error {
    tmp, err := strconv.Atoi(value)
    if err != nil {
        *i = append(*i, -1)
    } else {
        *i = append(*i, tmp)
    }
    return nil
}

var myints intslice
// End def

func sendMessage(reciever int, message string, fullUrl string) {
    v := url.Values{}
    v.Set("chat_id", strconv.Itoa(reciever))
    v.Set("text", message)
    v.Set("disable_web_page_preview", "1")

    resp, err := http.PostForm(fullUrl, v)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    body , err := ioutil.ReadAll(resp.Body)
    fmt.Printf("Status: %s\n", resp.Status)
    respReadable := string(body[:])
    fmt.Println(respReadable)
}

var chatId intslice

func main() {
    flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
    flags.Var(&chatId, "id", "Telegram Recievers id")
    token := flags.String("token", "", "Telegram token")
    telegramUrl := flags.String("url", "https://api.telegram.org/", "Telegram Api url")
    message := flags.String("message", "", "Telegram Message")
    dryRun := flags.Bool("dry", false,"Test mode, nothing will be sent to twitter")
    test := flags.Bool("test", false,"Test mode, nothing will be sent to twitter")
    flags.Parse(os.Args[1:]) 

    // Check if message is present
    if *message == "" { log.Fatal("You must provide a message") }

    // Assemble fullUrl
    fullUrl := *telegramUrl + "bot" + *token + "/sendMessage"

    // Posting message
    for _, to := range chatId {
        if *dryRun {
            fmt.Printf("Token: %s\nURL: %s\nMessage: %s\n", *token, fullUrl, *message)
        } else if *test {
            log.Printf("Sending to httpbin.org/post for: %d\n", to)
            fullUrl = "https://httpbin.org/post"
            sendMessage(to, *message, fullUrl)
        } else {
            log.Printf("Sending to Telegram id: %d\n", to)
            sendMessage(to, *message, fullUrl)
        }
    }
}
