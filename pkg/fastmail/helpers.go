package fastmail

// isEnabledToString returns a string representation of the enabled state.
func isEnabledToString(enabled bool) string {
	if enabled {
		return "enabled"
	}

	return "disabled"
}
