package types

type Role int

const (
	UserRole = Role(iota)
	AssistantRole
	ToolRole
	AlarmRole
	MemoryRole
)