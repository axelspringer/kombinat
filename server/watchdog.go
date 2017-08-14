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
	"os"
	"os/signal"
	"time"

	"github.com/Sirupsen/logrus"
)

// WatchDog holds the state of the watchdog
type WatchDog struct {
	close chan os.Signal
}

// NewWatchDog wrappes the creation of a new WatchDog
func NewWatchDog(timer *time.Ticker) *WatchDog {
	return mustWatchDog(timer)
}

// mustWatchDog create a new WatchDog
func mustWatchDog(timer *time.Ticker) *WatchDog {
	// generate channel to publish to
	osChan := make(chan os.Signal, 10)

	// starting watchdog, continue to run
	go func() {
		signal.Notify(osChan, os.Interrupt)
		// Wait for signaling from os
		<-osChan
		// stop timer, and exit
		timer.Stop()
		// debug
		logrus.Info("Stopping and cleaning. Bye, Bye!")
		// clean exit
		os.Exit(0)
	}()

	// debug
	logrus.Info("Watchdog started")

	// return ready WatchDog, you can also ignore that
	return &WatchDog{osChan}
}
