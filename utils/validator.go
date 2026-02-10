// utils/validator.go

package utils

import (
    "regexp"
    "strings"
)

var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
    if email == "" || len(email) > 254 {
        return false
    }

    if !EmailRegex.MatchString(email) {
        return false
    }

    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }

    local, domain := parts[0], parts[1]

    if len(local) == 0 || len(local) > 64 {
        return false
    }

    if len(domain) == 0 {
        return false
    }

    if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
        return false
    }

    if strings.Contains(domain, "..") {
        return false
    }

    return true
}

func IsValidPassword(password string) bool {
    if len(password) < 8 {
        return false
    }

    hasUpper := false
    hasLower := false
    hasDigit := false
    hasSpecial := false

    for _, c := range password {
        switch {
        case c >= 'A' && c <= 'Z':
            hasUpper = true
        case c >= 'a' && c <= 'z':
            hasLower = true
        case c >= '0' && c <= '9':
            hasDigit = true
        default:
            hasSpecial = true
        }
    }

    return hasUpper && hasLower && hasDigit && hasSpecial
}

func IsValidUsername(username string) bool {
    if len(username) < 3 || len(username) > 50 {
        return false
    }

    validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
    for _, c := range username {
        if !strings.ContainsRune(validChars, c) {
            return false
        }
    }

    return true
}

type ValidationError struct {
	Msg string
}

func (e ValidationError) Error() string {
	return e.Msg
}