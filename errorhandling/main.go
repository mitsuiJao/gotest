package main

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound   = errors.New("resource not found")
	ErrPermission = errors.New("permission denied")
)

type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %q: %s", e.Field, e.Msg)
}

func fetchUser(id int) (string, error) {
	if id < 0 {
		return "", &ValidationError{Field: "id", Msg: "must be non-negative"}
	}
	if id == 0 {
		return "", fmt.Errorf("fetchUser: %w", ErrNotFound)
	}
	if id == 999 {
		return "", fmt.Errorf("fetchUser: %w", ErrPermission)
	}
	return fmt.Sprintf("user-%d", id), nil
}

func validateInput(id int, name string) error {
	var errs []error
	if id < 0 {
		errs = append(errs, &ValidationError{Field: "id", Msg: "must be non-negative"})
	}
	if name == "" {
		errs = append(errs, &ValidationError{Field: "name", Msg: "must not be empty"})
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func main() {
	testIDs := []int{-1, 0, 999, 42}

	for _, id := range testIDs {
		_, err := fetchUser(id)
		if err == nil {
			fmt.Printf("id=%d: success\n", id)
			continue
		}

		switch {
		case errors.Is(err, ErrNotFound):
			fmt.Printf("id=%d: not found -> %v\n", id, err)
		case errors.Is(err, ErrPermission):
			fmt.Printf("id=%d: permission denied -> %v\n", id, err)
		default:
			var vErr *ValidationError
			if errors.As(err, &vErr) {
				fmt.Printf("id=%d: validation error on field %q -> %v\n", id, vErr.Field, err)
			} else {
				fmt.Printf("id=%d: unknown error -> %v\n", id, err)
			}
		}
	}

	fmt.Println("---")

	if err := validateInput(-1, ""); err != nil {
		fmt.Println("validateInput errors:")
		fmt.Println(err)
	}
}
