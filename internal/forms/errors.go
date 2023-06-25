package forms

type errors map[string][]string

// Add adds an error message to a given field
func (e errors) Add(field string, msg string) {
	e[field] = append(e[field], msg)
}

// Get returns the first errors message for a field
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
