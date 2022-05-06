package services

import gonanoid "github.com/matoous/go-nanoid"

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateSlug() (string, error) {
	return gonanoid.Generate(alphabet, 8)
}
