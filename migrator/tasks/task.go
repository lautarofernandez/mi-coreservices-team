package tasks

//Task is the interface for the Task
type Task interface {
	// Do is the function to do of the task
	 Do(data interface{}) error
}