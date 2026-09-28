package settings

// derefStr — dereference *string dengan default
func derefStr(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}
