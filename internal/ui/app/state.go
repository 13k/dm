package app

const (
	stateShowForm appState = iota
	stateShowDocument
	stateError
	stateDocumentSaved
)

type appState int

func (s appState) String() string {
	switch s {
	case stateError:
		return "error"
	case stateShowForm:
		return "show-form"
	case stateShowDocument:
		return "show-document"
	case stateDocumentSaved:
		return "document-saved"
	default:
		return "<unknown>"
	}
}
