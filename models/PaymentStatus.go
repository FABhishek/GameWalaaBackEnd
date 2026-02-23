package models

type PaymentStatus struct {
	OrderCreationId   string `json:"orderCreationId"`
	RazorpayPaymentId string `json:"razorpayPaymentId"`
	RazorpayOrderId   string `json:"razorpayOrderId"`
	RazorpaySignature string `json:"razorpaySignature"`
}

type PaymentAndGameStatus struct {
	PaymentDetails PaymentStatus `json:"paymentDetails"`
	GameStatus     GameStatus    `json:"gameStatus"`
}
