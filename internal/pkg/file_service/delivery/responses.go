package delivery

type ResponseURLBody struct {
	URL string `json:"url"`
}

type ResponseURL struct {
	Status int             `json:"status"`
	Body   ResponseURLBody `json:"body"`
}

func NewResponseURL(status int, URL string) *ResponseURL {
	return &ResponseURL{Status: status, Body: ResponseURLBody{URL: URL}}
}
