# Kafka Installed Version Collector

Download a Website, extract the version and publish the installed version to a kafka topic.

## Run version collector

```bash
go run main.go \
-kafka-brokers=kafka:9092 \
-kafka-topic=application_version_installed \
-kafka-schema-registry-url=http://localhost:8081 \
-app-name=Confluence \
-app-regex='<meta\s+name="ajs-version-number"\s+content="([^"]+)">' \
-app-url=https://confluence.benjamin-borbe.de \
-v=2
```
