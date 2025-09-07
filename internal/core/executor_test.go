package core

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thaison199py/multi-threaded-redis/internal/constant"
	"github.com/thaison199py/multi-threaded-redis/internal/data_structure"
)

// DecodeInt64 helper function to decode response into int64
func DecodeInt64(resp []byte, val *int64) error {
	str := string(resp)
	if !strings.HasPrefix(str, ":") {
		return errors.New("not an integer response")
	}
	str = strings.TrimRight(str[1:], "\r\n")
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*val = num
	return nil
}

func setupDictStore() *data_structure.Dict {
	d := data_structure.CreateDict()
	return d
}

func TestCmdExists(t *testing.T) {
	d := setupDictStore()
	d.Set("foo", d.NewObj("foo", "bar", -1))
	// Create expired key by setting expiry to current time - 1 second
	d.Set("baz", d.NewObj("baz", "qux", -1))
	d.SetExpiry("baz", -1000) // Set expiry to 1 second in the past
	dictStore = d

	// Test: 1 key exists and not expired
	res := cmdEXISTS([]string{"foo", "baz", "notfound"})
	if string(res) != string(Encode(int64(1), false)) {
		t.Errorf("expected 1, got %s", res)
	}

	// Test: no args
	res = cmdEXISTS([]string{})
	if string(res) != string(Encode(errors.New("(error) ERR wrong number of arguments for 'EXISTS' command"), false)) {
		t.Errorf("expected error for no args, got %s", res)
	}
}

func TestCmdPING(t *testing.T) {
	res := cmdPING([]string{})
	if string(res) != string(Encode("PONG", true)) {
		t.Errorf("expected PONG, got %s", res)
	}
	res = cmdPING([]string{"hello"})
	if string(res) != string(Encode("hello", false)) {
		t.Errorf("expected echo, got %s", res)
	}
	res = cmdPING([]string{"a", "b"})
	if string(res) != string(Encode(errors.New("ERR wrong number of arguments for 'ping' command"), false)) {
		t.Errorf("expected error, got %s", res)
	}
}

func TestCmdSET(t *testing.T) {
	d := setupDictStore()
	dictStore = d
	res := cmdSET([]string{"foo", "bar"})
	if string(res) != string(constant.RespOk) {
		t.Errorf("expected OK, got %s", res)
	}
	res = cmdSET([]string{"foo"})
	if string(res) != string(Encode(errors.New("(error) ERR wrong number of arguments for 'SET' command"), false)) {
		t.Errorf("expected error, got %s", res)
	}
}

func TestCmdGET(t *testing.T) {
	d := setupDictStore()
	d.Set("foo", d.NewObj("foo", "bar", -1))
	dictStore = d
	res := cmdGET([]string{"foo"})
	if string(res) != string(Encode("bar", false)) {
		t.Errorf("expected bar, got %s", res)
	}
	res = cmdGET([]string{"notfound"})
	if string(res) != string(constant.RespNil) {
		t.Errorf("expected nil, got %s", res)
	}
	res = cmdGET([]string{})
	if string(res) != string(Encode(errors.New("(error) ERR wrong number of arguments for 'GET' command"), false)) {
		t.Errorf("expected error, got %s", res)
	}
}

func TestCmdDEL(t *testing.T) {
	d := setupDictStore()
	d.Set("foo", d.NewObj("foo", "bar", -1))
	d.Set("baz", d.NewObj("baz", "qux", -1))
	dictStore = d
	res := cmdDEL([]string{"foo", "baz", "notfound"})
	if string(res) != string(Encode(int64(2), false)) {
		t.Errorf("expected 2, got %s", res)
	}
	res = cmdDEL([]string{})
	if string(res) != string(Encode(errors.New("(error) ERR wrong number of arguments for 'DEL' command"), false)) {
		t.Errorf("expected error, got %s", res)
	}
}

