//go:build mongo

package strfmt

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

func (d Date) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": d.String()})
}

func (d *Date) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if data, ok := m["data"].(string); ok {
		rd, err := time.ParseInLocation(RFC3339FullDate, data, DefaultTimeLocation)
		if err != nil {
			return err
		}
		*d = Date(rd)
		return nil
	}

	return errors.New("couldn't unmarshal bson bytes value as Date")
}

// MarshalBSON document from this value
func (b Base64) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": b.String()})
}

// UnmarshalBSON document into this value
func (b *Base64) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if bd, ok := m["data"].(string); ok {
		vb, err := base64.StdEncoding.DecodeString(bd)
		if err != nil {
			return err
		}
		*b = Base64(vb)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as base64")
}

func (d Duration) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": d.String()})
}

func (d *Duration) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if data, ok := m["data"].(string); ok {
		rd, err := ParseDuration(data)
		if err != nil {
			return err
		}
		*d = Duration(rd)
		return nil
	}

	return errors.New("couldn't unmarshal bson bytes value as Date")
}

// MarshalBSON renders the DateTime as a BSON document
func (t DateTime) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": t})
}

// UnmarshalBSON reads the DateTime from a BSON document
func (t *DateTime) UnmarshalBSON(data []byte) error {
	var obj struct {
		Data DateTime
	}

	if err := bson.Unmarshal(data, &obj); err != nil {
		return err
	}

	*t = obj.Data

	return nil
}

// MarshalBSONValue is an interface implemented by types that can marshal themselves
// into a BSON document represented as bytes. The bytes returned must be a valid
// BSON document if the error is nil.
//
// Marshals a DateTime as a bson.TypeDateTime, an int64 representing
// milliseconds since epoch.
func (t DateTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	// UnixNano cannot be used directly, the result of calling UnixNano on the zero
	// Time is undefined. Thats why we use time.Nanosecond() instead.

	tNorm := NormalizeTimeForMarshal(time.Time(t))
	i64 := tNorm.Unix()*1000 + int64(tNorm.Nanosecond())/1e6

	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(i64))

	return bson.TypeDateTime, buf, nil
}

// UnmarshalBSONValue is an interface implemented by types that can unmarshal a
// BSON value representation of themselves. The BSON bytes and type can be
// assumed to be valid. UnmarshalBSONValue must copy the BSON value bytes if it
// wishes to retain the data after returning.
func (t *DateTime) UnmarshalBSONValue(tpe bsontype.Type, data []byte) error {
	if tpe == bson.TypeNull {
		*t = DateTime{}
		return nil
	}

	if len(data) != 8 {
		return errors.New("bson date field length not exactly 8 bytes")
	}

	i64 := int64(binary.LittleEndian.Uint64(data))
	// TODO: Use bsonprim.DateTime.Time() method
	*t = DateTime(time.Unix(i64/1000, i64%1000*1000000))

	return nil
}

// MarshalBSON document from this value
func (u ULID) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": u.String()})
}

// UnmarshalBSON document into this value
func (u *ULID) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		id, err := ulid.ParseStrict(ud)
		if err != nil {
			return fmt.Errorf("couldn't parse bson bytes as ULID: %w", err)
		}
		u.ULID = id
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as ULID")
}

// MarshalBSON document from this value
func (u URI) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": u.String()})
}

// UnmarshalBSON document into this value
func (u *URI) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*u = URI(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as uri")
}

// MarshalBSON document from this value
func (e Email) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": e.String()})
}

// UnmarshalBSON document into this value
func (e *Email) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*e = Email(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as email")
}

// MarshalBSON document from this value
func (h Hostname) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": h.String()})
}

// UnmarshalBSON document into this value
func (h *Hostname) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*h = Hostname(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as hostname")
}

// MarshalBSON document from this value
func (u IPv4) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": u.String()})
}

// UnmarshalBSON document into this value
func (u *IPv4) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*u = IPv4(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as ipv4")
}

// MarshalBSON document from this value
func (u IPv6) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": u.String()})
}

// UnmarshalBSON document into this value
func (u *IPv6) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*u = IPv6(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as ipv6")
}

// MarshalBSON document from this value
func (u CIDR) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": u.String()})
}

// UnmarshalBSON document into this value
func (u *CIDR) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*u = CIDR(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as CIDR")
}

// MarshalBSON document from this value
func (u MAC) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": u.String()})
}

// UnmarshalBSON document into this value
func (u *MAC) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*u = MAC(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as MAC")
}

// MarshalBSON document from this value
func (r Password) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": r.String()})
}

// UnmarshalBSON document into this value
func (r *Password) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*r = Password(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as Password")
}

// MarshalBSON document from this value
func (u UUID) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": u.String()})
}

// UnmarshalBSON document into this value
func (u *UUID) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*u = UUID(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as UUID")
}

// MarshalBSON document from this value
func (u UUID3) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": u.String()})
}

// UnmarshalBSON document into this value
func (u *UUID3) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*u = UUID3(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as UUID3")
}

// MarshalBSON document from this value
func (u UUID4) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": u.String()})
}

// UnmarshalBSON document into this value
func (u *UUID4) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*u = UUID4(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as UUID4")
}

// MarshalBSON document from this value
func (u UUID5) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": u.String()})
}

// UnmarshalBSON document into this value
func (u *UUID5) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*u = UUID5(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as UUID5")
}

// MarshalBSON document from this value
func (u ISBN) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": u.String()})
}

// UnmarshalBSON document into this value
func (u *ISBN) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*u = ISBN(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as ISBN")
}

// MarshalBSON document from this value
func (u ISBN10) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": u.String()})
}

// UnmarshalBSON document into this value
func (u *ISBN10) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*u = ISBN10(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as ISBN10")
}

// MarshalBSON document from this value
func (u ISBN13) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": u.String()})
}

// UnmarshalBSON document into this value
func (u *ISBN13) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*u = ISBN13(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as ISBN13")
}

// MarshalBSON document from this value
func (u CreditCard) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": u.String()})
}

// UnmarshalBSON document into this value
func (u *CreditCard) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*u = CreditCard(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as CreditCard")
}

// MarshalBSON document from this value
func (u SSN) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": u.String()})
}

// UnmarshalBSON document into this value
func (u *SSN) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*u = SSN(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as SSN")
}

// MarshalBSON document from this value
func (h HexColor) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": h.String()})
}

// UnmarshalBSON document into this value
func (h *HexColor) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*h = HexColor(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as HexColor")
}

// MarshalBSON document from this value
func (r RGBColor) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.M{"data": r.String()})
}

// UnmarshalBSON document into this value
func (r *RGBColor) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	if ud, ok := m["data"].(string); ok {
		*r = RGBColor(ud)
		return nil
	}
	return errors.New("couldn't unmarshal bson bytes as RGBColor")
}
