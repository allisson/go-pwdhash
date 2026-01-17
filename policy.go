package pwdhash

// Policy represents a password hashing strength preset.
type Policy int

const (
	// PolicyInteractive balances CPU cost with low latency for login flows.
	PolicyInteractive Policy = iota
	// PolicyModerate increases cost for privileged accounts or admin portals.
	PolicyModerate
	// PolicySensitive maximizes cost for high-value secrets.
	PolicySensitive
)
