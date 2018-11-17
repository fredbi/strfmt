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
	"strings"
	"testing"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testFormat string

func (t testFormat) MarshalText() ([]byte, error) {
	return []byte(string(t)), nil
}

func (t *testFormat) UnmarshalText(b []byte) error {
	*t = testFormat(string(b))
	return nil
}

func (t testFormat) String() string {
	return string(t)
}

func isTestFormat(s string) bool {
	return strings.HasPrefix(s, "tf")
}

type tf2 string

func (t tf2) MarshalText() ([]byte, error) {
	return []byte(string(t)), nil
}

func (t *tf2) UnmarshalText(b []byte) error {
	*t = tf2(string(b))
	return nil
}

func istf2(s string) bool {
	return strings.HasPrefix(s, "af")
}

func (t tf2) String() string {
	return string(t)
}

type bf string

func (t bf) MarshalText() ([]byte, error) {
	return []byte(string(t)), nil
}

func (t *bf) UnmarshalText(b []byte) error {
	*t = bf(string(b))
	return nil
}

func (t bf) String() string {
	return string(t)
}

func isbf(s string) bool {
	return strings.HasPrefix(s, "bf")
}

func istf3(s string) bool {
	return strings.HasPrefix(s, "ff")
}

func init() {
	tf := testFormat("")
	Default.Add("test-format", &tf, isTestFormat)
}

func TestFormatRegistry(t *testing.T) {
	f2 := tf2("")
	f3 := bf("")
	registry := NewFormats()

	assert.True(t, registry.ContainsName("test-format"))
	assert.True(t, registry.ContainsName("testformat"))
	assert.False(t, registry.ContainsName("ttt"))

	assert.True(t, registry.Validates("testformat", "tfa"))
	assert.False(t, registry.Validates("testformat", "ffa"))

	assert.True(t, registry.Add("tf2", &f2, istf2))
	assert.True(t, registry.ContainsName("tf2"))
	assert.False(t, registry.ContainsName("tfw"))
	assert.True(t, registry.Validates("tf2", "afa"))

	assert.False(t, registry.Add("tf2", &f3, isbf))
	assert.True(t, registry.ContainsName("tf2"))
	assert.False(t, registry.ContainsName("tfw"))
	assert.True(t, registry.Validates("tf2", "bfa"))
	assert.False(t, registry.Validates("tf2", "afa"))

	assert.False(t, registry.Add("tf2", &f2, istf2))
	assert.True(t, registry.Add("tf3", &f2, istf3))
	assert.True(t, registry.ContainsName("tf3"))
	assert.True(t, registry.ContainsName("tf2"))
	assert.False(t, registry.ContainsName("tfw"))
	assert.True(t, registry.Validates("tf3", "ffa"))

	assert.True(t, registry.DelByName("tf3"))
	assert.True(t, registry.Add("tf3", &f2, istf3))

	assert.True(t, registry.DelByName("tf3"))
	assert.False(t, registry.DelByName("unknown"))
	assert.False(t, registry.Validates("unknown", ""))
}

type testStruct struct {
	D          Date       `json:"d,omitempty"`
	DT         DateTime   `json:"dt,omitempty"`
	Dur        Duration   `json:"dur,omitempty"`
	URIField   URI        `json:"uri,omitempty" mapstructure:"uri"`
	Eml        Email      `json:"eml,omitempty"`
	UUIDField  UUID       `json:"uuid,omitempty" mapstructure:"uuid"`
	UUID3Field UUID3      `json:"uuid3,omitempty" mapstructure:"uuid3"`
	UUID4Field UUID4      `json:"uuid4field,omitempty"`
	UUID5Field UUID5      `json:"uuid5field,omitempty"`
	Hn         Hostname   `json:"hn,omitempty"`
	Ipv4       IPv4       `json:"ipv4,omitempty"`
	Ipv6       IPv6       `json:"ipv6,omitempty"`
	Cidr       CIDR       `json:"cidr,omitempty"`
	Mac        MAC        `json:"mac,omitempty"`
	Isbn       ISBN       `json:"isbn,omitempty"`
	Isbn10     ISBN10     `json:"isbn10,omitempty"`
	Isbn13     ISBN13     `json:"isbn13,omitempty"`
	Creditcard CreditCard `json:"creditcard,omitempty"`
	Ssn        SSN        `json:"ssn,omitempty"`
	Hexcolor   HexColor   `json:"hexcolor,omitempty"`
	Rgbcolor   RGBColor   `json:"rgbcolor,omitempty"`
	B64        Base64     `json:"b64,omitempty"`
	Pw         Password   `json:"pw,omitempty"`
	ULID       ULID       `json:"ulid,omitempty"`
}

