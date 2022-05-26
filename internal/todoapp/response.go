package todoapp

type Response struct {
	// error message of response that return from api
	// in: string
	Error string `json:"error"`

	Data *ToDoList `json:"data"`
}
