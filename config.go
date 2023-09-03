package main

// Config holds the configuration for Idicon service.
type Config struct {
	Defaults Defaults     `toml:"defaults"`
	Users    []UserConfig `toml:"users"`
}

// Defaults holds default configuration values to be used as defaults for all users.
type Defaults struct {
	ColorScheme string `toml:"color-scheme"`
	Pattern     string `toml:"pattern"`
}

// UserConfig holds user specific configuration.
// ID is the id od the user in plain text,
// Alias is the alias to be used to generate the id icon.
type UserConfig struct {
	ID          string `toml:"id"`
	Alias       string `toml:"alias"`
	ColorScheme string `toml:"color-scheme"`
	Pattern     string `toml:"pattern"`
	Redirect    string `toml:"redirect"`
}
