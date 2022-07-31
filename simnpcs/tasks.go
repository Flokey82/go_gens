package simnpcs

type TaskType int

const (
	TaskNone TaskType = iota
	TaskFind
	TaskBuy
	TaskSell
	TaskCraft
)

func (t *TaskType) String() string {
	switch *t {
	case TaskFind:
		return "find"
	case TaskBuy:
		return "buy"
	case TaskSell:
		return "sell"
	case TaskCraft:
		return "craft"
	default:
		return "Unknown"
	}
}

type Task struct {
	ID   uint64
	Type TaskType
	Item *Item
}

func (t *Task) String() string {
	return t.Type.String() + " " + t.Item.Name
}

type Tasks []*Task

func (c *Tasks) AddTask(id uint64, item *Item, t TaskType) {
	*c = append(*c, &Task{
		ID:   id,
		Type: t,
		Item: item,
	})
}

func (c *Tasks) CompleteTask(t *Task) {
	c.RemoveTask(c.FindTask(t))
}

func (c *Tasks) RemoveTask(s int) {
	*c = append((*c)[:s], (*c)[s+1:]...)
}

func (c Tasks) FindTask(t *Task) int {
	for i, tt := range c {
		if tt == t {
			return i
		}
	}
	return -1
}
