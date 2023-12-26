package statuses

const (
	IntStatusPremiumNot = uint8(0)

	IntStatusPremiumPending = uint8(1)
	StatusPaymentPending    = "pending"

	IntStatusPremiumWaiting = uint8(2)
	StatusPaymentWaiting    = "waiting_for_capture"

	IntStatusPremiumSucceeded = uint8(3)
	StatusPaymentSucceeded    = "succeeded"

	IntStatusPremiumCanceled = uint8(4)
	StatusPaymentCanceled    = "canceled"
)

func IsStatusPaymentSuccessful(status string) bool {
	return status == StatusPaymentWaiting || status == StatusPaymentSucceeded
}

func IsIntStatusPremiumSuccessful(status uint8) bool {
	return status == IntStatusPremiumWaiting || status == IntStatusPremiumSucceeded
}
