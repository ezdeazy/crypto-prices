package main

import (
	"encoding/json"
	"fmt"
  	"io/ioutil"
  	"log"
  	"net/http"
  	"net/url"
  	"os"
  	"time"

  	"github.com/leekchan/accounting"
  	"github.com/joho/godotenv"
)


func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")
  
	if err != nil {
		// do nothing to load from environment variable
		// log.Fatalf("Error loading .env file")
	}
  
	return os.Getenv(key)
}


func main() {
  client := &http.Client{}
  req, err := http.NewRequest("GET","https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest", nil)
  if err != nil {
    log.Print(err)
    os.Exit(1)
  }

  api_key := goDotEnvVariable("API_KEY")
  q := url.Values{}
  //   q.Add("start", "1")
  //   q.Add("limit", "5000")
  //   q.Add("convert", "USD")
  q.Add("symbol", "BTC,ETH,DOGE")


  req.Header.Set("Accepts", "application/json")
  req.Header.Add("X-CMC_PRO_API_KEY", api_key)
  req.URL.RawQuery = q.Encode()


  resp, err := client.Do(req);
  if err != nil {
    fmt.Println("Error sending request to server")
    os.Exit(1)
  }
  // fmt.Println(resp.Status);
  respBody, _ := ioutil.ReadAll(resp.Body)

  // print the full response
  // fmt.Println(string(respBody));

  var result map[string]interface{}
  json.Unmarshal([]byte(respBody), &result)

  // uncomment the following to see pretty output of response
  // b, err := json.MarshalIndent(result, "", "  ")
  // if err != nil {
  //     fmt.Println("error:", err)
  // }

  // fmt.Print(string(b))
  // fmt.Print("\n\n")

  data := result["data"].(map[string]interface{})

  resp_status := result["status"].(map[string]interface{})
  resp_time := resp_status["timestamp"].(string)
  

  layout := "2006-01-02T15:04:05.000Z"
  t, err := time.Parse(layout, resp_time)

  if err != nil {
      fmt.Println(err)
  }


  layout_date := "January 02, 2006"
  layout_time := "03:04:05 PM"

  theDate := t.Format(layout_date)
  theTime := t.Format(layout_time)
  dateTime := theDate + " " + theTime


  btc := data["BTC"].(map[string]interface{})
  btc_name := btc["name"].(string)
  btc_symbol := btc["symbol"].(string)
  btc_quote := btc["quote"].(map[string]interface{})
  btc_usd := btc_quote["USD"].(map[string]interface{})
  btc_price := btc_usd["price"].(float64)
  // fmt.Printf("%T\n", btc_price)


  ac := accounting.Accounting{Symbol: "$", Precision: 2}

  btc_price_formatted := ac.FormatMoney(btc_price)
  btc_returnString := btc_name + " (" + btc_symbol + "): " + btc_price_formatted


  eth := data["ETH"].(map[string]interface{})
  eth_name := eth["name"].(string)
  eth_symbol := eth["symbol"].(string)
  eth_quote := eth["quote"].(map[string]interface{})
  eth_usd := eth_quote["USD"].(map[string]interface{})
  eth_price := eth_usd["price"].(float64)

  eth_price_formatted := ac.FormatMoney(eth_price)
  eth_returnString := eth_name + " (" + eth_symbol + "): " + eth_price_formatted


  doge := data["DOGE"].(map[string]interface{})
  doge_name := doge["name"].(string)
  doge_symbol := doge["symbol"].(string)
  doge_quote := doge["quote"].(map[string]interface{})
  doge_usd := doge_quote["USD"].(map[string]interface{})
  doge_price := doge_usd["price"].(float64)

  doge_price_formatted := ac.FormatMoney(doge_price)
  doge_returnString := doge_name + " (" + doge_symbol + "): " + doge_price_formatted


  fmt.Print("\n")
  fmt.Println("Time: ", dateTime)
  fmt.Print("\n")
  fmt.Println(btc_returnString)
  fmt.Println(eth_returnString)
  fmt.Println(doge_returnString)
  fmt.Print("\n\n")
}

