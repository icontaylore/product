package kafka

import (
	"encoding/base64"
	"math/big"
)

func decodeDebeziumDecimal(encoded string) (*big.Float, error) {
	if encoded == "" {
		return nil, nil
	}

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	// Debezium хранит Decimal как big-endian байты
	intVal := new(big.Int).SetBytes(data)

	// У нас scale=2 (из схемы), поэтому делим на 100
	result := new(big.Float).SetInt(intVal)
	result.Quo(result, big.NewFloat(100))

	return result, nil
}
