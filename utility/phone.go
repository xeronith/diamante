package utility

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

var _RegExPhone *regexp.Regexp

func init() {
	exp, err := regexp.Compile(`^\d*$`)
	if err != nil {
		log.Fatal(err.Error())
	}
	_RegExPhone = exp
}

func SanitizePhone(phoneNumber string) string {
	phoneNumber = strings.TrimLeft(phoneNumber, " +0")
	if !_RegExPhone.MatchString(phoneNumber) {
		return ""
	}

	phone, err := phonenumbers.Parse(phoneNumber, "IR")
	if err != nil {
		return phoneNumber
	}

	return fmt.Sprintf("%d%d", *phone.CountryCode, *phone.NationalNumber)
}
