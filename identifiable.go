package config

import "fmt"

type Identifiable interface {
	GetID() string
}

func FindByID[T Identifiable](slice []T, id string) (*T, error) {
	for i := range slice {
		if slice[i].GetID() == id {
			return &slice[i], nil
		}
	}
	return nil, fmt.Errorf("element with ID '%s' not found", id)
}