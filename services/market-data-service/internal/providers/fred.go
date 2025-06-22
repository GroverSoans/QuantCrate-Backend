package providers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// FREDResponse represents the API response structure
type FREDResponse struct {
	Observations []FREDObservation `json:"observations"`
}

// FREDObservation represents a single data point
type FREDObservation struct {
	Date  string `json:"date"`
	Value string `json:"value"`
}

// InterestRateData represents processed interest rate data
type InterestRateData struct {
	Date     time.Time
	Value    float64
	SeriesID string
}

// FetchFREDSeries fetches interest rate data from FRED for the given series ID
func FetchFREDSeries(seriesID string, startDate string, endDate string) ([]InterestRateData, error) {
	apiKey := os.Getenv("FRED_APIKEY")
	if apiKey == "" {
		return nil, fmt.Errorf("FRED_APIKEY environment variable not set")
	}

	url := fmt.Sprintf("https://api.stlouisfed.org/fred/series/observations?series_id=%s&api_key=%s&file_type=json&frequency=d&observation_start=%s&observation_end=%s",
		seriesID, apiKey, startDate, endDate)

	fmt.Printf("Fetching FRED data for series %s from %s to %s\n", seriesID, startDate, endDate)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from FRED: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var fredResp FREDResponse
	if err := json.Unmarshal(body, &fredResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	var data []InterestRateData
	for _, obs := range fredResp.Observations {
		// Skip missing or invalid values
		if obs.Value == "." || obs.Value == "" {
			continue
		}

		val, err := strconv.ParseFloat(obs.Value, 64)
		if err != nil {
			log.Printf("Warning: Could not parse value '%s' for date %s", obs.Value, obs.Date)
			continue
		}

		parsedDate, err := time.Parse("2006-01-02", obs.Date)
		if err != nil {
			log.Printf("Warning: Could not parse date '%s'", obs.Date)
			continue
		}

		data = append(data, InterestRateData{
			Date:     parsedDate,
			Value:    val,
			SeriesID: seriesID,
		})
	}

	fmt.Printf("Successfully fetched %d data points for series %s\n", len(data), seriesID)
	return data, nil
}

// GetLast15YearsInterestRates fetches the Daily Federal Funds Rate for the last 15 years
func GetLast15YearsInterestRates(db *sql.DB) {
	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(-15, 0, 0).Format("2006-01-02")

	fmt.Printf("Fetching Daily Federal Funds Rate (DFF) from %s to %s\n", startDate, endDate)

	// Fetch only DFF (Daily Federal Funds Rate)
	data, err := FetchFREDSeries("DFF", startDate, endDate)
	if err != nil {
		log.Printf("Error fetching DFF: %v", err)
		return
	}

	fmt.Printf("Successfully fetched %d Federal Funds Rate data points\n", len(data))

	// Print sample data
	fmt.Println("\nSample Federal Funds Rate data:")
	for i, point := range data {
		if i < 10 || i >= len(data)-5 {
			fmt.Printf("  %s: %.4f%%\n", point.Date.Format("2006-01-02"), point.Value)
		} else if i == 10 {
			fmt.Printf("  ... (%d more entries) ...\n", len(data)-15)
		}
	}

	// Save to database
	if db != nil {
		saveFederalFundsRate(db, data)
	}
}

// saveFederalFundsRate saves Federal Funds Rate data to database
func saveFederalFundsRate(db *sql.DB, data []InterestRateData) {
	if len(data) == 0 {
		fmt.Println("No Federal Funds Rate data to save")
		return
	}

	fmt.Printf("Saving %d Federal Funds Rate records to database...\n", len(data))

	// Start transaction for efficiency
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return
	}
	defer tx.Rollback()

	// Prepare statement
	stmt, err := tx.Prepare(`
		INSERT INTO interest_rates (series_id, series_name, date, rate)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (series_id, date) DO UPDATE SET
			rate = EXCLUDED.rate
	`)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return
	}
	defer stmt.Close()

	successCount := 0
	for _, point := range data {
		_, err := stmt.Exec(
			"DFF",
			"Daily Federal Funds Rate",
			point.Date.Format("2006-01-02"),
			point.Value,
		)
		if err != nil {
			log.Printf("Error inserting DFF on %s: %v", point.Date.Format("2006-01-02"), err)
			continue
		}
		successCount++
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return
	}

	fmt.Printf("Successfully saved %d Federal Funds Rate records\n", successCount)
}
