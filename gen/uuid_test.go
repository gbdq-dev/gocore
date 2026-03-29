package gen

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerators_ReturnNonNilAndNonNilUUID(t *testing.T) {
	tests := []struct {
		name string
		fn   func() UUIDGenerator
	}{
		{name: "UUID", fn: UUID},
		{name: "Generator", fn: Generator},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.fn()
			if g == nil {
				t.Fatalf("%s() returned nil generator", tc.name)
			}

			got := g.Next()
			if got == uuid.Nil {
				t.Fatalf("%s().Next() returned uuid.Nil", tc.name)
			}
		})
	}
}

func TestUUIDGeneratorNext_NilGenerator(t *testing.T) {
	var g UUIDGenerator

	got := g.Next()
	if got != uuid.Nil {
		t.Fatalf("Next() = %v, want %v", got, uuid.Nil)
	}
}

func TestUUIDGeneratorNext_UsesUnderlyingFunction(t *testing.T) {
	want := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	g := UUIDGenerator(func() uuid.UUID {
		return want
	})

	got := g.Next()
	if got != want {
		t.Fatalf("Next() = %v, want %v", got, want)
	}
}