func TestDecodeHook(t *testing.T) {
	registry := NewFormats()
	m := map[string]interface{}{
		"d":          "2014-12-15",
		"dt":         "2012-03-02T15:06:05.999999999Z",
		"dur":        "5s",
		"uri":        "http://www.dummy.com",
		"eml":        "dummy@dummy.com",
		"uuid":       "a8098c1a-f86e-11da-bd1a-00112444be1e",
		"uuid3":      "bcd02e22-68f0-3046-a512-327cca9def8f",
		"uuid4field": "025b0d74-00a2-4048-bf57-227c5111bb34",
		"uuid5field": "886313e1-3b8a-5372-9b90-0c9aee199e5d",
		"hn":         "somewhere.com",
		"ipv4":       "192.168.254.1",
		"ipv6":       "::1",
		"cidr":       "192.0.2.1/24",
		"mac":        "01:02:03:04:05:06",
		"isbn":       "0321751043",
		"isbn10":     "0321751043",
		"isbn13":     "978-0321751041",
		"hexcolor":   "#FFFFFF",
		"rgbcolor":   "rgb(255,255,255)",
		"pw":         "super secret stuff here",
		"ssn":        "111-11-1111",
		"creditcard": "4111-1111-1111-1111",
		"b64":        "ZWxpemFiZXRocG9zZXk=",
		"ulid":       "7ZZZZZZZZZZZZZZZZZZZZZZZZZ",
	}

	date, _ := time.Parse(RFC3339FullDate, "2014-12-15")
	dur, _ := ParseDuration("5s")
	dt, _ := ParseDateTime("2012-03-02T15:06:05.999999999Z")
	ulid, _ := ParseULID("7ZZZZZZZZZZZZZZZZZZZZZZZZZ")

	exp := &testStruct{
		D:          Date(date),
		DT:         dt,
		Dur:        Duration(dur),
		URIField:   URI("http://www.dummy.com"),
		Eml:        Email("dummy@dummy.com"),
		UUIDField:  UUID("a8098c1a-f86e-11da-bd1a-00112444be1e"),
		UUID3Field: UUID3("bcd02e22-68f0-3046-a512-327cca9def8f"),
		UUID4Field: UUID4("025b0d74-00a2-4048-bf57-227c5111bb34"),
		UUID5Field: UUID5("886313e1-3b8a-5372-9b90-0c9aee199e5d"),
		Hn:         Hostname("somewhere.com"),
		Ipv4:       IPv4("192.168.254.1"),
		Ipv6:       IPv6("::1"),
		Cidr:       CIDR("192.0.2.1/24"),
		Mac:        MAC("01:02:03:04:05:06"),
		Isbn:       ISBN("0321751043"),
		Isbn10:     ISBN10("0321751043"),
		Isbn13:     ISBN13("978-0321751041"),
		Creditcard: CreditCard("4111-1111-1111-1111"),
		Ssn:        SSN("111-11-1111"),
		Hexcolor:   HexColor("#FFFFFF"),
		Rgbcolor:   RGBColor("rgb(255,255,255)"),
		B64:        Base64("ZWxpemFiZXRocG9zZXk="),
		Pw:         Password("super secret stuff here"),
		ULID:       ulid,
	}

	test := new(testStruct)
	cfg := &mapstructure.DecoderConfig{
		DecodeHook: registry.MapStructureHookFunc(),
		// weakly typed will pass if this passes
		WeaklyTypedInput: false,
		Result:           test,
	}
	d, err := mapstructure.NewDecoder(cfg)
	require.NoError(t, err)
	err = d.Decode(m)
	require.NoError(t, err)
	assert.Equal(t, exp, test)
}

func TestDecodeDateTimeHook(t *testing.T) {
	testCases := []struct {
		Name  string
		Input string
	}{
		{
			"empty datetime",
			"",
		},
		{
			"invalid non empty datetime",
			"2019-01-01abc",
		},
	}
	registry := NewFormats()
	type layout struct {
		DateTime *DateTime `json:"datetime,omitempty"`
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			test := new(layout)
			cfg := &mapstructure.DecoderConfig{
				DecodeHook:       registry.MapStructureHookFunc(),
				WeaklyTypedInput: false,
				Result:           test,
			}
			d, err := mapstructure.NewDecoder(cfg)
			require.NoError(t, err)
			input := make(map[string]interface{})
			input["datetime"] = tc.Input
			err = d.Decode(input)
			require.Error(t, err, "error expected got none")
		})
	}
}

func TestDecode_ULID_Hook_Negative(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name  string
		Input string
	}{
		{
			"empty string for ulid",
			"",
		},
		{
			"invalid non empty ulid",
			"8000000000YYYYYYYYYYYYYYYY",
		},
	}
	registry := NewFormats()
	type layout struct {
		ULID *ULID `json:"ulid,omitempty"`
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			test := new(layout)
			cfg := &mapstructure.DecoderConfig{
				DecodeHook:       registry.MapStructureHookFunc(),
				WeaklyTypedInput: false,
				Result:           test,
			}
			d, err := mapstructure.NewDecoder(cfg)
			require.NoError(t, err)
			input := make(map[string]interface{})
			input["ulid"] = tc.Input
			err = d.Decode(input)
			require.Error(t, err, "error expected got none")
		})
	}
}