func TestCmdExpire(t *testing.T) {
	d := setupDictStore()
	d.Set("foo", d.NewObj("foo", "bar", -1))
	dictStore = d
	res := cmdEXPIRE([]string{"foo", "10"})
	if string(res) != string(constant.RespOk) {
		t.Errorf("expected OK, got %s", res)
	}
	res = cmdEXPIRE([]string{"foo"})
	if string(res) != string(Encode(errors.New("(error) ERR wrong number of arguments for 'EXPIRE' command"), false)) {
		t.Errorf("expected error, got %s", res)
	}
	res = cmdEXPIRE([]string{"foo", "-1"})
	if string(res) != string(Encode(errors.New("(error) ERR value is not an integer or out of range"), false)) {
		t.Errorf("expected error, got %s", res)
	}
}

func TestCmdTTL(t *testing.T) {
	d := setupDictStore()
	dictStore = d

	// Test non-existent key
	res := cmdTTL([]string{"nonexistent"})
	if string(res) != string(constant.TtlKeyNotExist) {
		t.Errorf("expected key not exist response for nonexistent key, got %s", res)
	}

	// Test key with no expiry
	d.Set("foo", d.NewObj("foo", "bar", -1))
	res = cmdTTL([]string{"foo"})
	if string(res) != string(constant.TtlKeyExistNoExpire) {
		t.Errorf("expected no expire response for key without TTL, got %s", res)
	}

	// Test key with expiry
	d.Set("temp", d.NewObj("temp", "value", 5000)) // 5 seconds
	res = cmdTTL([]string{"temp"})
	// Convert response to number for approximate comparison
	var ttl int64
	err := DecodeInt64(res, &ttl)
	if err != nil {
		t.Errorf("failed to decode TTL response: %v", err)
	}
	if ttl < 4 || ttl > 5 { // Allow for small timing differences
		t.Errorf("expected TTL around 5 seconds, got %d", ttl)
	}

	// Test expired key
	d.Set("expired", d.NewObj("expired", "value", -1))
	d.SetExpiry("expired", -1000) // Set to past
	res = cmdTTL([]string{"expired"})
	if string(res) != string(constant.TtlKeyNotExist) {
		t.Errorf("expected key not exist response for expired key, got %s", res)
	}

	// Test wrong number of arguments
	res = cmdTTL([]string{})
	if string(res) != string(Encode(errors.New("(error) ERR wrong number of arguments for 'TTL' command"), false)) {
		t.Errorf("expected error for no arguments, got %s", res)
	}
}

// mockWriter implements io.Writer for testing
type mockWriter struct {
	written []byte
	err     error
}

func (w *mockWriter) Write(p []byte) (n int, err error) {
	if w.err != nil {
		return 0, w.err
	}
	w.written = append(w.written, p...)
	return len(p), nil
}

func TestExecuteAndResponse(t *testing.T) {
	d := setupDictStore()
	dictStore = d

	testCases := []struct {
		name          string
		cmd           *Command
		setup         func()
		expectedWrite []byte
		writerErr     error
		verify        func(error, []byte)
	}{
		{
			name: "PING command",
			cmd: &Command{
				Cmd:  "PING",
				Args: []string{},
			},
			expectedWrite: Encode("PONG", true),
			verify: func(err error, written []byte) {
				if err != nil {
					t.Errorf("expected no error for PING, got %v", err)
				}
				if string(written) != string(Encode("PONG", true)) {
					t.Errorf("expected PONG response, got %s", written)
				}
			},
		},
		{
			name: "SET command",
			cmd: &Command{
				Cmd:  "SET",
				Args: []string{"key", "value"},
			},
			expectedWrite: constant.RespOk,
			verify: func(err error, written []byte) {
				if err != nil {
					t.Errorf("expected no error for SET, got %v", err)
				}
				if string(written) != string(constant.RespOk) {
					t.Errorf("expected OK response, got %s", written)
				}
				obj := dictStore.Get("key")
				if obj == nil || obj.Value != "value" {
					t.Error("SET command failed to store value")
				}
			},
		},
		{
			name: "GET command",
			cmd: &Command{
				Cmd:  "GET",
				Args: []string{"key"},
			},
			setup: func() {
				d.Set("key", d.NewObj("key", "value", -1))
			},
			expectedWrite: Encode("value", false),
			verify: func(err error, written []byte) {
				if err != nil {
					t.Errorf("expected no error for GET, got %v", err)
				}
				if string(written) != string(Encode("value", false)) {
					t.Errorf("expected value response, got %s", written)
				}
			},
		},
		{
			name: "Write error",
			cmd: &Command{
				Cmd:  "PING",
				Args: []string{},
			},
			writerErr: errors.New("write error"),
			verify: func(err error, written []byte) {
				if err == nil || err.Error() != "write error" {
					t.Errorf("expected write error, got %v", err)
				}
			},
		},
		{
			name: "Unknown command",
			cmd: &Command{
				Cmd:  "UNKNOWN",
				Args: []string{},
			},
			expectedWrite: []byte("-CMD NOT FOUND\r\n"),
			verify: func(err error, written []byte) {
				if err != nil {
					t.Errorf("expected no error for unknown command, got %v", err)
				}
				if string(written) != "-CMD NOT FOUND\r\n" {
					t.Errorf("expected command not found response, got %s", written)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}

			writer := &mockWriter{err: tc.writerErr}
			err := executeAndResponseWithWriter(tc.cmd, writer)
			tc.verify(err, writer.written)
		})
	}
}

