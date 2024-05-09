# Message Broker

Message Brokers are backbones of event-driven architectures. They enable communication between services in a decoupled manner. They are responsible for routing messages between services, ensuring that messages are delivered, and providing mechanisms for message persistence.

## Redis

Redis is an open-source (not anymore lol), in-memory data structure store, used as a database, cache, and message broker.

### Redis Streams

Redis Streams is a new data structure for managing message streams. It is a log data structure that allows you to store and consume messages in a fault-tolerant way. It is a perfect fit for building message brokers.
