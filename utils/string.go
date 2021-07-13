package utils

import "net/mail"

func ValidateEmail(s string) error {
	_, err :=  mail.ParseAddress(s)
	return err
}
