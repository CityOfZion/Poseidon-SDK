package api

// REFER TO README
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBalanceOfUsedAddress(t *testing.T) {
	nw := API{}
	response := nw.CheckBalance("AXym3Qc9mRbKF5HWAtmJUoHvB3F8brEdLe")
	assert.NotEqual(t, "not found", response.Address) // Test will fail, if the address has not been used yet
}

func TestGetBalanceOfNewAddress(t *testing.T) {
	nw := API{}
	address := "AT9DbViHVMy52FZnw2XDjT3TMkMNYyZZmm"
	response := nw.CheckBalance(address)
	assert.Equal(t, address, response.Address) // Test will fail, if the address has not been used yet
}
