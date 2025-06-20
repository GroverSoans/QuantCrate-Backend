import yfinance as yf
import pandas as pd
import psycopg2 
from psycopg2.extras import execute_values
import datetime

#Postgres Config
dbconn = {
    "dbname" : "marketData",
    "user" : "postgres",
    "password" : "C4killsu",
    "host" : "localhost",
    "port" : "5432"
}

tickers = ["NVDA", "GOOG", "MSFT", "AMZN", "META", "TSLA", "SPY", "BTC-USD", "AAPL"]

start_date = "2010-01-01"
end_date = datetime.date.today().strftime("%Y-%m-%d")

def insert_ohlcv(cursor, ticker, df):
    df = df.dropna()

    records = []
    for idx, row in df.iterrows():
        try:
            records.append((
                ticker,
                idx.date(),
                float(row["Open"].item() if hasattr(row["Open"], "item") else row["Open"]),
                float(row["High"].item() if hasattr(row["High"], "item") else row["High"]),
                float(row["Low"].item() if hasattr(row["Low"], "item") else row["Low"]),
                float(row["Close"].item() if hasattr(row["Close"], "item") else row["Close"]),
                int(row["Volume"].item() if hasattr(row["Volume"], "item") else row["Volume"]),
            ))
        except Exception as e:
            print(f"⚠️ Skipping row {idx} due to error: {e}")

    if records:
        sql = """
        INSERT INTO ohlcv (ticker, date, open, high, low, close, volume)
        VALUES %s
        ON CONFLICT (ticker, date) DO NOTHING;
        """
        execute_values(cursor, sql, records)
        print(f"✅ Inserted {len(records)} rows for {ticker}")
    else:
        print(f"⚠️ No valid data to insert for {ticker}")

def main():
    conn = psycopg2.connect(**dbconn)
    cursor = conn.cursor()

    for ticker in tickers:
        print(f"📥 Fetching data for {ticker}...")
        df = yf.download(ticker, start=start_date, end=end_date, interval="1d", auto_adjust=True)
        if not df.empty:
            insert_ohlcv(cursor, ticker, df)
        else:
            print(f"⚠️ No data returned for {ticker}")

    conn.commit()
    cursor.close()
    conn.close()
    print("🎉 All done.")

if __name__ == "__main__":
    main()