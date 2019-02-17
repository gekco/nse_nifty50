package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "time"
    "html/template"
)


var gainers map[string]interface{}
var losers map[string]interface{}

func getDataFromUrl(url string, result *map[string]interface{}) {
    response, err := http.Get(url)
    
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    // Create a goquery document from the HTTP response
    body, err := ioutil.ReadAll(response.Body)

    // Find and print image URLs
    bs := string(body)
 
    json.Unmarshal([]byte(bs), result)
}

func refreshGainersLosers() {
   for true {
        getDataFromUrl("https://www.nseindia.com/live_market/dynaContent/live_analysis/losers/niftyLosers1.json", &losers)
        getDataFromUrl("https://www.nseindia.com/live_market/dynaContent/live_analysis/gainers/niftyGainers1.json", &gainers)
        time.Sleep(300*time.Second)
    }
}

func mainView(w http.ResponseWriter, r *http.Request) {
    t,err := template.ParseFiles("main.html")
    if err != nil {
        fmt.Fprint(w, "ERROR !!!!! COME BACK LATER :)")
        fmt.Println(err)
        return
    }
    temp_args := map[string]interface{}{ "NIFTY 50 GAINERS" : gainers["data"], "NIFTY 50 LOSERS": losers["data"]}
    t.Execute(w, temp_args)
}

func main() {
    go refreshGainersLosers()
    http.HandleFunc("/", mainView)
    log.Fatal(http.ListenAndServe(":5000", nil))
}