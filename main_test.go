package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadSecretsCSV(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "test-csv")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a dummy CSV file
	csvContent := `name,value
my-secret-1,value1
my-secret-2,value2`
	csvFilePath := filepath.Join(tempDir, "secrets.csv")
	err = os.WriteFile(csvFilePath, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("failed to write dummy csv file: %v", err)
	}

	// Test case 1: Successful read
	secrets, err := readSecretsCSV(csvFilePath)
	if err != nil {
		t.Errorf("readSecretsCSV() returned an unexpected error: %v", err)
	}

	if len(secrets) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(secrets))
	}

	if secrets[0].Name != "my-secret-1" || secrets[0].Value != "value1" {
		t.Errorf("unexpected secret data: got %+v", secrets[0])
	}

	if secrets[1].Name != "my-secret-2" || secrets[1].Value != "value2" {
		t.Errorf("unexpected secret data: got %+v", secrets[1])
	}

	// Test case 2: File not found
	_, err = readSecretsCSV("nonexistent.csv")
	if err == nil {
		t.Error("readSecretsCSV() did not return an error for a nonexistent file")
	}

	// Test case 3: Invalid header
	invalidHeaderContent := `invalid_name,invalid_value
my-secret-1,value1`
	invalidHeaderFilePath := filepath.Join(tempDir, "invalid_header.csv")
	err = os.WriteFile(invalidHeaderFilePath, []byte(invalidHeaderContent), 0644)
	if err != nil {
		t.Fatalf("failed to write dummy csv file: %v", err)
	}

	_, err = readSecretsCSV(invalidHeaderFilePath)
	if err == nil {
		t.Error("readSecretsCSV() did not return an error for an invalid header")
	}
}
