package types

type AuthType int 

const (
	OAuth = AuthType(iota)
	ApiKey 
)