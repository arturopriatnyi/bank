package banking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMoney(t *testing.T) {
	for name, tc := range map[string]struct {
		units     int32
		nanos     int32
		currency  Currency
		wantMoney Money
		wantErr   bool
	}{
		"MoneyIsCreated_1.01USD": {
			units:    1,
			nanos:    1,
			currency: USD,
			wantMoney: Money{
				value:    101,
				currency: USD,
			},
			wantErr: false,
		},
		"NegativeNanosError": {
			units:     0,
			nanos:     -1,
			currency:  USD,
			wantMoney: Money{},
			wantErr:   true,
		},
		"NanosError": {
			units:     0,
			nanos:     101,
			currency:  USD,
			wantMoney: Money{},
			wantErr:   true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m, err := NewMoney(tc.units, tc.nanos, tc.currency)

			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.wantMoney, m)
		})
	}
}
