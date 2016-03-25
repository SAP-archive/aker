package uuid

import (
	"crypto/rand"
	"fmt"
)

// UUID represents a universally unique identifier
type UUID [16]byte

// String returns the canonical form of the UUID.
// The canonical form is represented by 32 lowercase hexadecimal digits,
// displayed in five groups separated by hyphens, in the form 8-4-4-4-12.
func (u UUID) String() string {
	// canonical form: 8-4-4-4-12
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		u[0:4],
		u[4:6],
		u[6:8],
		u[8:10],
		u[10:])
}

// Random generates random UUID.
// It uses CSPRNG.
func Random() (UUID, error) {
	u := UUID{}
	_, err := rand.Read(u[:])
	if err != nil {
		return UUID{}, err
	}
	return u, nil
}
