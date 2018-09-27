package redis

// Migration constructed by several up and down jobs.
type Migration struct {
	Up   []Job `yaml:"up"`
	Down []Job `yaml:"down"`
}

// Job is a job for specific keys.
type Job struct {
	// Migrant will use `SCAN {pattern}` to get keys.
	Pattern string `yaml:"pattern"`
	// Actions are the actions templates behaved on keys.
	Actions []string `yaml:"actions"`
}
