// Copyright 2021 Occamshub Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grypereceiver

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

const (
	typeStr = "grype"
)

var (
	defaultInclude = []string{"/"}
)

// NewFactory creates a factory for Grype receiver.
func NewFactory() component.ReceiverFactory {
	return receiverhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		receiverhelper.WithMetrics(createMetricsReceiver))
}

func createDefaultConfig() config.Receiver {
	return &Config{
		ScraperControllerSettings: scraperhelper.ScraperControllerSettings{
			ReceiverSettings:   config.NewReceiverSettings(config.NewComponentID(typeStr)),
			CollectionInterval: 24 * time.Hour,
		},
		Include: defaultInclude,
	}
}

func createMetricsReceiver(
	_ context.Context,
	params component.ReceiverCreateSettings,
	rConf config.Receiver,
	consumer consumer.Metrics,
) (component.MetricsReceiver, error) {
	cfg := rConf.(*Config)

	var ns = newGrypeScraper(params.Logger, cfg)
	scraper, err := scraperhelper.NewScraper(typeStr, ns.Scrape, scraperhelper.WithStart(ns.Start))
	if err != nil {
		return nil, err
	}

	return scraperhelper.NewScraperControllerReceiver(
		&cfg.ScraperControllerSettings, params, consumer,
		scraperhelper.AddScraper(scraper),
	)
}
