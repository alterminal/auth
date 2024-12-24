package api

type IdGenerator struct {
	Type string `json:"type"`
}

type CreateAccountRequest struct {
	Namespace   string      `json:"namespace" binding:"required"`
	ID          string      `json:"id"`
	IdGenerator IdGenerator `json:"idGenerator"`
	Account     string      `json:"account"`
	Email       string      `json:"email"`
	PhoneRegion string      `json:"phoneRegion"`
	PhoneNumber string      `json:"phoneNumber"`
	Password    string      `json:"password"`
}

type CreateSessionRequest struct {
	Namespace   string `json:"namespace" binding:"required"`
	Account     string `json:"account"`
	Password    string `json:"password" binding:"required"`
	Email       string `json:"email"`
	PhoneRegion string `json:"phoneRegion"`
	PhoneNumber string `json:"phoneNumber"`
	Idby        string `json:"idby"`
}

type SetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}
