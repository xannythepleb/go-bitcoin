package main

import (
       "encoding/json"
        "fmt"
        "net/http"
        "time"
)
type BitcoinPrice struct {
        USD int `json:"usd"`
}
type PriceData struct {
        Bitcoin BitcoinPrice `json:"bitcoin"`
}
type FeeData struct {
        FastestFee  int `json:"fastestFee"`
        HalfHourFee int `json:"halfHourFee"`
        HourFee     int `json:"hourFee"`
        EconomyFee  int `json:"economyFee"`
        MinimumFee  int `json:"minimumFee"`
}

// Simple script that shows current Bitcoin price and recommended transaction fees in the terminal
// Price data from Coingecko's API, transaction fee data from Mempool's API

// TODO:
// - Choice of APIs including support for own node
// - Bootstrap Tor to connect to own node through onion URL while outside LAN
// - User selected timezone or detect system time
// - Create some kind of UX, Android widget would be useful
// - Telegram bot would be cool
// - More data? Perhaps add flags to choose how much or little data you wanna see - Mempool's API has much more to offer!
// - Keep stacking sats
// - Overthrow government

func main() {
        // Grab current Bitcoin price from Coingecko API
        resp, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd")
        if err != nil {
                fmt.Println("Error downloading price data from Coingecko API:", err)
                return
        }
        defer resp.Body.Close()
        var priceData PriceData
        err = json.NewDecoder(resp.Body).Decode(&priceData)
        if err != nil {
                fmt.Println("Error decoding price JSON from Coingecko API:", err)
                return
        }
        price := priceData.Bitcoin.USD

        // Grab current tx fees from mempool.space API
        resp, err = http.Get("https://mempool.space/api/v1/fees/recommended")
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
        fmt.Printf("  1 BTC = $%d\n", price)
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
