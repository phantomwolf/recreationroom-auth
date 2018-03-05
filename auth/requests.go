package auth

type logoutRequest struct {
	UID uint64 `json:"uid"`
	SID string `json:"sid"`
}
