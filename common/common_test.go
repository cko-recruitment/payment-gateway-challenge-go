package common_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cko-recruitment/payment-gateway-challenge-go/common"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	assert := assert.New(t)
	for i := 0; i < 10; i++ {
		id, err := common.NextID()
		assert.Nil(err)
		sameId, err := common.ParseID(id.String())
		var tail string
		if err != nil {
			tail = ": " + err.Error()
		}
		assert.Nilf(err, "when reparsing stringified %v%s", id, tail)
		assert.Equal(sameId, id)
	}

	_, err := common.ParseID("Hello, world!")
	assert.NotNil(err)
}

func checkReader(assert *assert.Assertions, expected string, actual io.Reader, fmtAndArgs ...any) {
	content, err := common.ReaderContent(actual)
	if closer, ok := actual.(io.Closer); ok {
		defer closer.Close()
	}
	if len(fmtAndArgs) == 0 {
		assert.Nil(err)
		assert.Equal(expected, content)
	} else {
		fmtString, args := fmtAndArgs[0].(string), fmtAndArgs[1:]
		assert.Nilf(err, fmtString, args...)
		assert.Equalf(expected, content, fmtString, args...)
	}
}


func TestRefuse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	assert := assert.New(t)
	for _, reason := range [2]string{
		"A very bad request",
		"This request broke our server",
	} {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		common.Refuse(c, fmt.Errorf("%s", reason))
		res := rec.Result()
		assert.Equalf(http.StatusBadRequest, res.StatusCode, "when %s", reason)
		checkReader(assert, `{"error":"` + reason + `"}`, res.Body, "when %s", reason)
	}
}

type NoTags struct {
	X int
}

type testPair struct {
	src any
	expected string
}

func TestEncodeJson(t *testing.T) {
	assert := assert.New(t)
	for _, tp := range []testPair{
		testPair{map[string]any{"a": 42, "b": map[string]int{"e": 271, "p": 314}}, `{"a":42,"b":{"e":271,"p":314}}` + "\n"},
		testPair{42, "42\n"},
		testPair{[]float64{2.71, 3.1416}, "[2.71,3.1416]\n"},	// a tricky example, but should pass
		testPair{NoTags{42}, "{\"X\":42}\n"},
	} {
		rdr, err := common.EncodeJSON(tp.src)
		assert.Nilf(err, "when marshaling %v", tp.src)
		checkReader(assert, tp.expected, rdr)
	}
}
