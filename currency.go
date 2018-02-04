// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package strfmt

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

func init() {
	d := Currency("")
	// register this format in the default registry
	Default.Add("currency", &d, IsCurrency)
}

var (
	// List of ISO currencies as of Feb 2018
	isoCurrencies = map[string]bool{
		"AFA": true,
		"EUR": true,
		"USD": true,
		"GBP": true,
		"CHF": true,
		"SGD": true,
		"CAD": true,
	}

	currencyMatcher = regexp.MustCompile(`^[A-Za-z]{3}$`)
)

// IsCurrency returns true if the provided string is a valid currency
func IsCurrency(str string) bool {
	if currencyMatcher.MatchString(str) {
		if auth, found := isoCurrencies[strings.ToUpper(str)]; found {
			// currency map may be configured to forbid some currencies
			return auth
		}
	}
	return false
}

// Currency represents an ISO-3 currency
//
// swagger:strfmt currency
type Currency string

// MarshalText turns this instance into text
func (d Currency) MarshalText() ([]byte, error) {
	return []byte(d), nil
}

// UnmarshalText hydrates this instance from text
func (d *Currency) UnmarshalText(data []byte) error { // validation is performed later on
	*d = Currency(data)
	return nil
}

// Scan reads a Currency value from database driver type.
func (d *Currency) Scan(raw interface{}) error {
	switch v := raw.(type) {
	case []byte:
		*d = Currency(v)
	case string:
		*d = Currency(v)
	case nil:
		*d = Currency("")
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.Currency from: %#v", v)
	}

	return nil
}

// Value converts Currency to a primitive value ready to be written to a database.
func (d Currency) Value() (driver.Value, error) {
	return driver.Value(string(d)), nil
}

// String converts this currency to a string
func (d Currency) String() string {
	return string(d)
}

// MarshalJSON returns the Currency as JSON
func (d Currency) MarshalJSON() ([]byte, error) {
	var w jwriter.Writer
	d.MarshalEasyJSON(&w)
	return w.BuildBytes()
}

// MarshalEasyJSON writes the Currency to a easyjson.Writer
func (d Currency) MarshalEasyJSON(w *jwriter.Writer) {
	w.String(string(d))
}

// UnmarshalJSON sets the Currency from JSON
func (d *Currency) UnmarshalJSON(data []byte) error {
	l := jlexer.Lexer{Data: data}
	d.UnmarshalEasyJSON(&l)
	return l.Error()
}

// UnmarshalEasyJSON sets the Currency from a easyjson.Lexer
func (d *Currency) UnmarshalEasyJSON(in *jlexer.Lexer) {
	if data := in.String(); in.Ok() {
		*d = Currency(data)
	}
}

// GetBSON returns the Currency a bson.M{} map.
func (d *Currency) GetBSON() (interface{}, error) {
	return bson.M{"data": string(*d)}, nil
}

// SetBSON sets the Currency from raw bson data
func (d *Currency) SetBSON(raw bson.Raw) error {
	var m bson.M
	if err := raw.Unmarshal(&m); err != nil {
		return err
	}

	if data, ok := m["data"].(string); ok {
		*d = Currency(data)
		return nil
	}

	return errors.New("couldn't unmarshal bson raw value as Currency")
}
