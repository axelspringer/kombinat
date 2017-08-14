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

package server

import (
	"context"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	DefaultPollInterval = 60 // how often to poll api
	DefaultErrTimeout   = 10 // back-off timeout on error
	DefaultEtcdEndpoint = "http://127.0.0.1:2379"
	DefaultEtcdEnable   = false
)

var (
	Verbose      bool
	EtcdEnable   bool
	EtcdEndpoint string
	PollInterval int
	ErrTimeout   int
)

// Run is taking the commands is running the server
func Run(cmd *cobra.Command, args []string) error {
	// Timer
	pollTimer := time.NewTicker(1000 * time.Millisecond)

	// WatchDog
	wd := NewWatchDog(pollTimer)

	// WatchCat
	ctx, cancel := context.WithCancel(context.TODO())
	_, err := NewWatchCat(ctx, 10*time.Second, 10*time.Second) // has to changed
	defer cancel()

	// ErrWatchCat
	if err == (ErrWatchCat) {
		logrus.Info(ErrWatchCat)
	}

	// wait for Watchdog
	<-wd.close

	// return default no error
	return nil
}
