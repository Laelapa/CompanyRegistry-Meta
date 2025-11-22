package logging

const (

	// Service related fields --------------------------

	// FieldService is the service identifier
	FieldService = "service"
	// FieldEnv is the environment the service is running in / the logger is configured for
	FieldEnv = "env"
	// FieldLoggingLevel is the logging level the logger is configured for
	FieldLoggingLevel = "logging_level"

	// Request related fields --------------------------

	// FieldRemoteAddr is the IP address of the client making the request
	FieldRemoteAddr = "remote_addr"
	// FieldMethod is the HTTP method of the request (GET, POST, etc.)
	FieldMethod = "method"
	// FieldPath is the URL path of the request
	FieldPath = "path"
	// FieldReferer is the Referer header from the request
	FieldReferer = "referer"

	FieldUserID = "user_id"

	// HTTP Server related fields ----------------------

	FieldServerAddr = "server_addr"

	FieldServerPort = "server_port"

	// Kafka related fields ----------------------------

	FieldKafkaTopic   = "kafka_topic"
	FieldKafkaBrokers = "kafka_brokers"

	// Other common fields -----------------------------

	FieldError = "error"
)
