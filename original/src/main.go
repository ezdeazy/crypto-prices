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

	"github.com/joho/godotenv"
	"github.com/leekchan/accounting"
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
	SOL  SOL  `json:"SOL"`
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

type SOL struct {
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

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Hello world from our test page!")
}

func doNothing(w http.ResponseWriter, r *http.Request) {}

func getTime() string {
	// returns the current time in PST
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

	return time.Now().In(pacificTime).Format(time.RFC1123)
}

func makeRequest() coins {
	var c coins

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	api_key := goDotEnvVariable("API_KEY")
	q := url.Values{}
	q.Add("symbol", "BTC,ETH,DOGE,SOL")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", api_key)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Err(err).Msg("error reading resp")
		os.Exit(1)
	}

	if err = json.Unmarshal(body, &c); err != nil { // Parse []byte to go struct pointer
		log.Err(err).Msg("Can not unmarshal JSON")
		os.Exit(1)
	}

	return c
}

func getCoinPrice(coin string, c coins) string {
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	var price string

	if coin == "BTC" {
		price = ac.FormatMoneyBigFloat(big.NewFloat(c.Data.BTC.Quote.USD.Price))
	} else if coin == "DOGE" {
		price = ac.FormatMoneyBigFloat(big.NewFloat(c.Data.DOGE.Quote.USD.Price))
	} else if coin == "ETH" {
		price = ac.FormatMoneyBigFloat(big.NewFloat(c.Data.ETH.Quote.USD.Price))
	} else if coin == "SOL" {
		price = ac.FormatMoneyBigFloat(big.NewFloat(c.Data.SOL.Quote.USD.Price))
	} else {
		price = ac.FormatMoneyBigFloat(big.NewFloat(c.Data.BTC.Quote.USD.Price))
	}

	return price
}

func CryptoPrices(w http.ResponseWriter, r *http.Request) {
	// log.Info().Msg("Hello world! Serving crypto prices server")

	dateTime := getTime()
	req := makeRequest()

	btcPrice := getCoinPrice("BTC", req)
	ethPrice := getCoinPrice("ETH", req)
	dogePrice := getCoinPrice("DOGE", req)
	solPrice := getCoinPrice("SOL", req)

	log.Info().Msg("crypto -- BTC: " + btcPrice + "  ETH: " + ethPrice + "  DOGE: " + dogePrice + "  SOL: " + solPrice)

	returnString := "Crypto prices -- " + dateTime + "\n\nBTC:  " + btcPrice + "\nETH:  " + ethPrice + "\nDOGE: " + dogePrice + "\nSOL:  " + solPrice + "\n"

	fmt.Fprintf(w, returnString)
}

func LandingPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Info().Msg("404 -- These are not the droids you are looking for...")
        errorHandler(w, r, http.StatusNotFound)
        return
    }
	log.Warn().Msg(r.URL.Path)
	fmt.Printf("%+v\n", r)
	fmt.Println("Testing")
	log.Info().Msg("landing -- Crypto prices landing page")

	returnString := "Hello World!"

	fmt.Fprintf(w, returnString)
}

func Vault(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("vault -- nothing to see here")

	bop5362 := "https://opensea.io/assets/0x495f947276749ce646f68ac8c248420045cb7b5e/7158387443764788461856806462490708796861921043890951688982477332577228161025"
	aa1591 := "https://opensea.io/assets/0x92b1289ee1c70cb2e51fa63a448ceef2734ec6ff/1591"

	aw147 := "https://solscan.io/token/7RjK5U2Vu5rdJX8qzLxYKT5AEeJ85Rt5YspTsnpXxEjX"
	aw1994 := "https://solscan.io/token/HhPBpvkYZx2PGugo8wjck3LtWP6prRmMJBtpQHck8WQ1"
	aw2745 := "https://solscan.io/token/vCrcerhJB3ct3din5x6ZiygstzzNu5Fn5hii34dbVrT"
	aw2770 := "https://solscan.io/token/9aVQXBzY4PvWDCxSJz2Yho7Q12W2VbuK84347msoRZvV"
	aw2794 := "https://solscan.io/token/BPgEvvSKnW35hxQf398pBBmS1cWJb2hcH3DWSRF5rSBW"
	aw3571 := "https://solscan.io/token/6YkUND1u1rwJwFYP7xaandtzuAHvC1bxWz69sr9s2oiJ"

	ethString := "\nETH NFTs:\n" + bop5362 + "\n" + aa1591 + "\n"
	returnString := ethString + "\nSOL NFTs:\n" + aw147 + "\n" + aw1994 + "\n" + aw2745 + "\n" + aw2770 + "\n" + aw2794 + "\n" + aw3571

	fmt.Fprint(w, returnString)

}

func Health(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("health -- app is healthy! hiya")

	returnString := `( ･∀･)ﾉ ---=== + + +`

	fmt.Fprintf(w, returnString)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
    w.WriteHeader(status)
    if status == http.StatusNotFound {
        fmt.Fprint(w, "These aren't the droids you are looking for...")
    }
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: logTimeFormat}).Level(convertLevel("debug")).With().Caller().Logger()
	log.Info().Msg("Hello world! Crypto prices server is now running")
	http.HandleFunc("/favicon.ico", doNothing)
	http.HandleFunc("/", LandingPage)
	http.HandleFunc("/crypto", CryptoPrices)
	http.HandleFunc("/vault", Vault)
	http.HandleFunc("/health", Health)
	http.ListenAndServe(":8080", nil)
}
