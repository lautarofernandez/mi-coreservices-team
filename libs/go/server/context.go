package server

// Environment is a string that contains the current environment
// in which the application should boot.
type Environment string

const (
	// EnvProduction is the environment used in production environments
	EnvProduction Environment = "production"

	// EnvDevelop is the environment used in development or staging environments
	EnvDevelop = "develop"

	// EnvTest is the environment used in testing environment
	EnvTest = "test"

	// EnvIntegration is the environment used in integration environments
	EnvIntegration = "integration"
)

// Role is a string that contains the role in which the application
// should bootstrap.
type Role string

const (
	// RoleIndexer role will bootstrap the server in indexing mode.
	// This mode should be used only by endpoints that receive data from BigQ and index a entity
	RoleIndexer Role = "indexer"

	// RoleRead role will bootstrap the server in read mode.
	// This mode does not provide searching capabilities, but instead provides
	// reading assets by primary id.
	RoleRead = "read"

	// RoleSearch role will bootstrap the server in searching mode.
	// This mode setups the necessary endpoints so that clients can connect with us
	// and consume movements information.
	RoleSearch = "search"

	// RoleWrite role will bootstrap the server in write mode.
	// This mode enables the endpoints needed for writing things to a backing store.
	RoleWrite = "write"

	// RoleWorker role will bootstrap the server in worker mode.
	// This mode should be used only by endpoints that receive data from BigQ and do a specific task
	RoleWorker = "worker"
)

// ApplicationContext contains the necessary information for bootstraping the server
type ApplicationContext struct {
	Environment Environment `json:"environment"`
	Role        Role        `json:"role"`
	Tag         string      `json:"tag"`
}
