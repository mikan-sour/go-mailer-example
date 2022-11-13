#!/bin/bash

export $(grep -v '^#' .env | xargs)
KAFKA_DOMAIN=$KAFKA_HOST:$KAFKA_PORT

# Create topic
sh /opt/bitnami/kafka/bin/kafka-topics.sh  \
    --create --topic $KAFKA_LISTEN_TOPIC \
    --partitions 3 \
    --replication-factor 1 \
    --bootstrap-server $KAFKA_DOMAIN

# List
sh /opt/bitnami/kafka/bin/kafka-topics.sh  --list  --bootstrap-server $KAFKA_DOMAIN

# Produce
kafka-console-producer.sh --broker-list $KAFKA_DOMAIN --topic $KAFKA_LISTEN_TOPIC