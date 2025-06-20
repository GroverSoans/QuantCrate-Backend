package providers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"os"

	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
)

// InitializePolygon uses the same logic you had — now calls reusable fetcher functions.
func InitializePolygon(db *sql.DB) {
	apiKey := os.Getenv("POLYGON_APIKEY")
	if apiKey == "" {
		log.Fatal("Missing POLYGON_APIKEY env var")
	}
	
	polygonClient := polygon.New(apiKey)
	getPreviousCloseOHLCV(polygonClient, db)

}

func getPreviousCloseOHLCV (client *polygon.Client, db *sql.DB) {
	tickers := []string{"NVDA", "GOOG", "MSFT", "AMZN", "META", "TSLA", "SPY", "X:BTCUSD", "AAPL"}

	// Map polygon tickers to Postgres tickers
	tickerMap := map[string]string {
		"X:BTCUSD" : "BTC-USD",
	}
	
	for _, ticker := range tickers {
		params := models.GetPreviousCloseAggParams{
			Ticker:   ticker,
		}.WithAdjusted(true)

		//fetch data
		res, err := client.GetPreviousCloseAgg(context.Background(), params)
		if err != nil {
			log.Printf("Error fetching %s: %v", ticker, err)
			continue
		}

		if len(res.Results) == 0 {
			log.Printf("No results for %s", ticker)
			continue
		}
		
		data := res.Results[0]

		//Extract UTC date from timestamp (Polygon gives me timestamp)
		var millis int64
		b, _ := json.Marshal(data.Timestamp)
		_ = json.Unmarshal(b, &millis)
		
		t := time.UnixMilli(millis).UTC()
		date := t.Format("2006-01-02")

		//Normalize ticker name for Postgres
		dbTicker, ok := tickerMap[ticker]
		if !ok {
			dbTicker = ticker
		}

		//Validate required fields
		if data.Open == 0 || data.Close == 0 {
			log.Printf("Skipping %s on %s due to missing data", dbTicker, date)
			continue
		}

		// Insert Data into database
		_, err = db.Exec(`
			INSERT INTO ohlcv (ticker, date, open, high, low, close, volume)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (ticker, date) DO NOTHING
		`, dbTicker, date, data.Open, data.High, data.Low, data.Close, data.Volume)

		if err != nil {
			log.Printf("DB insert error for %s on %s: %v", dbTicker, date, err)
			continue
		}

		log.Printf("Inserted %s for %s", date, dbTicker)

		log.Printf("Previous close for %s: %+v\n", ticker, res)
		// Wait to avoid rate limiting (4 request per min)
		time.Sleep(15 * time.Second)
	}
}