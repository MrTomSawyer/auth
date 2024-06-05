package rbac

type Action string

const (
	Read   Action = "read"
	Insert Action = "insert"
	Update Action = "update"
	Delete Action = "delete"
	Audit  Action = "audit"
)

type Permission struct {
	Name        string
	Description string
	Entity      string
	Action      Action
}
