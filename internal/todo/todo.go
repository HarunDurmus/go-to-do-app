package todo

type Todo struct {
	ID        string `json:"id"`
	Title     string `json:"task" bson:"task"`
	Task      string `json:"title" bson:"title"`
	Completed bool   `json:"completed" bson:"completed"`
}
