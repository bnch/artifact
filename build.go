package main

const (
	statusBuilding = iota
	statusBuilt
	statusFailed
)

// Build is a collection of artifacts from a build.
type Build struct {
	ID           int
	RepositoryID int
	Commit       string
	Status       int
}
