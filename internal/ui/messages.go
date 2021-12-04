package ui

var _ error = (*ErrorMsg)(nil)

type ErrorMsg struct {
	Err error
}

func NewErrorMsg(err error) *ErrorMsg {
	return &ErrorMsg{Err: err}
}

func (e *ErrorMsg) Error() string {
	return e.Err.Error()
}

type FormSubmittedMsg struct {
	Entries []string
}

func NewFormSubmittedMsg(entries []string) *FormSubmittedMsg {
	return &FormSubmittedMsg{Entries: entries}
}

type DocumentParsedMsg struct {
	Entries []string
}

func NewDocumentParsedMsg(entries []string) *DocumentParsedMsg {
	return &DocumentParsedMsg{Entries: entries}
}

type DocumentRenderedMsg struct {
	Body        string
	BodyColored string
}

func NewDocumentRenderedMsg(body, bodyColored string) *DocumentRenderedMsg {
	return &DocumentRenderedMsg{Body: body, BodyColored: bodyColored}
}

type ClipboardDocumentMsg struct {
	Body string
}

func NewClipboardDocumentMsg(body string) *ClipboardDocumentMsg {
	return &ClipboardDocumentMsg{Body: body}
}

type ClipboardWrittenMsg struct{}

func NewClipboardWrittenMsg() *ClipboardWrittenMsg {
	return &ClipboardWrittenMsg{}
}

type SaveDocumentMsg struct {
	Body string
}

func NewSaveDocumentMsg(body string) *SaveDocumentMsg {
	return &SaveDocumentMsg{Body: body}
}

type DocumentSavedMsg struct{}

func NewDocumentSavedMsg() *DocumentSavedMsg {
	return &DocumentSavedMsg{}
}
