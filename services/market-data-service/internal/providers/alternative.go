package providers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Structs for parsing the JSON response
type FearGreedResponse struct {
	Name string          `json:"name"`
	Data []FearGreedData `json:"data"`
}

type FearGreedData struct {
	Value               string `json:"value"`
	ValueClassification string `json:"value_classification"`
	Timestamp           string `json:"timestamp"`
	TimeUntilUpdate     string `json:"time_until_update,omitempty"`
}

func GetFearAndGreedIndex(db *sql.DB) {
	res, err := http.Get("https://api.alternative.me/fng/?limit=1") //Gets most resent index
	if err != nil {
		log.Fatalf("Error fetching fear and greed index from alternative: %v", err)
	}
	defer res.Body.Close() 

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Print the raw JSON response
	fmt.Printf("Fear and Greed Index Response: %s\n", string(body))

	// Parse the JSON response
	var fearGreedResponse FearGreedResponse
	if err := json.Unmarshal(body, &fearGreedResponse); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	// Save to database
	SaveToPostgreSQL(db, fearGreedResponse.Data)
}

func SaveToPostgreSQL(db *sql.DB, data []FearGreedData) {
	// Prepare the SQL statement
	stmt, err := db.Prepare(`
		INSERT INTO fgi_index (value, value_classification, timestamp_unix, resolved_date)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (resolved_date) DO UPDATE SET
			value = EXCLUDED.value,
			value_classification = EXCLUDED.value_classification,
			timestamp_unix = EXCLUDED.timestamp_unix
	`)
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return
	}
	defer stmt.Close()

	// Insert each data point
	for _, item := range data {
		// Convert string value to integer
		value, err := strconv.Atoi(item.Value)
		if err != nil {
			log.Printf("Error converting value to int: %v", err)
			continue
		}

		// Convert timestamp to integer
		timestampUnix, err := strconv.ParseInt(item.Timestamp, 10, 64)
		if err != nil {
			log.Printf("Error converting timestamp to int: %v", err)
			continue
		}

		// Convert timestamp to date
		resolvedDate := time.Unix(timestampUnix, 0).UTC().Format("2006-01-02")

		// Execute the insert
		_, err = stmt.Exec(value, item.ValueClassification, timestampUnix, resolvedDate)
		if err != nil {
			log.Printf("Error inserting data into database: %v", err)
			continue
		}

		fmt.Printf("Successfully saved Fear & Greed Index: %d (%s) for date %s\n",
			value, item.ValueClassification, resolvedDate)
	}
}
