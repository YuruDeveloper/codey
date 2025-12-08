package auth

type Auth interface {
	Update()
	Key() string
}