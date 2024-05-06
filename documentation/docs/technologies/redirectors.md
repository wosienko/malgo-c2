# Redirectors

## Go

Right now, Redirectors are written in Go, as there is no need for introducing new languages. Go is a statically typed, compiled language that is designed for simplicity and efficiency. It is a great choice for writing high-performance servers. It is a perfect balance between performance and ease of development.

### gRPC

Redirectors are gRPC clients. They communicate with the Gateway Server using gRPC.

### Other technologies

Redirectors might use various other technologies to enable communications over covert channels. For example, HTTP, DNS, ICMP, etc.