// executeAndResponseWithWriter is a test helper that uses an io.Writer instead of a file descriptor
func executeAndResponseWithWriter(cmd *Command, writer io.Writer) error {
	var res []byte

	switch cmd.Cmd {
	case "PING":
		res = cmdPING(cmd.Args)
	case "SET":
		res = cmdSET(cmd.Args)
	case "GET":
		res = cmdGET(cmd.Args)
	case "TTL":
		res = cmdTTL(cmd.Args)
	case "DEL":
		res = cmdDEL(cmd.Args)
	case "EXPIRE":
		res = cmdEXPIRE(cmd.Args)
	case "EXISTS":
		res = cmdEXISTS(cmd.Args)
	case "CMS.INITBYDIM":
		res = cmdCMSINITBYDIM(cmd.Args)
	case "CMS.INITBYPROB":
		res = cmdCMSINITBYPROB(cmd.Args)
	case "CMS.INCRBY":
		res = cmdCMSINCRBY(cmd.Args)
	case "CMS.QUERY":
		res = cmdCMSQUERY(cmd.Args)
	case "BF.RESERVE":
		res = cmdBFRESERVE(cmd.Args)
	case "BF.MADD":
		res = cmdBFMADD(cmd.Args)
	case "BF.EXISTS":
		res = cmdBFEXISTS(cmd.Args)
	default:
		res = []byte("-CMD NOT FOUND\r\n")
	}
	_, err := writer.Write(res)
	return err
}

func TestCmdCMSINITBYDIM(t *testing.T) {
	// Clear the store before each test
	cmsStore = make(map[string]*data_structure.CMS)

	// Test case 1: Valid arguments
	res := cmdCMSINITBYDIM([]string{"mycms", "100", "5"})
	assert.Equal(t, string(constant.RespOk), string(res))
	assert.NotNil(t, cmsStore["mycms"])

	// Test case 2: Key already exists
	res = cmdCMSINITBYDIM([]string{"mycms", "200", "10"})
	assert.Equal(t, string(Encode(errors.New("CMS: key already exists"), false)), string(res))

	// Test case 3: Wrong number of arguments
	res = cmdCMSINITBYDIM([]string{"mycms", "100"})
	assert.Equal(t, string(Encode(errors.New("(error) ERR wrong number of arguments for 'CMS.INITBYDIM' command"), false)), string(res))

	// Test case 4: Invalid width
	res = cmdCMSINITBYDIM([]string{"mycms2", "abc", "5"})
	assert.Contains(t, string(res), "width must be a integer number")

	// Test case 5: Invalid height
	res = cmdCMSINITBYDIM([]string{"mycms3", "100", "xyz"})
	assert.Contains(t, string(res), "height must be a integer number")
}

