package delivery

import (
	"fmt"
	"sync"
)

const (
	startLenMapIdempotencyPayment = 100
	maxLenKeyIdempotencyPayment   = 64
)

func generateString() string {
	resultString := ""

	for len(resultString) < maxLenKeyIdempotencyPayment {
		a := new(int8)

		resultString += fmt.Sprintf("%v", a)
	}

	return resultString[:maxLenKeyIdempotencyPayment]
}

type KeyIdempotencyPayment string

type MapIdempotencePayment struct {
	mapIdempotence map[MetadataPayment]KeyIdempotencyPayment
	mu             *sync.RWMutex
}

func NewMapIdempotence() *MapIdempotencePayment {
	return &MapIdempotencePayment{
		mapIdempotence: make(map[MetadataPayment]KeyIdempotencyPayment, startLenMapIdempotencyPayment),
		mu:             &sync.RWMutex{},
	}
}

func (m *MapIdempotencePayment) AddPayment(metadata *MetadataPayment) KeyIdempotencyPayment {
	m.mu.RLock()

	keyIdempotencyPayment, ok := m.mapIdempotence[*metadata]
	if ok {
		return keyIdempotencyPayment
	}

	m.mu.RUnlock()

	keyIdempotencyPayment = KeyIdempotencyPayment(generateString())

	m.mu.Lock()
	m.mapIdempotence[*metadata] = keyIdempotencyPayment
	m.mu.Unlock()

	return keyIdempotencyPayment
}
