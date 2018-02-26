package gk

// DBOperation are the available operation we can do to a Datastore
type DBOperation string

const (
	// Select is the operation of reading a value(s) from a Datastore
	Select DBOperation = "SELECT"
	// Insert is the operation of sending a value to a Datastore
	Insert DBOperation = "INSERT"
	// Update is the operation of updating a value(s) from a Datastore
	Update DBOperation = "UPDATE"
	// Delete is the operation of deleting a value(s) from a Datastore
	Delete DBOperation = "DELETE"
)
