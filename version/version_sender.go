// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/bborbe/kafka-installed-version-collector/avro"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/seibert-media/go-kafka/schema"
)

type Sender struct {
	KafkaBrokers   string
	KafkaTopic     string
	SchemaRegistry interface {
		SchemaId(subject string, schema string) (uint32, error)
	}
}

func (s *Sender) Send(ctx context.Context, version avro.ApplicationVersionInstalled) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	client, err := sarama.NewClient(strings.Split(s.KafkaBrokers, ","), config)
	if err != nil {
		return errors.Wrap(err, "create client failed")
	}
	defer client.Close()

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return errors.Wrap(err, "create sync producer failed")
	}
	defer producer.Close()

	schemaId, err := s.SchemaRegistry.SchemaId(fmt.Sprintf("%s-value", s.KafkaTopic), version.Schema())
	if err != nil {
		return errors.Wrap(err, "get schema id failed")
	}
	buf := &bytes.Buffer{}
	if err := version.Serialize(buf); err != nil {
		return errors.Wrap(err, "serialize version failed")
	}
	partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: s.KafkaTopic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%s-%s-%s", version.Url, version.App, version.Version)),
		Value: &schema.AvroEncoder{SchemaId: schemaId, Content: buf.Bytes()},
	})
	if err != nil {
		return errors.Wrap(err, "send message to kafka failed")
	}
	glog.V(3).Infof("send message successful to %s with partition %d offset %d", s.KafkaTopic, partition, offset)
	return nil
}
