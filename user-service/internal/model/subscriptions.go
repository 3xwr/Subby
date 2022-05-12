package model

type SubSuccessResponse struct {
	Subscriber string `json:"subscriber"`
	Subscribed string `json:"subscribed"`
	SubSuccess bool   `json:"sub_success"`
}

type UnsubSuccessResponse struct {
	Subscriber   string `json:"subscriber"`
	Unsubscribed string `json:"unsubscribed"`
	UnsubSuccess bool   `json:"unsub_success"`
}
