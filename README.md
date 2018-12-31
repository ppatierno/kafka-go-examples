# Apache Kafka Go examples

Apache Kafka sample applications in Go language using the [segmentio Kafka Go library](https://github.com/segmentio/kafka-go).

* consumer: consumer application.
* producer: producer application.
* util: utility package used by the consumer/producer applications.

It also provides the related Dockerfile for building Docker container images but it's needed to build the applications before building these images.

Finally, it provides YAML files describing Kubernetes Deployment(s) for the consumer and producer applications.
