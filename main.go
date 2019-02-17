package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "time"
    "html/template"
    "os"
)


var gainers map[string]interface{}
var losers map[string]interface{}

func getDataFromUrl(url string, result *map[string]interface{}) {
    fmt.Println("__+++______HERE 1")
    response, err := http.Get(url)
    fmt.Println("__+++______HERE 2")

    if err != nil {
        log.Fatal(err)
        fmt.Println("__+++______",err)
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
        fmt.Println("__+++______HERE 0")
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
    fmt.Println(temp_args)
    t.Execute(w, temp_args)
}

func main() {
    go refreshGainersLosers()
    http.HandleFunc("/", mainView)
    var port = os.Getenv("PORT")
    if port == "" {
        port = "5000"
    }
    port =":"+port
    fmt.Println("______",port)
    log.Fatal(http.ListenAndServe(port, nil))
}