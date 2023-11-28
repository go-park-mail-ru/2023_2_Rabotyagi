package statuses

const (
	// StatusResponseSuccessful uses for indicates successful status of request
	StatusResponseSuccessful = 2000

	// StatusRedirectAfterSuccessful uses when need redirect in frontend user to another resource
	StatusRedirectAfterSuccessful = 3003

	MinValueClientError = 4000
	// StatusBadFormatRequest uses when get bad request from frontend and errors with this status need frontend developers
	StatusBadFormatRequest = 4000

	// StatusBadContentRequest uses when user has entered incorrect data and needs to show him this error
	StatusBadContentRequest = 4400
	MaxValueClientError     = 4999

	// StatusInternalServer uses for indicates internal error status in server
	StatusInternalServer = 5000
)
