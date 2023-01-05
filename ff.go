package ff

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	initLock sync.RWMutex
	randInit bool
)

const (
	minValue = 0.0
	maxValue = 1.0
)

// FlagStorage represents the required functionality of
// a storage solution for feature flags.
type FlagStorage interface {
	GetFlag(id string) (*FeatureFlag, error)
	PutFlag(ff *FeatureFlag) error
	DeleteFlag(id string) error
}

// FeatureFlag represents a feature flag with an identifier and a value
type FeatureFlag struct {
	ID    string
	Value float64
}

// NewFeatureFlag returns a new feature flag
func NewFeatureFlag(id string, value float64) (*FeatureFlag, error) {
	if value < minValue || value > maxValue {
		return nil, fmt.Errorf("Feature flag value must be between 0.0 and 1.0 inclusive, got %f", value)
	}
	return &FeatureFlag{ID: id, Value: value}, nil
}

// IsEnabled returns true with probability equal to the value of the flag.
func (ff *FeatureFlag) IsEnabled() bool {
	ensureRandIsSeeded()

	return rand.Float64() < ff.Value
}

// IsEnabledForUser returns true if the treatment represented by
// this flag should be applied for a given user identifier. If the
// flag's value does not change, this will always return the same
// result for any given user identifier.
func (ff *FeatureFlag) IsEnabledForUser(userID string) bool {
	return toUniformFloat64(fmt.Sprintf("%s%s", ff.ID, userID)) < ff.Value
}

// leverages the uniform distribution of cryptographic hashing
// functions in order to produce a random, but deterministically
// repeatable float64 for the provided string
func toUniformFloat64(str string) float64 {
	h := sha256.New()
	h.Write([]byte(str))
	hash := h.Sum(nil)
	return float64(binary.LittleEndian.Uint64(hash)) / float64((1<<64)-1)
}

// ensures that the random number generator has been seeded
func ensureRandIsSeeded() {
	initLock.RLock()
	seeded := randInit
	initLock.RUnlock()

	if seeded {
		return
	}

	rand.Seed(time.Now().UTC().UnixNano())

	initLock.Lock()
	randInit = true
	initLock.Unlock()
}
