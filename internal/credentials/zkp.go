package credentials

import (
	"fmt"
	"math"
	"math/big"

	"github.com/0xdecaf/zkrp/bulletproofs"
)

// ZKProof holds a Bulletproof range proof with its minimum threshold
type ZKProof struct {
	Proof  bulletproofs.BulletProof `json:"proof"`
	MinAge uint64                   `json:"minAge"`
}

// nextPowerOfTwo returns the smallest power-of-two >= n
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

// GenerateAgeProof creates a zero-knowledge proof that 'age' >= 'minAge'
// using the bulletproofs subpackage. It chooses a range end of 2^e where e
// is the next power-of-two exponent that covers (age - minAge + 1).
func GenerateAgeProof(age uint64, minAge uint64) (*ZKProof, error) {
	if age < minAge {
		return nil, fmt.Errorf("age %d is below minimum %d", age, minAge)
	}
	// Compute diff = age - minAge
	diff := new(big.Int).SetUint64(age - minAge)
	// Compute required range size = diff + 1
	size := age - minAge + 1
	// Determine exponent bits0 = ceil(log2(size))
	bits0 := uint64(math.Ceil(math.Log2(float64(size))))
	// Compute exponent must be power-of-two
	e := nextPowerOfTwo(bits0)
	// Range end = 2^e
	rangeEnd := uint64(1) << e
	// Setup bulletproof params for this range end
	params, err := bulletproofs.Setup(int64(rangeEnd))
	if err != nil {
		return nil, fmt.Errorf("setting up bulletproof params: %w", err)
	}
	// Prove that diff in [0, rangeEnd-1]
	proof, err := bulletproofs.Prove(diff, params)
	if err != nil {
		return nil, fmt.Errorf("proving range: %w", err)
	}
	return &ZKProof{Proof: proof, MinAge: minAge}, nil
}

// VerifyAgeProof verifies that the provided proof demonstrates age ≥ MinAge
func VerifyAgeProof(z *ZKProof) error {
	ok, err := z.Proof.Verify()
	if err != nil {
		return fmt.Errorf("verifying range proof: %w", err)
	}
	if !ok {
		return fmt.Errorf("range proof invalid: age < %d", z.MinAge)
	}
	return nil
}
