package dbSeed

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomSchool generates a random school name
func RandomSchool() string {
	return "I.P.L.T. " + RandomString(6) + " " + RandomString(6)
}

// RandomEmail returns a random email for testing purposes
func RandomEmail() string {
	return fmt.Sprintf("%s@utm.com", RandomString(5))
}

// RandomGender returns either M or F
func RandomGender() string {
	m := rand.Intn(2)
	if m == 0 {
		return "M"
	} else {
		return "F"
	}

}

// RandomGender returns either M or F
func RandomPhoneNumber() string {
	phoneNum := "+373"
	// ascii numbs between 48 and 57
	for i := 0; i < 9; i++ {
		phoneNum = phoneNum + strconv.FormatInt(RandomInt(0, 9), 10)
	}
	return phoneNum
}

// RandomResidence returns a random street
func RandomResidence() string {
	return fmt.Sprintf("Moldova, Chisinau, str. %s %s", RandomString(5), RandomString(7))
}

// RandomBirthDate returns a random date
func RandomBirthDate() time.Time {
	min := time.Date(1990, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2013, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}
