# DNS Redirector

## Possible messages

Before any message gets processed, the following command needs to be invoked:

```go
// removeDomain removes the last 4 subdomains from the domain
// e.g. "projectId.sessionId.<4 random characters>.subdomain.domain.com." -> "projectId.sessionId"
func (h *Handler) removeDomain(domain string) string {
    parts := strings.Split(domain, ".")
    if len(parts) < 5 {
        return ""
    }
    return strings.Join(parts[:len(parts)-5], ".")
}
```

### AAAA

Session registration. Request looks as follows:

```plaintext
projectId.sessionId.<4 random characters>.subdomain.domain.com
```

If session is successfully registered (or already registered), the response is:

```plaintext
4efc:3425:412b:0bc0:99a1:9f87:d39f:9c84
```

### TXT

Any server response is Base64 encoded. Results are direct JSON objects from gRPC response, i.e.:

```protobuf
message CommandInfoResponse {
    string command_id = 1;
    string type = 2;
    int64 command_length = 3;
}

message CommandDetailsChunkResponse {
    string data = 1;
    bool is_last_chunk = 2;
}
```

#### Querying for new commands

Request looks as follows:

```plaintext
<SessionID>.<4 random characters>.subdomain.domain.com
```

#### Querying for command details

Request looks as follows:

```plaintext
<char>.<offset>.<CommandID>.<4 random characters>.subdomain.domain.com
```

Char is necessary to conform with DNS query not starting with a number. Offset is the offset of the chunk to be retrieved. If the chunk is the last one, the response will contain `is_last_chunk` set to `true`.

### A

It always responds with `183.216.123.191` if there is no error.

#### Setting result length

Before exfiltrating the data, result length needs to be set. Request looks as follows:

```plaintext
<char>.<result length>.<CommandID>.<4 random characters>.subdomain.domain.com
```

Once again, char is necessary to conform with DNS query not starting with a number.

#### Exfiltrating data

Request looks as follows:

```plaintext
<data chunkS with multiple dots>.<offset>.<CommandID>.<4 random characters>.subdomain.domain.com
```

Chunks are hex encoded.
