//go:build mongo

package strfmt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestBSONDate(t *testing.T) {
	dateOriginal := Date(time.Date(2014, 10, 10, 0, 0, 0, 0, time.UTC))

	bsonData, err := bson.Marshal(&dateOriginal)
	require.NoError(t, err)

	var dateCopy Date
	err = bson.Unmarshal(bsonData, &dateCopy)
	require.NoError(t, err)
	assert.Equal(t, dateOriginal, dateCopy)
}

func TestBSONBase64(t *testing.T) {
	const b64 string = "This is a byte array with unprintable chars, but it also isn"
	b := []byte(b64)
	subj := Base64(b)

	bsonData, err := bson.Marshal(subj)
	require.NoError(t, err)

	var b64Copy Base64
	err = bson.Unmarshal(bsonData, &b64Copy)
	require.NoError(t, err)
	assert.Equal(t, subj, b64Copy)
}

func TestBSONDuration(t *testing.T) {
	dur := Duration(42)
	bsonData, err := bson.Marshal(&dur)
	require.NoError(t, err)

	var durCopy Duration
	err = bson.Unmarshal(bsonData, &durCopy)
	require.NoError(t, err)
	assert.Equal(t, dur, durCopy)
}

func TestBSONDateTime(t *testing.T) {
	for caseNum, example := range testCases {
		t.Logf("Case #%d", caseNum)
		dt := DateTime(example.time)

		bsonData, err := bson.Marshal(&dt)
		require.NoError(t, err)

		var dtCopy DateTime
		err = bson.Unmarshal(bsonData, &dtCopy)
		require.NoError(t, err)
		// BSON DateTime type loses timezone information, so compare UTC()
		assert.Equal(t, time.Time(dt).UTC(), time.Time(dtCopy).UTC())

		// Check value marshaling explicitly
		m := bson.M{"data": dt}
		bsonData, err = bson.Marshal(&m)
		require.NoError(t, err)

		var mCopy bson.M
		err = bson.Unmarshal(bsonData, &mCopy)
		require.NoError(t, err)

		data, ok := m["data"].(DateTime)
		assert.True(t, ok)
		assert.Equal(t, time.Time(dt).UTC(), time.Time(data).UTC())
	}
}

func TestBSONULID(t *testing.T) {
	t.Parallel()
	t.Run("positive", func(t *testing.T) {
		t.Parallel()
		ulid, _ := ParseULID(testUlid)

		bsonData, err := bson.Marshal(&ulid)
		require.NoError(t, err)

		var ulidUnmarshaled ULID
		err = bson.Unmarshal(bsonData, &ulidUnmarshaled)
		require.NoError(t, err)
		assert.Equal(t, ulid, ulidUnmarshaled)

		// Check value marshaling explicitly
		m := bson.M{"data": ulid}
		bsonData, err = bson.Marshal(&m)
		require.NoError(t, err)

		var mUnmarshaled bson.M
		err = bson.Unmarshal(bsonData, &mUnmarshaled)
		require.NoError(t, err)

		data, ok := m["data"].(ULID)
		assert.True(t, ok)
		assert.Equal(t, ulid, data)
	})
	t.Run("negative", func(t *testing.T) {
		t.Parallel()
		uuid := UUID("00000000-0000-0000-0000-000000000000")
		bsonData, err := bson.Marshal(&uuid)
		require.NoError(t, err)

		var ulidUnmarshaled ULID
		err = bson.Unmarshal(bsonData, &ulidUnmarshaled)
		require.Error(t, err)
	})
}