func Test_RegistryForGenerator(t *testing.T) {
	reg := Default
	_ = Default.DelByName("test-format")

	genReg, ok := reg.(RegistryForGenerator)
	if !assert.True(t, ok) {
		t.FailNow()
	}

	assert.Equal(t, "datetime", genReg.Normalized("date-time"))
	assert.Equal(t, "", genReg.Normalized("unknown"))

	assert.ElementsMatch(t, genReg.Formats(), []string{
		"bsonobjectid",
		"byte",
		"creditcard",
		"date",
		"datetime",
		"duration",
		"email",
		"hexcolor",
		"hostname",
		"ipv4",
		"ipv6",
		"isbn",
		"isbn10",
		"isbn13",
		"mac",
		"password",
		"rgbcolor",
		"ssn",
		"uri",
		"uuid",
		"uuid3",
		"uuid4",
		"uuid5",
	})

	assert.ElementsMatch(t, genReg.Types(), []string{
		"Base64",
		"CreditCard",
		"Date",
		"DateTime",
		"Duration",
		"Email",
		"HexColor",
		"Hostname",
		"IPv4",
		"IPv6",
		"ISBN",
		"ISBN10",
		"ISBN13",
		"MAC",
		"ObjectId",
		"Password",
		"RGBColor",
		"SSN",
		"URI",
		"UUID",
		"UUID3",
		"UUID4",
		"UUID5",
	})

	assert.Equal(t, map[string]string{
		"bsonobjectid": "ObjectId",
		"byte":         "Base64",
		"creditcard":   "CreditCard",
		"date":         "Date",
		"datetime":     "DateTime",
		"duration":     "Duration",
		"email":        "Email",
		"hexcolor":     "HexColor",
		"hostname":     "Hostname",
		"ipv4":         "IPv4",
		"ipv6":         "IPv6",
		"isbn":         "ISBN",
		"isbn10":       "ISBN10",
		"isbn13":       "ISBN13",
		"mac":          "MAC",
		"password":     "Password",
		"rgbcolor":     "RGBColor",
		"ssn":          "SSN",
		"uri":          "URI",
		"uuid":         "UUID",
		"uuid3":        "UUID3",
		"uuid4":        "UUID4",
		"uuid5":        "UUID5",
	}, genReg.FormatToTypes())

	assert.Equal(t, map[string][]string{
		"ObjectId":   []string{"bsonobjectid"},
		"ISBN13":     []string{"isbn13"},
		"SSN":        []string{"ssn"},
		"Date":       []string{"date"},
		"Email":      []string{"email"},
		"IPv4":       []string{"ipv4"},
		"ISBN":       []string{"isbn"},
		"CreditCard": []string{"creditcard"},
		"Password":   []string{"password"},
		"ISBN10":     []string{"isbn10"},
		"URI":        []string{"uri"},
		"Hostname":   []string{"hostname"},
		"IPv6":       []string{"ipv6"},
		"MAC":        []string{"mac"},
		"UUID":       []string{"uuid"},
		"UUID3":      []string{"uuid3"},
		"UUID5":      []string{"uuid5"},
		"RGBColor":   []string{"rgbcolor"},
		"Duration":   []string{"duration"},
		"DateTime":   []string{"datetime"},
		"UUID4":      []string{"uuid4"},
		"HexColor":   []string{"hexcolor"},
		"Base64":     []string{"byte"},
	}, genReg.TypeToFormats())

	assert.Equal(t, map[string]string{
		"Date":       "Date{}",
		"Hostname":   "Hostname(\"\")",
		"RGBColor":   "RGBColor(\"rgb(0,0,0)\")",
		"UUID3":      "UUID3(\"\")",
		"UUID4":      "UUID4(\"\")",
		"HexColor":   "HexColor(\"#000000\")",
		"IPv4":       "IPv4(\"\")",
		"IPv6":       "IPv6(\"\")",
		"ISBN13":     "ISBN13(\"\")",
		"MAC":        "MAC(\"\")",
		"CreditCard": "CreditCard(\"\")",
		"ISBN10":     "ISBN10(\"\")",
		"Password":   "Password(\"\")",
		"UUID5":      "UUID5(\"\")",
		"ObjectId":   "ObjectId(\"\")",
		"SSN":        "SSN(\"\")",
		"URI":        "URI(\"\")",
		"Base64":     "Base64([]byte(nil))",
		"DateTime":   "DateTime{}",
		"Duration":   "Duration(0)",
		"Email":      "Email(\"\")",
		"ISBN":       "ISBN(\"\")",
		"UUID":       "UUID(\"\")",
	}, genReg.ZeroExpressions())
}
