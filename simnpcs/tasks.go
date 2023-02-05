package simnpcs

// TaskType represents the type of task.
type TaskType int

// Task types.
const (
	TaskNone TaskType = iota
	TaskFind
	TaskBuy
	TaskSell
	TaskCraft
)

// String returns the string representation of the task type.
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

// Task represents a task.
type Task struct {
	ID   uint64   // Unique id
	Type TaskType // Type of task
	Item *Item    // Item to find/buy/sell/craft
}

// String returns the string representation of the task.
func (t *Task) String() string {
	return t.Type.String() + " " + t.Item.Name
}

// Tasks represents a list of tasks.
type Tasks []*Task

// AddTask adds a task to the list.
func (c *Tasks) AddTask(id uint64, item *Item, t TaskType) {
	*c = append(*c, &Task{
		ID:   id,
		Type: t,
		Item: item,
	})
}

// CompleteTask completes a task.
func (c *Tasks) CompleteTask(t *Task) {
	c.RemoveTask(c.FindTask(t))
}

// RemoveTask removes a task from the list.
func (c *Tasks) RemoveTask(s int) {
	*c = append((*c)[:s], (*c)[s+1:]...)
}

// FindTask finds a task in the list.
func (c Tasks) FindTask(t *Task) int {
	for i, tt := range c {
		if tt == t {
			return i
		}
	}
	return -1
}
