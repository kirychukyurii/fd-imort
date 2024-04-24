package config

// S3 represents a configuration for accessing an S3 bucket.
//
// AccessKeyID is the access key for authenticating with AWS.
//
// SecretAccessKey is the secret access key for authenticating with AWS.
//
// Region is the AWS region where the S3 bucket is located.
//
// Bucket is the name of the S3 bucket.
type S3 struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
}

type Server struct {
	Address string
	Token   string
}

// Config represents a configuration object that is used to store application settings.
//
// LogLevel is the log level for the application.
//
// ExportedPath is the base path where exported files are stored.
//
// Domain is the domain name for the application.
//
// DSN is the database connection string.
//
// S3 is the configuration for accessing an S3 bucket. Refer to the documentation of the S3 type for more details.
type Config struct {
	LogLevel      string
	LogFile       string
	Workers       int
	ExportedPath  string
	AttachmentDir string
	Domain        string
	DSN           string
	S3            *S3
	Server        *Server
}

func New() *Config {
	return &Config{
		S3:     &S3{},
		Server: &Server{},
	}
}
