#!/bin/bash

rmdir -rf /var/git/repos/GitOperatorTestRepo

docker exec kafka kafka-topics --zookeeper zookeeper:2181 \
 --create --topic gitoperator-in \
 --partitions 1 \
 --replication-factor 1 \
#  --min.compaction.lag.ms 86400000 \
#  --max.compaction.lag.ms 432000000 \
#  --retention.ms 86400000


docker exec kafka kafka-topics --zookeeper zookeeper:2181 \
 --create --topic gitoperator-out \
 --partitions 1 \
 --replication-factor 1 \
#  --cleanup.policy compact,delete \
#  --min.compaction.lag.ms 86400000 \
#  --max.compaction.lag.ms 432000000 \
#  --retention.ms 86400000

