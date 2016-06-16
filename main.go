package main

import (
    "net/http"
    "io/ioutil"
    "encoding/json"
    "strings"
    "io"
    "log"
    "time"
    "os"
)

// LastTimestamp //
var LastTimestamp int

func main() {
    if len(os.Args) != 3 {
        log.Println("Usage: VKFakeOnline <login> <password>")
        return
    }

    token := VKAuth(os.Args[1], os.Args[2]).AccessToken
    if token != "" {
        log.Println("VKFakeOnline has been launched")
    } else {
        log.Println("Error: An error occurred when authentication")
        return 
    }
    OfflineLoop(token)
}

// Response //
type Response struct {
    AccessToken string  `json:"access_token"`
    ExpiresIn   int     `json:"expires_in"`
    UserID      int     `json:"user_id"`
}

// OfflineLoop //
func OfflineLoop(token string) {
    if LastTimestamp == 0 {
        LastTimestamp = VKSendOffline(token)
        time.Sleep(time.Second * 7)
        OfflineLoop(token)
    }

    if LastTimestamp != VKSendOffline(token) {
        log.Println("Warning: The latest online time has been changed.")
    }

    time.Sleep(time.Second * 7)
    OfflineLoop(token)
}

// VKSendOffline //
func VKSendOffline(token string) int {
    resp, err := http.Get("https://api.vk.com/method/execute?code=API.account.setOffline()%3Bvar+u%3DAPI.users.get(%7Bfields%3A%22online%2Clast_seen%22%7D)%5B0%5D%3Breturn{time:u.last_seen.time,}%3B&access_token="+token)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    type JSONBody struct {
        Response struct {
            Timestamp int `json:"time"`
        } `json:"response"`
    }

    var JSONresp JSONBody

    if err = json.Unmarshal(body, &JSONresp); err != nil {
        return -1
    }

    return JSONresp.Response.Timestamp
}

// VKAuth //
func VKAuth(login string, password string) Response {
    resp, err := http.Get("https://oauth.vk.com/token?grant_type=password&client_id=2274003&client_secret=hHbZxrka2uZ6jB1inYsH&username=" + login + "&password="+password+"&v=5.45")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    dec := json.NewDecoder(strings.NewReader(string(body)))
    var m Response
    for {
        if err := dec.Decode(&m); err == io.EOF {
            break
        } else if err != nil {
            log.Fatal(err)
        }
        return m
    }
    return m
}