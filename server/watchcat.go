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
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"

	amzn "github.com/aws/aws-sdk-go/aws"
	"github.com/axelspringer/kombinat/aws"
	"github.com/axelspringer/kombinat/etcd"
	"github.com/coreos/etcd/client"
)

// WatchCat wraps the state of the watch cat
type WatchCat struct {
	timer       *time.Ticker
	errTimeout  time.Duration
	pollTimeout time.Duration
	ctx         context.Context
}

// NewWatchCat wraps the generation of a new watchcat
func NewWatchCat(ctx context.Context, pollTimeout time.Duration, errTimeout time.Duration) (*WatchCat, error) {
	// create new WatchCat
	wc, err := mustNewWatchCat(ctx, pollTimeout, errTimeout)

	// try to create cat
	if err != nil {
		return nil, ErrWatchCat
	}

	// tickle the cat, the
	wc.startWatchCat()

	// return WatchCat
	return wc, nil
}

// mustTickleWatchCat
func (w *WatchCat) startWatchCat() {
	// init session
	initSession := aws.NewSession(&amzn.Config{})

	// create a new instance
	ec2, err := aws.NewEC2(initSession)

	// exit on error
	if err != nil {
		os.Exit(1)
	}

	// start go routine for checking
	go func() {
		// loop
		for {
			// get ec2 peers
			ec2Peers, err := ec2.GetAutoScalingGroupPeers()

			// exit on error
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1) // methodology, let fail, and restart via systemd
			}

			// get ids
			ec2PeersIds := make([]*string, 0)
			for _, peer := range ec2Peers {
				ec2PeersIds = append(ec2PeersIds, peer.InstanceId)
			}

			// if etcd is enabled
			if EtcdEnable {
				if err := watchEtcd(&w.ctx); err != nil {
					logrus.Errorln(err.Error())
					os.Exit(1)
				}
			}

			// listen to channels
			select {
			case <-w.ctx.Done():
				return
			default:
				// waiting for timer
				<-w.timer.C
			}
		}
	}()

	return
}

// mustNewWatchCat
func mustNewWatchCat(ctx context.Context, pollTimeout time.Duration, errTimeout time.Duration) (*WatchCat, error) {
	// create time
	pollTimer := time.NewTicker(pollTimeout)

	// return a new WatchCat
	return &WatchCat{pollTimer, pollTimeout, errTimeout, ctx}, nil
}

// watchEtcd is watching an etcd cluster
func watchEtcd(c *context.Context) error {
	// create client
	cli, err := etcd.New(client.Config{
		Endpoints: []string{EtcdEndpoint},
	})

	// check error
	if err != nil {
		return err
	}

	// get members api
	etcdMembersAPI := client.NewMembersAPI(*cli)

	// check for etcd peers
	etcdMembers, err := etcdMembersAPI.List(*c)

	// exit on error
	if err != nil {
		return err
	}

	// list unhealthy etcd members
	unhealtyMembers := make([]client.Member, 0)

	// search in etcd members
	for _, m := range etcdMembers {
		if err := etcd.CheckMembersHealth(&m); err != nil {
			unhealtyMembers = append(unhealtyMembers, m)
		}
	}

	// remove unhealthy members
	for _, m := range unhealtyMembers {
		// remove unhealthy members
		err := etcdMembersAPI.Remove(*c, m.ID)
		// check for error in remove, but fire and forget
		if err != nil {
			// log message, but forget
			fmt.Println(err.Error())
		}
	}

	return nil
}