func TestCmdCMSINITBYPROB(t *testing.T) {
	// Clear the store before each test
	cmsStore = make(map[string]*data_structure.CMS)

	// Test case 1: Valid arguments
	res := cmdCMSINITBYPROB([]string{"mycms", "0.01", "0.001"})
	assert.Equal(t, string(constant.RespOk), string(res))
	assert.NotNil(t, cmsStore["mycms"])

	// Test case 2: Key already exists
	res = cmdCMSINITBYPROB([]string{"mycms", "0.01", "0.001"})
	assert.Equal(t, string(Encode(errors.New("CMS: key already exists"), false)), string(res))

	// Test case 3: Wrong number of arguments
	res = cmdCMSINITBYPROB([]string{"mycms", "0.01"})
	assert.Equal(t, string(Encode(errors.New("(error) ERR wrong number of arguments for 'CMS.INITBYPROB' command"), false)), string(res))

	// Test case 4: Invalid error rate (not float)
	res = cmdCMSINITBYPROB([]string{"mycms2", "abc", "0.001"})
	assert.Contains(t, string(res), "errRate must be a floating point number")

	// Test case 5: Invalid error rate (out of range)
	res = cmdCMSINITBYPROB([]string{"mycms3", "1.5", "0.001"})
	assert.Contains(t, string(res), "invalid overestimation value")

	// Test case 6: Invalid probability (not float)
	res = cmdCMSINITBYPROB([]string{"mycms4", "0.01", "xyz"})
	assert.Contains(t, string(res), "probability must be a floating poit number")

	// Test case 7: Invalid probability (out of range)
	res = cmdCMSINITBYPROB([]string{"mycms5", "0.01", "1.5"})
	assert.Contains(t, string(res), "invalid prob value")
}

func TestCmdCMSINCRBYAndCMSQUERY(t *testing.T) {
	// Clear the store before each test
	cmsStore = make(map[string]*data_structure.CMS)

	// Initialize a CMS filter
	cmdCMSINITBYDIM([]string{"mycms", "100", "5"})

	// Test IncrBy valid operations
	res := cmdCMSINCRBY([]string{"mycms", "item1", "5", "item2", "10"})
	assert.Contains(t, string(res), "5")  // Check for item1 count
	assert.Contains(t, string(res), "10") // Check for item2 count

	// Test Query valid operations
	res = cmdCMSQUERY([]string{"mycms", "item1", "item2", "item3"})
	assert.Contains(t, string(res), "5")  // item1 count
	assert.Contains(t, string(res), "10") // item2 count
	assert.Contains(t, string(res), "0")  // item3 (non-existent) count

	// Test IncrBy: key does not exist
	res = cmdCMSINCRBY([]string{"nonexistent", "item1", "5"})
	assert.Equal(t, string(Encode(errors.New("CMS: key does not exist"), false)), string(res))

	// Test Query: key does not exist
	res = cmdCMSQUERY([]string{"nonexistent", "item1"})
	assert.Equal(t, string(Encode(errors.New("CMS: key does not exist"), false)), string(res))

	// Test IncrBy: wrong number of arguments
	res = cmdCMSINCRBY([]string{"mycms", "item1"})
	assert.Equal(t, string(Encode(errors.New("(error) ERR wrong number of arguments for 'CMS.INCBY' command"), false)), string(res))
	res = cmdCMSINCRBY([]string{"mycms", "item1", "5", "item2"})
	assert.Equal(t, string(Encode(errors.New("(error) ERR wrong number of arguments for 'CMS.INCBY' command"), false)), string(res))

	// Test IncrBy: invalid increment value
	res = cmdCMSINCRBY([]string{"mycms", "item1", "abc"})
	assert.Contains(t, string(res), "increment must be a non negative integer number")

	// Test Query: wrong number of arguments
	res = cmdCMSQUERY([]string{"mycms"})
	assert.Equal(t, string(Encode(errors.New("(error) ERR wrong number of arguments for 'CMS.QUERY' command"), false)), string(res))

	// Test IncrBy overflow (This test case requires a large increment and might be slow or difficult to reliably simulate without a mock CMS)
	// For now, assume IncrBy handles MaxUint32 internally, and the command correctly returns the overflow message.
	// To properly test overflow, we'd need to mock data_structure.CMS.IncrBy to return MaxUint32.
	// For simplicity, we'll skip direct overflow simulation for now and rely on the underlying data_structure tests.

	// Example for a large increment that might cause overflow if not handled (conceptual)
	// cmsStore["mycms"].IncrBy("large_item", math.MaxUint32 - 1)
	// res = cmdCMSINCRBY([]string{"mycms", "large_item", "2"})
	// assert.Contains(t, string(res), "CMS: INCRBY overflow")
}

