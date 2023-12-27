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

func ConvertToIntStatus(status string) uint8 {
	switch {
	case status == StatusPaymentPending:
		return IntStatusPremiumPending
	case status == StatusPaymentWaiting:
		return IntStatusPremiumWaiting
	case status == StatusPaymentSucceeded:
		return IntStatusPremiumSucceeded
	case status == StatusPaymentCanceled:
		return IntStatusPremiumCanceled
	default:
		return IntStatusPremiumNot
	}
}
