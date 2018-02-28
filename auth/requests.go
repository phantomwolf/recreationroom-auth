package auth

type registerRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type registerResponse struct {
	Err error `json:"error,omitempty"`
}

type unregisterRequest struct {
	UID      uint64 `json:"uid"`
	Password string `json:"password"`
	SID      string `json:"sid"`
}

type unregisterResponse struct {
	Err error `json:"error,omitempty"`
}

type loginRequest struct {
	NameOrEmail string `json:"name_or_email"`
	Password    string `json:"password"`
}

type loginResponse struct {
	Err error `json:"error,omitempty"`
}

type logoutRequest struct {
	UID uint64 `json:"uid"`
	SID string `json:"sid"`
}
