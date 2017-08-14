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
	"github.com/coreos/etcd/client"
)

// New creates a new client
func New(cfg client.Config) (*client.Client, error) {
	// creating new client
	cli, err := mustNewClient(&cfg)

	// check for error in creating client
	if err != nil {
		return nil, err // forward error
	}

	// returning new client
	return cli, nil
}

// mustNewClient returns a new client, and error
func mustNewClient(cfg *client.Config) (*client.Client, error) {
	// really creating new client
	cli, err := client.New(*cfg)

	// check for errors on client
	if err != nil {
		return nil, ErrNoEtcdClient
	}

	// returning client
	return &cli, nil
}
