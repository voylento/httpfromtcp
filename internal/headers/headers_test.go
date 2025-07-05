package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParsing(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("  Host:    localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 28, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	// Test: 1st of 2 headers
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nUser-Agent: curl/7.81.0\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: 2nd of 2 headers
	n2, done, err := headers.Parse(data[n:])
	require.NoError(t, err)
	assert.Equal(t, 2, len(headers))
	require.Equal(t, "curl/7.81.0", headers["user-agent"])
	assert.Equal(t, 25, n2)
	assert.False(t, done)
	
	// Test: done
	n, done, err = headers.Parse(data[n+n2:])
	require.NoError(t, err)
	assert.True(t, done)
	

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("   Host : localhost:42069   \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid character in header name
	headers = NewHeaders()
	data = []byte("H@st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid multiple values for same header
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nSet-Person: testity1\r\nSet-Person: testity2\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	n2, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	assert.False(t, done)
	n3, done, err := headers.Parse(data[n+n2:])
	require.NoError(t, err)
	assert.False(t, done)
	_, done, err = headers.Parse(data[n+n2+n3:])
	require.NoError(t, err)
	assert.True(t, done)
	assert.Equal(t, "testity1, testity2", headers["set-person"])

	// Test: headers with empty values
	headers = NewHeaders()
	data = []byte("Set-Person:\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, "", headers["set-person"])
}

