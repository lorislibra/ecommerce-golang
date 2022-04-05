package dtos

type JwtResponse struct {
	AccessToken string `json:"access_token"`
	ExpireIn    int64  `json:"expire_in"`
}
