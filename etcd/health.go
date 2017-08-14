// Copyright © 2017 Axel Springer SE
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package etcd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/pkg/transport"
)

// CheckMembersHealth is checking the health of an etcd member
func CheckMembersHealth(m *client.Member) error {
	// check for available client urls
	if len(m.ClientURLs) == 0 {
		fmt.Printf("member %s is unreachable: no available published client urls\n", m.ID)
		return nil
	}

	dialTimeout := 30 * time.Second
	tr, err := transport.NewTransport(transport.TLSInfo{}, dialTimeout)
	// check for nil
	if err != nil {
		os.Exit(1)
	}

	// http client
	hc := http.Client{
		Transport: tr,
	}

	// checked := false
	for _, url := range m.ClientURLs {
		// check for health route
		if _, err := hc.Get(url + "/health"); err != nil {
			return ErrEtcdMemberUnhealty
		}
	}

	return nil
}
