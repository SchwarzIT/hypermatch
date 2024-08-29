package hypermatch

// Property represents a property and a slice of values inside an event. An Event is defined as a slice of Property objects.
type Property struct {
	Path   string
	Values []string
}
