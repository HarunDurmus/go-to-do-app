package todo

type Todo struct {
	ID        string `json:"id"`
	Task      string `json:"task" bson:"task"`
	Completed bool   `json:"completed" bson:"copleted"`
}
