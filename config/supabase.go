package config

import (
	"fmt"
	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
	"os"
)

var SupabaseClient *supabase.Client

func InitSupabase() error {

	// Add debug logging
	fmt.Printf("Initializing Supabase client...\n")

	// The URL should include the full path
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_ANON_KEY")

	fmt.Printf("Supabase URL: %s\n", supabaseURL)
	fmt.Printf("Supabase Key length: %d\n", len(supabaseKey))

	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return fmt.Errorf("failed to create Supabase client: %w", err)
	}

	SupabaseClient = client
	fmt.Printf("Supabase client initialized successfully\n")

	return nil
}
func GetSupabaseClient() *supabase.Client {
	return SupabaseClient
}

var dbClient *postgrest.Client

func InitPostgres() error {
	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_ANON_KEY")

	headers := map[string]string{
		"apikey":        supabaseKey,
		"Authorization": fmt.Sprintf("Bearer %s", supabaseKey),
		"Content-Type":  "application/json",
		"Prefer":        "return=minimal",
	}

	client := postgrest.NewClient(fmt.Sprintf("%s/rest/v1", supabaseUrl), "public", headers)

	fmt.Printf("DB Client created with URL: %s\n", supabaseUrl)

	if client.ClientError != nil {
		return fmt.Errorf("failed to create Postgres client: %w", client.ClientError)
	}

	dbClient = client
	return nil
}

func GetDBClient() *postgrest.Client {
	return dbClient
}
