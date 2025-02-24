# Useful tools
I have found the follows tools have helped me the most in my software engineering practice.

## Language
`Go` https://go.dev/ is an easy to learn, yet very capable language suitable for building backend services.

## Serialization, Persistence and Storage
JSON / YAML / TOML serialization and persistence as human-readable text files.

`Go`'s GOB https://pkg.go.dev/encoding/gob for both serialization and persistence.

`protobuf` cross-platform serialization with good support for schema evolution.

BoltDB https://github.com/boltdb/bolt provides an embeddable key-value store with prefix and range search capability.

SQL databases: sqlite, Postgres and MariaDB.

## Communication, Messaging and Streams

NATS https://pkg.go.dev/github.com/nats-io/nats-server/v2@v2.10.25/server provides an embeddable messaging system with pub-sub, request-reply, work queues, streams, key-value and object stores.

MQTT https://mqtt.org/ provides a lightweight messaging system for internet of things applications.

## Vector Databases, RAG and LLMs
chromem-go https://github.com/philippgille/chromem-go vector database with bindings to ollama

ollama https://ollama.com/ allows running LLMs locally.