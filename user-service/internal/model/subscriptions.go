package model

import "github.com/google/uuid"

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

type CheckSubscriptionResponse struct {
	Subscribed bool `json:"subscribed"`
}

type CheckSubscriptionRequest struct {
	Subscriber uuid.UUID `json:"subscriber"`
	Subscribed uuid.UUID `json:"subscribed"`
}
