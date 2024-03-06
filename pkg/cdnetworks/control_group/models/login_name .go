package models

// LoginName represents a login name.
type LoginName struct {
	LoginName *string `json:"loginName,omitempty" xml:"loginName,omitempty"`
}

// SetLoginName sets the login name.
func (s *LoginName) SetLoginName(v string) *LoginName {
	s.LoginName = &v
	return s
}
