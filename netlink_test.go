// Copyright 2017 Elasticsearch Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build linux

package libaudit

import (
	"bytes"
	"encoding/binary"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Validate NetlinkClient implements NetlinkSendReceiver.
var _ NetlinkSendReceiver = &NetlinkClient{}

func TestParseNetlinkErrorDataTooShort(t *testing.T) {
	assert.Error(t, ParseNetlinkError(nil), "too short")
}

func TestParseNetlinkErrorErrno(t *testing.T) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, -1*int32(NLE_MSG_TOOSHORT))
	assert.Equal(t, ParseNetlinkError(buf.Bytes()), NLE_MSG_TOOSHORT)
}

func TestNewNetlinkClient(t *testing.T) {
	c, err := NewNetlinkClient(syscall.NETLINK_AUDIT, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	assert.Len(t, c.readBuf, os.Getpagesize())

	// First PID assigned by the kernel will be our actual PID.
	assert.EqualValues(t, os.Getpid(), c.pid)

	c2, err := NewNetlinkClient(syscall.NETLINK_AUDIT, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c2.Close()

	// Second PID assigned by kernel will be random.
	assert.NotEqual(t, 0, c2.pid)
	assert.NotEqual(t, c.pid, c2.pid)
}
