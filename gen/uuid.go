package gen

import (
	"github.com/google/uuid"
)

// UUIDGenerator produces random unique identifiers (UUIDs).
type UUIDGenerator func() uuid.UUID

// UUID returns a UUID generator.
func UUID() UUIDGenerator {
	return Generator()
}

// Next returns the next unique identifier.
func (g UUIDGenerator) Next() uuid.UUID {
	if g == nil {
		return uuid.Nil
	}

	return g()
}

// RealUUIDGenerator returns a generator that produces new UUIDs via uuid.New.
func Generator() UUIDGenerator {
	return func() uuid.UUID {
		return uuid.New()
	}
}
