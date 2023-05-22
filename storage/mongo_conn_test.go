package storage

import (
	"testing"
)

func TestConnectToMongoDB(t *testing.T) {

	// Call the function being tested
	client, err := ConnectToMongoDB()

	// Check for errors
	if err != nil {
		t.Errorf("ConnectToMongoDB returned an error: %v", err)
	}

	// Check that the client is not nil
	if client == nil {
		t.Errorf("ConnectToMongoDB returned a nil client")
	}

}
