package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"net/http"
	"strconv"
	"strings"
)


// Simple script that shows current Bitcoin price and recommended transaction fees in the terminal
// Price data from Coingecko's API, transaction fee data from Mempool's API

// TODO:
// - Choice of APIs NOW IMPLEMENTED!
// - Add support for own node
// - Bootstrap Tor to connect to own node through onion URL while outside LAN
// - User selected timezone or detect system time
// - Create some kind of UX, Android widget would be useful
// - Telegram bot would be cool
// - More data? Perhaps add flags to choose how much or little data you wanna see - Mempool's API has much more to offer!
// - Keep stacking sats
// - Overthrow government


type FeeData struct {
        FastestFee  int `json:"fastestFee"`
        HalfHourFee int `json:"halfHourFee"`
        HourFee     int `json:"hourFee"`
        EconomyFee  int `json:"economyFee"`
        MinimumFee  int `json:"minimumFee"`
}

type CoinDeskResponse struct {
	Bpi struct {
		USD struct {
			Rate string `json:"rate"`
		} `json:"USD"`
	} `json:"bpi"`
}

type CoinGeckoResponse struct {
	Bitcoin struct {
		Usd float64 `json:"usd"`
	} `json:"bitcoin"`
}

type BitfinexResponse struct {
	LastPrice string `json:"last_price"`
}

type KrakenResponse struct {
	Result struct {
		XXBTZUSD struct {
			C []string `json:"c"`
		} `json:"XXBTZUSD"`
	} `json:"result"`
}

type BinanceResponse struct {
	Price string `json:"price"`
}

type MempoolResponse map[string]int

func getBitcoinPriceFromCoinDesk() (float64, error) {
	resp, err := http.Get("https://api.coindesk.com/v1/bpi/currentprice/BTC.json")
	if err != nil {
		return 0, fmt.Errorf("Could not fetch from CoinDesk: %w", err)
	}
	defer resp.Body.Close()

	var data CoinDeskResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("Could not parse CoinDesk response: %w", err)
	}

	price, err := strconv.ParseFloat(strings.Replace(data.Bpi.USD.Rate, ",", "", -1), 64)
	if err != nil {
		return 0, fmt.Errorf("Could not parse CoinDesk price: %w", err)
	}

	return price, nil
}

func getBitcoinPriceFromCoinGecko() (float64, error) {
	resp, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd")
	if err != nil {
		return 0, fmt.Errorf("Could not fetch from CoinGecko: %w", err)
	}
	defer resp.Body.Close()

	var data CoinDeskResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("Could not parse CoinGecko response: %w", err)
	}

	price, err := strconv.ParseFloat(strings.Replace(data.Bpi.USD.Rate, ",", "", -1), 64)
	if err != nil {
		return 0, fmt.Errorf("Could not parse CoinGecko price: %w", err)
	}

	return price, nil
}

func getBitcoinPriceFromBitfinex() (float64, error) {
	resp, err := http.Get("https://api.bitfinex.com/v1/pubticker/btcusd")
	if err != nil {
		return 0, fmt.Errorf("Could not fetch from Bitfinex: %w", err)
	}
	defer resp.Body.Close()

	var data BitfinexResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal(err)
	}

	price, err := strconv.ParseFloat(data.LastPrice, 64)
	if err != nil {
		log.Fatal(err)
	}

	return price, nil
}

func getBitcoinPriceFromKraken() (float64, error) {
	resp, err := http.Get("https://api.kraken.com/0/public/Ticker?pair=BTCUSD")
	if err != nil {
		return 0, fmt.Errorf("Could not fetch from Kraken: %w", err)
	}
	defer resp.Body.Close()

	var data KrakenResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal(err)
	}

	price, err := strconv.ParseFloat(data.Result.XXBTZUSD.C[0], 64)
	if err != nil {
		log.Fatal(err)
	}

	return price, nil
}

func getBitcoinPriceFromBinance() (float64, error) {
	resp, err := http.Get("https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT")
	if err != nil {
		return 0, fmt.Errorf("Could not fetch from Binance: %w", err)
	}
	defer resp.Body.Close()

	var data BinanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal(err)
	}

	price, err := strconv.ParseFloat(data.Price, 64)
	if err != nil {
		log.Fatal(err)
	}
	return price, nil
}

