package subscription

import "fmt"

// Unit represents a currency unit (e.g., USD, EUR).
type Currency string

// Predefined currency units (constants for common usage).
const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	NGN Currency = "NGN" // Nigerian Naira
	// Add other currencies as needed
)

// IsValid checks if the currency unit is one of the predefined valid units.
// You might want a more comprehensive list or external data for this in a real app.
func (u Currency) IsValid() bool {
	switch u {
	case USD, EUR, GBP, NGN:
		return true
	default:
		return false
	}
}

// String implements the fmt.Stringer interface.
func (u Currency) String() string {
	return string(u)
}

// ExchangeRates is a simple map for demonstration.
// In a real application, you'd fetch these from a reliable external source.
// This example uses a very simplified, fixed exchange rate.
// For production, consider a service that provides real-time, bidirectional rates.
var ExchangeRates = map[Currency]map[Currency]float64{
	USD: {
		EUR: 0.92,    // 1 USD = 0.92 EUR
		GBP: 0.79,    // 1 USD = 0.79 GBP
		NGN: 1470.00, // 1 USD = 1470 NGN (Example rate, fluctuates)
	},
	EUR: {
		USD: 1.09,    // 1 EUR = 1.09 USD
		GBP: 0.86,    // 1 EUR = 0.86 GBP
		NGN: 1600.00, // 1 EUR = 1600 NGN
	},
	// Add more conversion rates for other currencies
}

// Convert converts an amount from a source currency to a target currency.
// This is a utility function, not part of the Price object itself.
func ConvertCurrency(amount float64, fromUnit, toUnit Currency) (float64, error) {
	if fromUnit == toUnit {
		return amount, nil // No conversion needed
	}

	if !fromUnit.IsValid() {
		return 0, fmt.Errorf("invalid source currency unit: %s", fromUnit)
	}
	if !toUnit.IsValid() {
		return 0, fmt.Errorf("invalid target currency unit: %s", toUnit)
	}

	rates, ok := ExchangeRates[fromUnit]
	if !ok {
		return 0, fmt.Errorf("no exchange rates found for source currency: %s", fromUnit)
	}

	rate, ok := rates[toUnit]
	if !ok {
		return 0, fmt.Errorf("no exchange rate found from %s to %s", fromUnit, toUnit)
	}

	return amount * rate, nil
}
