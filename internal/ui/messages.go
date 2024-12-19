package ui

var _ error = (*ErrorMsg)(nil)

type ErrorMsg struct { //nolint:errname // UI message representing an error
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

type CopyToClipboardMsg struct {
	Text string
}

func NewClipboardDocumentMsg(text string) *CopyToClipboardMsg {
	return &CopyToClipboardMsg{Text: text}
}

type ClipboardWrittenMsg struct{}

func NewClipboardWrittenMsg() *ClipboardWrittenMsg {
	return &ClipboardWrittenMsg{}
}

type PublishSlackMsg struct {
	Channel string
	Message string
}

func NewPublishSlackMsg(channel, msg string) *PublishSlackMsg {
	return &PublishSlackMsg{
		Channel: channel,
		Message: msg,
	}
}

type PublishedSlackMsg struct {
	Channel   string
	Timestamp string
}

func NewPublishedSlackMsg(channel, timestamp string) *PublishedSlackMsg {
	return &PublishedSlackMsg{
		Channel:   channel,
		Timestamp: timestamp,
	}
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
