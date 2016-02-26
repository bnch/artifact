package main

// Repository is a github repository connected with artifact.
type Repository struct {
	ID        int
	User      string
	Repo      string
	SecretKey string
}
