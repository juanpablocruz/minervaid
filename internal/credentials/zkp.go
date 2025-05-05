package credentials

import (
	"fmt"
	"math"
	"math/big"

	"github.com/0xdecaf/zkrp/bulletproofs"
)

// RangeProof holds a Bulletproof range proof with its minimum threshold and range end.
type RangeProof struct {
	Type  string                   `json:"type"`
	Min   uint64                   `json:"min"`
	Range uint64                   `json:"range"`
	Proof bulletproofs.BulletProof `json:"proof"`
}

// nextPowerOfTwo returns the smallest power-of-two >= n.
func nextPowerOfTwo(n uint64) uint64 {
	if n == 0 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	n++
	return n
}

// GenerateRangeProof creates a zero-knowledge proof that 'value' >= 'min' using Bulletproofs.
// It computes a range end of 2^e where e is the smallest exponent such that 2^e >= (value - min + 1).
func GenerateRangeProof(value, min uint64) (*RangeProof, error) {
	if value < min {
		return nil, fmt.Errorf("value %d is below minimum %d", value, min)
	}

	// Compute difference: value - min
	diff := new(big.Int).SetUint64(value - min)
	// Range size = diff + 1
	size := value - min + 1
	// Calculate exponent bits0 = ceil(log2(size))
	bits0 := uint64(math.Ceil(math.Log2(float64(size))))
	// Next power-of-two exponent = nextPowerOfTwo(bits0)
	e := nextPowerOfTwo(bits0)
	// Range end = 2^e
	rangeEnd := uint64(1) << e

	// Setup bulletproof parameters for this range end
	params, err := bulletproofs.Setup(int64(rangeEnd))
	if err != nil {
		return nil, fmt.Errorf("setting up bulletproof params: %w", err)
	}
	// Prove that diff in [0, rangeEnd-1]
	proof, err := bulletproofs.Prove(diff, params)
	if err != nil {
		return nil, fmt.Errorf("proving range: %w", err)
	}

	return &RangeProof{
		Type:  "BulletproofRangeProof",
		Min:   min,
		Range: rangeEnd,
		Proof: proof,
	}, nil
}

// VerifyRangeProof verifies that the provided RangeProof demonstrates value >= Min.
func VerifyRangeProof(r *RangeProof) (err error) {
	// Recover from panics during proof verification
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("panic verifying range proof: %v", rec)
		}
	}()

	// Perform the cryptographic verification
	ok, verr := r.Proof.Verify()
	if verr != nil {
		return fmt.Errorf("verifying range proof: %w", verr)
	}
	if !ok {
		return fmt.Errorf("range proof invalid: value < %d", r.Min)
	}
	return nil
}
