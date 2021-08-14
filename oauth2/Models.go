package oauth2

type ClaimSet struct {
	Iss   string   `json:"iss"`             // email address of the client_id of the application making the access token request
	Scope []string `json:"scope,omitempty"` // space-delimited list of the permissions the application requests
	Aud   string   `json:"aud"`             // descriptor of the intended target of the assertion (Optional).
	Exp   int64    `json:"exp"`             // the expiration time of the assertion (seconds since Unix epoch)
	Iat   int64    `json:"iat"`             // the time the assertion was issued (seconds since Unix epoch)
	Typ   string   `json:"typ,omitempty"`   // token type (Optional).

	// Email for which the application is requesting delegated access (Optional).
	Sub           string `json:"sub,omitempty"`
	UserFirstName string `json:"firstName,omitempty"`
	UserLastName  string `json:"lastName,omitempty"`
	DevUserID     int64  `json:"devOptiiUserId,omitempty"`
	TestUserID    int64  `json:"testOptiiUserId"`
	ProdUserId    int64  `json:"optiiUserId"`

	// The old name of Sub. Client keeps setting Prn to be
	// complaint with legacy OAuth 2.0 providers. (Optional)
	Prn string `json:"prn,omitempty"`

	// See http://tools.ietf.org/html/draft-jones-json-web-token-10#section-4.3
	// This array is marshalled using custom code (see (c *ClaimSet) encode()).
	PrivateClaims map[string]interface{} `json:"-"`
}
