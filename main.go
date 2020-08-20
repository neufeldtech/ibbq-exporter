/*
   Copyright 2018 the original author or authors
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at
       http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package main

import (
	"context"

	"time"

	"github.com/go-ble/ble"
	"github.com/sworisbreathing/go-ibbq/v2"
	"go.uber.org/zap"
)

var logger, _ = zap.NewProduction()

func temperatureReceived(temperatures []float64) {
	logger.Info("yep", zap.Float64("temps", temperatures[0]))
}
func batteryLevelReceived(batteryLevel int) {
	logger.Info("Received battery data batteryPct", zap.Int("batterylevel", batteryLevel))
}
func statusUpdated(status ibbq.Status) {
	logger.Info("Status updated status")
}

func disconnectedHandler(cancel func(), done chan struct{}) func() {
	return func() {
		logger.Info("Disconnected")
		cancel()
		close(done)
	}
}

func main() {

	var err error
	logger.Info("initializing context")
	ctx1, cancel := context.WithCancel(context.Background())
	defer cancel()
	registerInterruptHandler(cancel)
	ctx := ble.WithSigHandler(ctx1, cancel)
	logger.Info("context initialized")
	var bbq ibbq.Ibbq
	logger.Info("instantiating ibbq struct")
	done := make(chan struct{})
	var config ibbq.Configuration
	if config, err = ibbq.NewConfiguration(60*time.Second, 5*time.Minute); err != nil {
		logger.Error("woops", zap.Error(err))
	}
	if bbq, err = ibbq.NewIbbq(ctx, config, disconnectedHandler(cancel, done), temperatureReceived, batteryLevelReceived, statusUpdated); err != nil {
		logger.Error("woops", zap.Error(err))
	}
	logger.Info("instantiated ibbq struct")
	logger.Info("Connecting to device")
	if err = bbq.Connect(); err != nil {
		logger.Fatal("Error connecting to device")
	}
	logger.Info("Connected to device")
	<-ctx.Done()
	<-done
	logger.Info("Exiting")
}
