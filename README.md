# Apache Kafka Go examples

Apache Kafka sample applications in Go language using the [segmentio](https://github.com/segmentio/kafka-go) and [shopify sarama](https://github.com/Shopify/sarama) Kafka Go libraries.

* segmentio:
    * consumer: consumer application.
    * producer: producer application.
* sarama:
    * sync-producer: synchronous producer application
    * partition-consumer: partition consumer application
* util: utility package used by the consumer/producer applications.

It also provides the related Dockerfile for building Docker container images but it's needed to build the applications before building these images.

Finally, it provides YAML files describing Kubernetes Deployment(s) for the consumer and producer applications.
