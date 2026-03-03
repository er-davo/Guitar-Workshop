#!/bin/sh
# kafka-init.sh

echo "Waiting for Kafka to be ready..."

until kafka-topics --bootstrap-server kafka:29092 --list >/dev/null 2>&1; do
  sleep 2
done

echo "Kafka is ready. Creating topics..."

for t in \
  audio.separation.start \
  audio.separation.complete \
  tab.generation.request \
  tab.generation.start
do
  kafka-topics --bootstrap-server kafka:29092 \
    --create \
    --if-not-exists \
    --topic "$t" \
    --partitions 1 \
    --replication-factor 1
done

echo "Kafka initialization finished"