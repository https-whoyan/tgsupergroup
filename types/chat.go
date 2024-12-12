package types

type ChatID = int64
type MessageID = int64
type ChatType string

func (ct ChatType) String() string {
	return string(ct)
}

const SuperGroupType ChatType = "supergroup"

func (ct ChatType) IsSuperGroup() bool {
	return ct == SuperGroupType
}

type Chat struct {
	ChatID   ChatID
	ChatType ChatType
}

func (c Chat) IsSuperGroup() bool {
	return c.ChatType.IsSuperGroup()
}
