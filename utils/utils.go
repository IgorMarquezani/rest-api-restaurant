package utils

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var nonValidCookieValue []byte = []byte{'"', ',', ';', '\\'}

func IsInvalidCookieValue(char byte) bool {
  for i := 0; i < len(nonValidCookieValue); i++ {
    if char == nonValidCookieValue[i] {
      return true 
    }
  }

  return false
}

func HashString(str string) string {
	if str == "" {
		panic("String is empty")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(str), 8)
	if err != nil {
		panic(err)
	}

	return string(hashed)
}

func RandomCharacter() byte {
	var character byte

	rand.Seed(int64(time.Now().Nanosecond()))

	for character < 33 || character == '"' {
		character = byte(rand.Intn(127))
	}
	return character
}

func RandomByteArray() []byte {
	length := rand.Int31n(60)
	array := make([]byte, length)

	for i := 0; i < int(length); i++ {
    random := RandomCharacter()

    if IsInvalidCookieValue(random) {
      i--
      continue
    }

		array[i] = random
	}
	return array
}

func Invert(s string) string {
  var newStr []byte

	for i := len(s) - 1; i > -1; i-- {
		newStr = append(newStr, s[i])
	}
	return string(newStr)
}
