#!/usr/bin/env python3
"""
Export script to create 3 separate CSV files from PostgreSQL database for ML model training.

This script exports:
1. Bitcoin OHLCV data → data/bitcoin_ohlcv.csv
2. Federal Funds Rate → data/federal_funds_rate.csv  
3. Fear & Greed Index → data/fear_greed_index.csv

Each dataset is saved as a separate CSV file in the data/ directory.
"""

import os
from dotenv import load_dotenv
import pandas as pd
import psycopg2
from datetime import datetime
import sys

load_dotenv()
# Database connection parameters
DB_CONFIG = os.getenv("DATABASE_URL")

def connect_to_database():
    """Connect to PostgreSQL database"""
    try:
        if DB_CONFIG:
            # Using DATABASE_URL connection string
            conn = psycopg2.connect(DB_CONFIG)
            print(f"✅ Connected to database using DATABASE_URL")
        else:
            print(f"❌ DATABASE_URL environment variable not set")
            print("Please set DATABASE_URL environment variable:")
            print("  export DATABASE_URL='postgresql://username:password@host:port/database'")
            sys.exit(1)
        return conn
    except Exception as e:
        print(f"❌ Error connecting to database: {e}")
        print("Make sure your database is running and DATABASE_URL is correct:")
        print("  export DATABASE_URL='postgresql://username:password@host:port/database'")
        sys.exit(1)

def export_bitcoin_data(conn):
    """Export Bitcoin OHLCV data to CSV"""
    print("\n📊 Exporting Bitcoin OHLCV data...")
    
    query = """
    SELECT ticker, date, open, high, low, close, volume
    FROM ohlcv 
    WHERE ticker = 'BTC' OR ticker = 'BTCUSD' OR ticker LIKE '%BTC%'
    ORDER BY date ASC
    """
    
    try:
        df = pd.read_sql(query, conn)
        
        if df.empty:
            print("⚠️  No Bitcoin data found in ohlcv table")
            # Create empty CSV with headers
            df = pd.DataFrame(columns=['ticker', 'date', 'open', 'high', 'low', 'close', 'volume'])
        
        output_path = 'data/bitcoin_ohlcv.csv'
        df.to_csv(output_path, index=False)
        print(f"✅ Bitcoin data exported: {output_path} ({len(df)} rows)")
        
        if not df.empty:
            print(f"   Date range: {df['date'].min()} to {df['date'].max()}")
            print(f"   Tickers: {df['ticker'].unique().tolist()}")
        
        return df
    except Exception as e:
        print(f"❌ Error exporting Bitcoin data: {e}")
        return pd.DataFrame()

def export_interest_rates(conn):
    """Export Federal Funds Rate data to CSV"""
    print("\n📈 Exporting Interest Rates data...")
    
    query = """
    SELECT series_id, series_name, date, rate
    FROM interest_rates 
    WHERE series_id = 'DFF'
    ORDER BY date ASC
    """
    
    try:
        df = pd.read_sql(query, conn)
        
        if df.empty:
            print("⚠️  No interest rate data found")
            df = pd.DataFrame(columns=['series_id', 'series_name', 'date', 'rate'])
        
        output_path = 'data/federal_funds_rate.csv'
        df.to_csv(output_path, index=False)
        print(f"✅ Interest rates exported: {output_path} ({len(df)} rows)")
        
        if not df.empty:
            print(f"   Date range: {df['date'].min()} to {df['date'].max()}")
            print(f"   Rate range: {df['rate'].min():.4f}% to {df['rate'].max():.4f}%")
        
        return df
    except Exception as e:
        print(f"❌ Error exporting interest rates: {e}")
        return pd.DataFrame()

def export_fear_greed_index(conn):
    """Export Fear & Greed Index data to CSV"""
    print("\n😨 Exporting Fear & Greed Index data...")
    
    # Try both possible table names
    queries = [
        """
        SELECT value, value_classification, resolved_date as date, timestamp_unix
        FROM fgi_index 
        ORDER BY resolved_date ASC
        """,
        """
        SELECT value, value_classification, resolved_date as date, timestamp_unix
        FROM btc_fear_greed_index 
        ORDER BY resolved_date ASC
        """
    ]
    
    df = pd.DataFrame()
    
    for i, query in enumerate(queries):
        try:
            df = pd.read_sql(query, conn)
            table_name = "fgi_index" if i == 0 else "btc_fear_greed_index"
            print(f"✅ Found data in table: {table_name}")
            break
        except Exception as e:
            if i == len(queries) - 1:  # Last query failed
                print(f"⚠️  No Fear & Greed Index data found in either table")
                df = pd.DataFrame(columns=['value', 'value_classification', 'date', 'timestamp_unix'])
    
    output_path = 'data/fear_greed_index.csv'
    df.to_csv(output_path, index=False)
    print(f"✅ Fear & Greed Index exported: {output_path} ({len(df)} rows)")
    
    if not df.empty:
        print(f"   Date range: {df['date'].min()} to {df['date'].max()}")
        print(f"   Value range: {df['value'].min()} to {df['value'].max()}")
        print(f"   Classifications: {df['value_classification'].unique().tolist()}")
    
    return df


def main():
    """Main execution function"""
    print("🚀 Starting data export process...")
    print(f"Timestamp: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    
    # Create data directory if it doesn't exist
    os.makedirs('data', exist_ok=True)
    
    # Connect to database
    conn = connect_to_database()
    
    try:
        # Export individual datasets
        export_bitcoin_data(conn)
        export_interest_rates(conn)
        export_fear_greed_index(conn)
        
        print(f"\n🎉 Data export completed successfully!")
        print(f"📁 Check the 'data/' directory for CSV files")
        
    except Exception as e:
        print(f"❌ Unexpected error during export: {e}")
    finally:
        conn.close()
        print("🔌 Database connection closed")

if __name__ == "__main__":
    main() 