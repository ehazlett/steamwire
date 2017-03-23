package version

var (
	name        = "steamwire"
	version     = "0.1.0"
	description = "steam application news updates"
	gitCommit   = "HEAD"
)

// Name returns the name of the application
func Name() string {
	return name
}

// Version returns the version including GitCommit
func Version() string {
	return version + " (" + gitCommit + ")"
}

// Description returns the application description
func Description() string {
	return description
}

// FullVersion returns the application name and version
func FullVersion() string {
	return Name() + " " + Version()
}
