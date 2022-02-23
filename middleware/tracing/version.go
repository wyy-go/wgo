package tracing

func Version() string {
	return "0.25.0"
}

func SemVersion() string {
	return "semver:" + Version()
}
