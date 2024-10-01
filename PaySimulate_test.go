package PaySimulate

import (
	"fmt"
	"testing"
)

// ANSI escape codes for green text
const (
	green = "\033[32m"
	reset = "\033[0m"
	red   = "\033[31m"
)

// Test function
func TestCardErrors(t *testing.T) {
	cardTest := CardDetails{"12345678901234567", "John Doe", "12/23", "123", "123 Main St", "US"}

	_, err := PaymentSuccess(cardTest, true)
	if err == nil || err.Error() != "card number must be 16 digits" {
		t.Errorf("Expected error 'card number must be 16 digits' but got %v", err)
	} else {
		fmt.Printf("%sAPPROVED: testCardErrors err->%s%s\n", green, err, reset)
	}

	cardTest = CardDetails{"1234567890123456", "Jo", "12/23", "123", "123 Main St", "US"}
	_, err = PaymentSuccess(cardTest, true)
	if err == nil || err.Error() != "name must be at least 3 characters" {
		t.Errorf("Expected error 'name must be at least 3 characters' but got %v", err)
	} else {
		fmt.Printf("%sAPPROVED: testCardErrors err->%s%s\n", green, err, reset)
	}

	cardTest = CardDetails{"1234567890123456", "John Doe", "12/2", "123", "123 Main St", "US"}
	_, err = PaymentSuccess(cardTest, true)
	if err == nil || err.Error() != "invalid expiration date" {
		t.Errorf("Expected error 'invalid expiration date' but got %v", err)
	} else {
		fmt.Printf("%sAPPROVED: testCardErrors err->%s%s\n", green, err, reset)
	}

	cardTest = CardDetails{"1234567890123456", "John Doe", "12/24", "1234", "123 Main St", "US"}
	_, err = PaymentSuccess(cardTest, true)
	if err == nil || err.Error() != "invalid CVV" {
		t.Errorf("Expected error 'invalid CVV' but got %v", err)
	} else {
		fmt.Printf("%sAPPROVED: testCardErrors err->%s%s\n", green, err, reset)
	}

	cardTest = CardDetails{"1234567890123456", "John Doe", "12/24", "123", "123", "US"}
	_, err = PaymentSuccess(cardTest, true)
	if err == nil || err.Error() != "invalid billing address" {
		t.Errorf("Expected error 'invalid billing address' but got %v", err)
	} else {
		fmt.Printf("%sAPPROVED: testCardErrors err->%s%s\n", green, err, reset)
	}

	cardTest = CardDetails{"1234567890123456", "John Doe", "12/24", "123", "123 Main St", "USA"}
	_, err = PaymentSuccess(cardTest, true)
	if err == nil || err.Error() != "invalid country code" {
		t.Errorf("Expected error 'invalid country code' but got %v", err)
	} else {
		fmt.Printf("%sAPPROVED: testCardErrors err->%s%s\n", green, err, reset)
	}
}

func TestPaymentSuccess(t *testing.T) {
	cardTest := CardDetails{"1234567890123456", "John Doe", "12/23", "123", "123 Main St", "US"}

	result, err := PaymentSuccess(cardTest, true)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if result == "" {
		t.Errorf("Expected payment success but got failure")
	} else {
		fmt.Printf("%sAPPROVED: TestPaymentSuccess%s\n", green, reset)
	}
}

func TestWebServer(t *testing.T) {
	StartWebServer(":8080")
}
