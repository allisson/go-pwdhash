package argon2

import "errors"

// PolicyParams captures the Argon2 tuning knobs for a policy preset.
type PolicyParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
}

// ParamsForPolicy returns the Argon2 parameters associated with the policy id.
func ParamsForPolicy(p int) (PolicyParams, error) {
	switch p {
	case 0: // Interactive
		return PolicyParams{
			Memory:      64 * 1024, // 64 MB
			Iterations:  3,
			Parallelism: 4,
		}, nil

	case 1: // Moderate
		return PolicyParams{
			Memory:      128 * 1024, // 128 MB
			Iterations:  4,
			Parallelism: 4,
		}, nil

	case 2: // Sensitive
		return PolicyParams{
			Memory:      256 * 1024, // 256 MB
			Iterations:  5,
			Parallelism: 8,
		}, nil
	}

	return PolicyParams{}, errors.New("unknown password policy")
}
