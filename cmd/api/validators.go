package main

import (
	"fmt"
	"strings"
)

// var (
// 	ErrTitleEmpty      = errors.New("title required")
// 	ErrTitleTooShort   = errors.New("title too short")
// 	ErrTitleTooLong    = errors.New("title too long")
// 	ErrContentEmpty    = errors.New("content required")
// 	ErrContentTooLong  = errors.New("content too long")
// 	ErrContentTooShort = errors.New("content too short")
// )

func validateRequired(v string, f string) error {
	l := len(v)

	if l == 0 {
		return fmt.Errorf("%s is required", f)
	}
	return nil
}

func validateMin(v string, min int, f string) error {
	l := len(v)
	if l < min {
		return fmt.Errorf("%s must be at least %d", f, min)
	}
	return nil
}

func validateMax(v string, max int, f string) error {
	l := len(v)
	if l > max {
		return fmt.Errorf("%s maximum %d", f, max)
	}
	return nil
}

func validateEmail(email string, f string) error {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fmt.Errorf("%s must be a valid email", f)
	}
	return nil
}

// func validateTitle(title string) error {
// 	l := len(title)

// 	switch {
// 	case l == 0:
// 		return ErrTitleEmpty
// 	case l < minTitle:
// 		return ErrTitleTooShort
// 	case l > maxTitle:
// 		return ErrTitleTooLong
// 	default:
// 		return nil
// 	}
// }

// func validateContent(content string) error {
// 	l := len(content)

// 	switch {
// 	case l == 0:
// 		return ErrContentEmpty
// 	case l < minTitle:
// 		return ErrContentTooShort
// 	case l > maxTitle:
// 		return ErrContentTooLong
// 	default:
// 		return nil
// 	}
// }
