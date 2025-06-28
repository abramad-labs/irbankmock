package version

// Value of this variable automatically gets replaced with the correct version of the
// binary in build-time using -ldflags flag in go build
var ServerVersion string = "dev"