func TestCmdBFRESERVE(t *testing.T) {
	// Clear the store before each test
	bloomStore = make(map[string]*data_structure.Bloom)

	// Test case 1: Valid arguments
	res := cmdBFRESERVE([]string{"mybf", "0.01", "100"})
	assert.Equal(t, string(constant.RespOk), string(res))
	assert.NotNil(t, bloomStore["mybf"])

	// Test case 2: Key already exists
	res = cmdBFRESERVE([]string{"mybf", "0.01", "100"})
	assert.Contains(t, string(res), "Bloom filter with key 'mybf' already exist")

	// Test case 3: Wrong number of arguments (too few)
	res = cmdBFRESERVE([]string{"mybf2", "0.01"})
	assert.Contains(t, string(res), "ERR wrong number of arguments for 'BF.RESERVE' command")

	// Test case 4: Valid number of arguments (with unimplemented options)
	res = cmdBFRESERVE([]string{"mybf3", "0.01", "100", "NONEXIST", "1"})
	assert.Equal(t, string(constant.RespOk), string(res))

	// Test case 5: Invalid error rate
	res = cmdBFRESERVE([]string{"mybf4", "abc", "100"})
	assert.Contains(t, string(res), "error rate must be a floating point number")

	// Test case 6: Invalid capacity
	res = cmdBFRESERVE([]string{"mybf5", "0.01", "xyz"})
	assert.Contains(t, string(res), "capacity must be an integer number")
}

func TestCmdBFMADDAndBFEXISTS(t *testing.T) {
	// Clear the store before each test
	bloomStore = make(map[string]*data_structure.Bloom)

	// Test case 1: MADD to a non-existent bloom filter (should auto-create)
	res := cmdBFMADD([]string{"mybf", "item1", "item2"})
	assert.Contains(t, string(res), "1")
	assert.NotNil(t, bloomStore["mybf"])

	// Verify existence
	res = cmdBFEXISTS([]string{"mybf", "item1"})
	assert.Equal(t, string(constant.RespOne), string(res))
	res = cmdBFEXISTS([]string{"mybf", "item2"})
	assert.Equal(t, string(constant.RespOne), string(res))
	res = cmdBFEXISTS([]string{"mybf", "nonexistent"})
	assert.Equal(t, string(constant.RespZero), string(res))

	// Test case 2: MADD to an existing bloom filter
	res = cmdBFMADD([]string{"mybf", "item3"})
	assert.Contains(t, string(res), "1")

	// Verify existence of newly added item
	res = cmdBFEXISTS([]string{"mybf", "item3"})
	assert.Equal(t, string(constant.RespOne), string(res))

	// Test case 3: BF.EXISTS on a non-existent bloom filter
	res = cmdBFEXISTS([]string{"nonexistentbf", "item"})
	assert.Equal(t, string(constant.RespZero), string(res))

	// Test case 4: BF.MADD wrong number of arguments
	res = cmdBFMADD([]string{"mybf"})
	assert.Contains(t, string(res), "ERR wrong number of arguments for 'BF.MADD' command")

	// Test case 5: BF.EXISTS wrong number of arguments
	res = cmdBFEXISTS([]string{"mybf"})
	assert.Contains(t, string(res), "ERR wrong number of arguments for 'BF.EXISTS' command")
}