func main() {
	var bitcoinPrice float64
	var err error

	fmt.Println("Choose an API to grab the Bitcoin price from:")
	fmt.Println("1. CoinDesk")
	fmt.Println("2. CoinGecko")
	fmt.Println("3. Bitfinex")
	fmt.Println("4. Kraken")
	fmt.Println("5. Binance")

	var apiChoice int
	_, err = fmt.Scan(&apiChoice)
	if err != nil {
		log.Fatal(err)
	}

	switch apiChoice {
	case 1:
		bitcoinPrice, err = getBitcoinPriceFromCoinDesk()
	case 2:
		bitcoinPrice, err = getBitcoinPriceFromCoinGecko()
	case 3:
		bitcoinPrice, err = getBitcoinPriceFromBitfinex()
	case 4:
		bitcoinPrice, err = getBitcoinPriceFromKraken()
	case 5:
		bitcoinPrice, err = getBitcoinPriceFromBinance()
	default:
		log.Fatal("Invalid choice")
	}

	if err != nil {
		log.Fatalf("Failed to get Bitcoin price: %v", err)
	}


	// Format the bitcoin price with commas
	formattedBitcoinPrice := strconv.FormatFloat(bitcoinPrice, 'f', 2, 64)
	parts := strings.Split(formattedBitcoinPrice, ".")
	integerPart := parts[0]
	decimalPart := parts[1]

	// Split the integer part into groups of three
	var groups []string
	for len(integerPart) > 0 {
		if len(integerPart) < 3 {
			groups = append([]string{integerPart}, groups...)
			break
		} else {
			groups = append([]string{integerPart[len(integerPart)-3:]}, groups...)
			integerPart = integerPart[:len(integerPart)-3]
		}
	}

	// Join the groups with commas
	formattedBitcoinPrice = strings.Join(groups, ",") + "." + decimalPart

        // Grab current tx fees from mempool.space API
        resp, err := http.Get("https://mempool.space/api/v1/fees/recommended")
        if err != nil {
                fmt.Println("Error downloading tx fee data from Mempool API:", err)
                return
        }
        defer resp.Body.Close()
        var feeData FeeData
        err = json.NewDecoder(resp.Body).Decode(&feeData)
        if err != nil {
                fmt.Println("Error decoding tx fee JSON from Mempool API:", err)
                return
        }
        
        // Setting it to UK timezone cuz that's where I live innit fam
        loc, err := time.LoadLocation("Europe/London")
        if err != nil {
                fmt.Println("Error loading timezone location:", err)
                return
        }
        // Grab the time in above timezone
        t := time.Now().In(loc)
        // Date format
        dateStr := fmt.Sprintf("%d%s of %s, %d", t.Day(), getDaySuffix(t.Day()), t.Month().String(), t.Year())
        // Time format
        timeStr := t.Format("03:04pm")
        
	// Finally we can print the data!
	fmt.Printf("\nOn the %s at %s:\n", dateStr, timeStr)
        fmt.Printf("  1 BTC = $%s\n", formattedBitcoinPrice)
        fmt.Printf("  1 BTC = 1 BTC\n")
        fmt.Printf("  The recommended tx fees are:\n")
        fmt.Printf("    - Fast: %d sat/byte\n", feeData.FastestFee)
        fmt.Printf("    - Half hour: %d sat/byte\n", feeData.HalfHourFee)
        fmt.Printf("    - Hour: %d sat/byte\n", feeData.HourFee)
        fmt.Printf("    - Economy: %d sat/byte\n", feeData.EconomyFee)
        fmt.Printf("    - Minimum: %d sat/byte\n\n", feeData.MinimumFee)
        fmt.Printf("   #FreeRoss\n\n")
}

// One more function for date suffix (e.g. "st", "nd", "rd", "th")
func getDaySuffix(day int) string {
        if day >= 11 && day <= 13 {
                return "th"
        }
        switch day % 10 {
        case 1:
                return "st"
        case 2:
                return "nd"
        case 3:
                return "rd"
        default:
                return "th"
        }

}
