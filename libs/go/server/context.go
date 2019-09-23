package server

// Environment is a string that contains the current environment
// in which the application should boot.
type Environment string

const (
	// EnvProduction is the environment used in production environments
	EnvProduction Environment = "production"

	// EnvSandbox is the environment used in sandbox environments
	EnvSandbox Environment = "sandbox"

	// EnvDevelop is the environment used in development or staging environments
	EnvDevelop Environment = "develop"

	// EnvTest is the environment used in testing environment
	EnvTest Environment = "test"

	// EnvIntegration is the environment used in integration environments
	EnvIntegration Environment = "integration"
)

// Role is a string that contains the role in which the application
// should bootstrap.
type Role string

const (
	// RoleIndexer role will bootstrap the server in indexing mode.
	// This mode should be used only by endpoints that receive data from BigQ.
	RoleIndexer Role = "indexer"

	// RoleRead role will bootstrap the server in read mode.
	// This mode does not provide searching capabilities, but instead provides
	// reading assets by primary id.
	RoleRead Role = "read"

	// RoleSearch role will bootstrap the server in searching mode.
	// This mode setups the necessary endpoints so that clients can connect with us
	// and consume movements information.
	RoleSearch Role = "search"

	// RoleWrite role will bootstrap the server in write mode.
	// This mode enables the endpoints needed for writing things to a backing store.
	RoleWrite Role = "write"

	// RoleWorker role will bootstrap the server in worker mode.
	// This mode should be used only by endpoints that receive data from BigQ and do a specific task
	RoleWorker Role = "worker"

	// RoleMiddleEnd role will bootstrap the server in middleend mode.
	// This mode should be used by middle end apps. It allow read & write traffic, it jumps over mlauth validation
	RoleMiddleEnd Role = "middleend"
)

// ApplicationContext contains the necessary information for bootstraping the server
type ApplicationContext struct {
	Environment Environment `json:"environment"`
	Role        Role        `json:"role"`
	Tag         string      `json:"tag"`
}
