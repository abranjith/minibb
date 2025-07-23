package utils

import (
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"strings"
)

// GenerateTripcode generates a 4chan-style newstyle tripcode
// Format: Name##Password -> Name !TripCode
func GenerateTripcode(name string) string {
	parts := strings.Split(name, "##")
	if len(parts) != 2 {
		return name // No tripcode password provided
	}

	username := parts[0]
	password := parts[1]

	if password == "" {
		return username
	}

	// Generate tripcode using 4chan newstyle algorithm
	tripcode := generateNewstyleTripcode(password)

	return username + " !" + tripcode
}

func generateNewstyleTripcode(password string) string {
	// Hash the password with MD5
	hash := md5.Sum([]byte(password + "H."))

	// Use first 8 bytes as DES key
	key := hash[:8]

	// Create DES cipher
	cipher, err := des.NewCipher(key)
	if err != nil {
		// Fallback to simple hash if DES fails
		return base64.StdEncoding.EncodeToString(hash[:])[:10]
	}

	// Encrypt the password
	plaintext := []byte(password + "........") // Pad to 8 bytes
	if len(plaintext) > 8 {
		plaintext = plaintext[:8]
	} else {
		for len(plaintext) < 8 {
			plaintext = append(plaintext, '.')
		}
	}

	ciphertext := make([]byte, 8)
	cipher.Encrypt(ciphertext, plaintext)

	// Encode and take first 10 characters
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	if len(encoded) > 10 {
		encoded = encoded[:10]
	}

	// Replace characters that might cause issues
	encoded = strings.ReplaceAll(encoded, "+", ".")
	encoded = strings.ReplaceAll(encoded, "/", ".")
	encoded = strings.ReplaceAll(encoded, "=", ".")

	return encoded
}
