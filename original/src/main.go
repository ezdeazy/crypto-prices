package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/leekchan/accounting"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)


type coins struct {
	Data data `json:"data"`
}

type data struct {
	BTC  BTC  `json:"BTC"`
	ETH  ETH  `json:"ETH"`
	DOGE DOGE `json:"DOGE"`
}

type BTC struct {
	Id                int         `json:"id"`
	Name              string      `json:"name"`
	Symbol            string      `json:"symbol"`
	Slug              string      `json:"slug"`
	NumMarketPairs    int         `json:"num_market_pairs"`
	DateAdded         time.Time   `json:"date_added"`
	MaxSupply         int         `json:"max_supply"`
	CirculatingSupply int         `json:"circulating_supply"`
	TotalSupply       int         `json:"total_supply"`
	IsActive          int         `json:"is_active"`
	Platform          interface{} `json:"platform"`
	CmcRank           int         `json:"cmc_rank"`
	IsFiat            int         `json:"is_fiat"`
	LastUpdated       time.Time   `json:"last_updated"`
	Quote             Quote       `json:"quote"`
}

type ETH struct {
	Id                int         `json:"id"`
	Name              string      `json:"name"`
	Symbol            string      `json:"symbol"`
	Slug              string      `json:"slug"`
	NumMarketPairs    int         `json:"num_market_pairs"`
	DateAdded         time.Time   `json:"date_added"`
	MaxSupply         int         `json:"max_supply"`
	CirculatingSupply float64     `json:"circulating_supply"`
	TotalSupply       float64     `json:"total_supply"`
	IsActive          int         `json:"is_active"`
	Platform          interface{} `json:"platform"`
	CmcRank           int         `json:"cmc_rank"`
	IsFiat            int         `json:"is_fiat"`
	LastUpdated       time.Time   `json:"last_updated"`
	Quote             Quote       `json:"quote"`
}

type DOGE struct {
	Id                int         `json:"id"`
	Name              string      `json:"name"`
	Symbol            string      `json:"symbol"`
	Slug              string      `json:"slug"`
	NumMarketPairs    int         `json:"num_market_pairs"`
	DateAdded         time.Time   `json:"date_added"`
	MaxSupply         int         `json:"max_supply"`
	CirculatingSupply float64     `json:"circulating_supply"`
	TotalSupply       float64     `json:"total_supply"`
	IsActive          int         `json:"is_active"`
	Platform          interface{} `json:"platform"`
	CmcRank           int         `json:"cmc_rank"`
	IsFiat            int         `json:"is_fiat"`
	LastUpdated       time.Time   `json:"last_updated"`
	Quote             Quote       `json:"quote"`
}

type Quote struct {
	USD USD `json:"USD"`
}

type USD struct {
	Price                 float64   `json:"price"`
	Volume24H             float64   `json:"volume_24h"`
	VolumeChange24H       float64   `json:"volume_change_24h"`
	PercentChange1H       float64   `json:"percent_change_1h"`
	PercentChange24H      float64   `json:"percent_change_24h"`
	PercentChange7D       float64   `json:"percent_change_7d"`
	PercentChange30D      float64   `json:"percent_change_30d"`
	PercentChange60D      float64   `json:"percent_change_60d"`
	PercentChange90D      float64   `json:"percent_change_90d"`
	MarketCap             float64   `json:"market_cap"`
	MarketCapDominance    float64   `json:"market_cap_dominance"`
	FullyDilutedMarketCap float64   `json:"fully_diluted_market_cap"`
	LastUpdated           time.Time `json:"last_updated"`
}

const logTimeFormat = "2006-01-02T15:04:05.000Z07:00"


func goDotEnvVariable(key string) string {

	err := godotenv.Load(".env")
  
	if err != nil {
		// do nothing to load from environment variable
	}
  
	return os.Getenv(key)
}

func convertLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "trace":
		 return zerolog.TraceLevel
	case "debug":
		 return zerolog.DebugLevel
	case "info":
		 return zerolog.InfoLevel
	case "warn":
		 return zerolog.WarnLevel
	case "error":
		 return zerolog.ErrorLevel
	case "fatal":
		 return zerolog.FatalLevel
	case "panic":
		 return zerolog.PanicLevel
	default:
		 return zerolog.InfoLevel
	}
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: logTimeFormat}).Level(convertLevel("debug")).With().Caller().Logger()
	log.Info().Msg("Listening on port :8080")
	http.HandleFunc("/favicon.ico", doNothing)
  http.HandleFunc("/", CryptoPrices)
	http.HandleFunc("/testing", HelloWorld)
	http.ListenAndServe(":8080", nil)
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Hello world from our test page!")
}

func CryptoPrices(w http.ResponseWriter, r *http.Request) {

	log.Info().Msg("Hello world! Serving crypto prices server")

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
		log.Print(err)
	  os.Exit(1)
	}
	defer resp.Body.Close()

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


	var c coins

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		 log.Err(err).Msg("error reading resp")
		 os.Exit(1)
	}

	if err = json.Unmarshal(body, &c); err != nil { // Parse []byte to go struct pointer
		 log.Err(err).Msg("Can not unmarshal JSON")
		 os.Exit(1)
	}

	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	btc_price := ac.FormatMoneyBigFloat(big.NewFloat(c.Data.BTC.Quote.USD.Price))
	eth_price := ac.FormatMoneyBigFloat(big.NewFloat(c.Data.ETH.Quote.USD.Price))
	doge_price := ac.FormatMoneyBigFloat(big.NewFloat(c.Data.DOGE.Quote.USD.Price))
	
	log.Info().Msg("Returning BTC Price: " + btc_price)
	log.Info().Msg("Returning ETH Price: " + eth_price)
	log.Info().Msg("Returning DOGE Price: " + doge_price)

	returnString := "Crypto prices -- " + dateTime + "\n\nBTC:  " + btc_price + "\nETH:  " + eth_price + "\nDOGE: " + doge_price + "\n"

	fmt.Fprintf(w, returnString)
}


func doNothing(w http.ResponseWriter, r *http.Request){}
