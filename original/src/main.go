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

	err := godotenv.Load(".env")
  
	if err != nil {
		// do nothing to load from environment variable
	}
  
	return os.Getenv(key)
}


func main() {
	log.Printf("Listening on port :8080")
	http.HandleFunc("/favicon.ico", doNothing)
  http.HandleFunc("/", CryptoPrices)
	http.HandleFunc("/testing", HelloWorld)
	http.ListenAndServe(":8080", nil)
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	log.Printf("Hello world from our test page!")
}

func CryptoPrices(w http.ResponseWriter, r *http.Request) {
	log.Printf("Hello world! Serving crypto prices server")
	// fmt.Fprintf(w, "directory, %s!", r.URL.Path[1:])

	client := &http.Client{}
	req, err := http.NewRequest("GET","https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest", nil)
	if err != nil {
	  log.Print(err)
	  os.Exit(1)
	}
  
	api_key := goDotEnvVariable("API_KEY")
	q := url.Values{}
	q.Add("symbol", "BTC,ETH,DOGE")
  
  
	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", api_key)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req);
	if err != nil {
	  fmt.Println("Error sending request to server")
	  os.Exit(1)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal([]byte(respBody), &result)


	data := result["data"].(map[string]interface{})


	// manually set time zone
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Printf("error loading location '%s': %v\n", tz, err)
		}
	}
  
	// convert UTC to pacific
	pacificTime, errTime := time.LoadLocation("America/Los_Angeles")
	if errTime != nil {
		fmt.Println("err: ", errTime.Error())
	}
	dateTime := time.Now().In(pacificTime).Format(time.RFC1123)


	btc := data["BTC"].(map[string]interface{})
	btc_name := btc["name"].(string)
	btc_symbol := btc["symbol"].(string)
	btc_quote := btc["quote"].(map[string]interface{})
	btc_usd := btc_quote["USD"].(map[string]interface{})
	btc_price := btc_usd["price"].(float64)


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


	returnString := "Crypto prices -- " + dateTime + "\n\n" + btc_returnString + "\n" + eth_returnString + "\n" + doge_returnString
	fmt.Fprintf(w, returnString)
}


func doNothing(w http.ResponseWriter, r *http.Request){}
