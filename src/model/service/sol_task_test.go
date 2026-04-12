package service

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/tidwall/gjson"
)

func TestSolClientHealthy(t *testing.T) {
	bodyData, err := SolRetryClient("getHealth", nil)
	if err != nil {
		t.Fatalf("SolRetryClient failed: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(bodyData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	status, ok := result["result"].(string)
	if !ok {
		t.Fatalf("Unexpected response format: %v", result)
	}

	t.Logf("RPC Health Status: %s", status)

	if status != "ok" {
		t.Errorf("Expected health status 'ok', got '%s'", status)
	}
}

func TestSolClientGetSignaturesForAddress(t *testing.T) {
	// Example wallet address (replace with actual test address)
	address := "2uFTf9TZ8gd7Kg6hkb79TxfaeNpaAgpJ8uVHguv2Yweu"

	bodyData, err := SolRetryClient("getSignaturesForAddress", []interface{}{address, map[string]interface{}{"commitment": "finalized", "limit": 100}})
	if err != nil {
		t.Fatalf("SolRetryClient failed: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(bodyData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	signatures, ok := result["result"].([]interface{})
	if !ok {
		t.Fatalf("Unexpected response format: %v", result)
	}

	t.Logf("Found %d signatures for address %s", len(signatures), address)

}

func TestSolClientGetTransaction(t *testing.T) {
	// Example transaction signature (replace with actual test signature)
	sig := "2aEoNykk4ZJ27C3y7EDJiQUc7GFnnsMe7ofFzB73swGL8kTxSBFCnwzWw3jzr3BND7k8hx15fZHUUAbG1XemNFe5"

	txData, err := SolRetryClient("getTransaction", []interface{}{sig, map[string]interface{}{"encoding": "jsonParsed", "commitment": "finalized"}})
	if err != nil {
		t.Fatalf("SolRetryClient failed: %v", err)
	}
	fmt.Printf("%v\n", string(txData))

	var result map[string]interface{}
	err = json.Unmarshal(txData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	txInfo, ok := result["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("Unexpected response format: %v", result)
	}

	t.Logf("Transaction Info for signature %s: %v", sig, txInfo)
}

func TestFindATAAddress(t *testing.T) {
	tests := []struct {
		name  string
		owner string
		mint  string
		want  string
	}{
		{
			name:  "RAY token ATA",
			owner: "2uFTf9TZ8gd7Kg6hkb79TxfaeNpaAgpJ8uVHguv2Yweu",
			mint:  "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R",
			want:  "GgmJrwuP946uV8qAwsnXxzYrJqEwW6eGnsVnQZFS5rp4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ata, err := FindATAAddress(tt.owner, tt.mint)
			if err != nil {
				t.Fatalf("FindATAAddress failed: %v", err)
			}

			t.Logf("Owner: %s", tt.owner)
			t.Logf("Mint: %s", tt.mint)
			t.Logf("ATA: %s", ata)

			if tt.want != "" && ata != tt.want {
				t.Errorf("Expected ATA %s, got %s", tt.want, ata)
			}
		})
	}
}

func TestMatchATAAddress(t *testing.T) {
	owner := "2uFTf9TZ8gd7Kg6hkb79TxfaeNpaAgpJ8uVHguv2Yweu"
	mint := "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R" // ray token
	expectedATA := "GgmJrwuP946uV8qAwsnXxzYrJqEwW6eGnsVnQZFS5rp4"

	ok := MatchAtaAddress(owner, mint, expectedATA)
	t.Logf("Owner: %s", owner)
	t.Logf("Mint: %s", mint)
	t.Logf("Expected ATA: %s", expectedATA)
	t.Logf("Match result: %v", ok)

	if !ok {
		t.Error("Expected ATA to match, but it didn't")
	}
}

func TestMatchUsdtAtaAddress(t *testing.T) {
	// Example wallet address (replace with actual test address)
	owner := "2uFTf9TZ8gd7Kg6hkb79TxfaeNpaAgpJ8uVHguv2Yweu"

	ata, err := FindATAAddress(owner, USDT_Mint)
	if err != nil {
		t.Fatalf("FindATAAddress failed: %v", err)
	}

	t.Logf("Owner: %s", owner)
	t.Logf("USDT Mint: %s", USDT_Mint)
	t.Logf("USDT ATA: %s", ata)

	ok := MatchUsdtAtaAddress(owner, ata)
	if !ok {
		t.Error("Expected USDT ATA to match")
	}
}

func TestMatchUsdcAtaAddress(t *testing.T) {
	// Example wallet address (replace with actual test address)
	owner := "2uFTf9TZ8gd7Kg6hkb79TxfaeNpaAgpJ8uVHguv2Yweu"

	ata, err := FindATAAddress(owner, USDC_Mint)
	if err != nil {
		t.Fatalf("FindATAAddress failed: %v", err)
	}

	t.Logf("Owner: %s", owner)
	t.Logf("USDC Mint: %s", USDC_Mint)
	t.Logf("USDC ATA: %s", ata)

	ok := MatchUsdcAtaAddress(owner, ata)
	if !ok {
		t.Error("Expected USDC ATA to match")
	}
}

func TestAdjustAmount(t *testing.T) {
	tests := []struct {
		name     string
		amount   uint64
		decimals int
		want     float64
	}{
		{
			name:     "USDT amount (6 decimals)",
			amount:   123456789,
			decimals: 6,
			want:     123.46,
		},
		{
			name:     "USDC amount (6 decimals)",
			amount:   1000000,
			decimals: 6,
			want:     1.0,
		},
		{
			name:     "SOL amount (9 decimals)",
			amount:   1000000000,
			decimals: 9,
			want:     1.0,
		},
		{
			name:     "Zero amount",
			amount:   0,
			decimals: 6,
			want:     0,
		},
		{
			name:     "Small amount",
			amount:   1,
			decimals: 6,
			want:     0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adjusted := ADJustAmount(tt.amount, tt.decimals)
			t.Logf("Raw amount: %d, Decimals: %d, Adjusted: %.2f", tt.amount, tt.decimals, adjusted)

			if adjusted != tt.want {
				t.Errorf("Expected %.2f, got %.2f", tt.want, adjusted)
			}
		})
	}
}

func TestParseTransferInfoFromInstruction_SplTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// SPL Token "transfer" (no mint in instruction, must look up from postTokenBalances)
	sig := "3tZTwLrvmiZ59h4UzyMHPd7DPux7t9eXZgkUvEfquaoSuERrPSRNzWuSHKQM2fbiCWFDGNqoLpu2kLZnfoegVpqN"
	txData, err := SolGetTransaction(sig)
	if err != nil {
		t.Fatalf("SolGetTransaction failed: %v", err)
	}

	instructions := gjson.GetBytes(txData, "result.transaction.message.instructions").Array()
	var found bool
	for _, inst := range instructions {
		info, err := ParseTransferInfoFromInstruction(inst, txData)
		if err != nil {
			t.Logf("parse error (ok to skip): %v", err)
			continue
		}
		if info == nil {
			continue
		}
		found = true
		t.Logf("SPL transfer: source=%s dest=%s mint=%s amount=%.6f raw=%d blockTime=%d",
			info.Source, info.Destination, info.Mint, info.Amount, info.RawAmount, info.BlockTime)

		if info.Mint == "" {
			t.Error("Expected mint to be resolved from postTokenBalances")
		}
		if info.RawAmount != 50000 {
			t.Errorf("Expected raw amount 50000, got %d", info.RawAmount)
		}
		if info.BlockTime == 0 {
			t.Error("Expected non-zero blockTime")
		}
	}
	if !found {
		t.Error("No transfer instruction found in transaction")
	}
}

func TestParseTransferInfoFromInstruction_TransferChecked(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// SPL Token "transferChecked" (has mint and tokenAmount in instruction)
	sig := "2aEoNykk4ZJ27C3y7EDJiQUc7GFnnsMe7ofFzB73swGL8kTxSBFCnwzWw3jzr3BND7k8hx15fZHUUAbG1XemNFe5"
	txData, err := SolGetTransaction(sig)
	if err != nil {
		t.Fatalf("SolGetTransaction failed: %v", err)
	}

	instructions := gjson.GetBytes(txData, "result.transaction.message.instructions").Array()
	var found bool
	for _, inst := range instructions {
		info, err := ParseTransferInfoFromInstruction(inst, txData)
		if err != nil {
			t.Logf("parse error (ok to skip): %v", err)
			continue
		}
		if info == nil {
			continue
		}
		found = true
		t.Logf("TransferChecked: source=%s dest=%s mint=%s amount=%.6f raw=%d blockTime=%d",
			info.Source, info.Destination, info.Mint, info.Amount, info.RawAmount, info.BlockTime)

		if info.Mint != USDT_Mint {
			t.Errorf("Expected USDT mint %s, got %s", USDT_Mint, info.Mint)
		}
		if info.RawAmount != 300000 {
			t.Errorf("Expected raw amount 300000, got %d", info.RawAmount)
		}
		if info.Amount != 0.3 {
			t.Errorf("Expected amount 0.3, got %f", info.Amount)
		}
		if info.BlockTime == 0 {
			t.Error("Expected non-zero blockTime")
		}
	}
	if !found {
		t.Error("No transferChecked instruction found in transaction")
	}
}

func TestParseTransferInfoFromInstruction_SystemTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// System program SOL transfer
	sig := "5pNMonUBvLVpxXTmyd5CGVBs49W6781g2ACnrCXhbmtz58KENYA7HSqu6hQkQweg3qQboRd8WAscphNAtiq9UtZZ"
	txData, err := SolGetTransaction(sig)
	if err != nil {
		t.Fatalf("SolGetTransaction failed: %v", err)
	}

	instructions := gjson.GetBytes(txData, "result.transaction.message.instructions").Array()
	transferCount := 0
	for _, inst := range instructions {
		info, err := ParseTransferInfoFromInstruction(inst, txData)
		if err != nil {
			t.Logf("parse error (ok to skip): %v", err)
			continue
		}
		if info == nil {
			continue
		}
		transferCount++
		t.Logf("System transfer #%d: source=%s dest=%s mint=%s amount=%.9f raw=%d blockTime=%d",
			transferCount, info.Source, info.Destination, info.Mint, info.Amount, info.RawAmount, info.BlockTime)

		if info.Mint != "SOL" {
			t.Errorf("Expected mint SOL, got %s", info.Mint)
		}
		if info.RawAmount == 0 {
			t.Error("Expected non-zero raw amount")
		}
		if info.BlockTime == 0 {
			t.Error("Expected non-zero blockTime")
		}
	}
	if transferCount == 0 {
		t.Error("No system transfer instruction found")
	}
	t.Logf("Found %d system transfers", transferCount)
}
