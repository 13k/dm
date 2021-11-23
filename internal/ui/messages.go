package ui

type ErrorMsg struct {
	Err error
}

func NewErrorMsg(err error) ErrorMsg {
	return ErrorMsg{Err: err}
}

func (e ErrorMsg) Error() string {
	return e.Err.Error()
}

type FormSubmittedMsg struct {
	Entries []string
}

type DocumentRenderedMsg struct {
	Body        string
	BodyColored string
}

type SaveDocumentMsg struct {
	Body string
}

type DocumentSavedMsg struct{}
