# Gateway Server

Gateway server is, up to a certain point, very similar to the Websocket Server - it also creates commands that are sent to appropriate streams.

Where it differs is the usage of gRPC.

## gRPC

```proto
service GatewayService {
  rpc RegisterNewSession(RegisterNewSessionRequest) returns (EmptyResponse) {} // Synchronous
  rpc SessionHeartbeat(SessionHeartbeatRequest) returns (EmptyResponse) {} // Asynchronous
  rpc CommandInfo(CommandInfoRequest) returns (CommandInfoResponse) {} // Synchronous
  rpc CommandDetailsChunk(CommandDetailsChunkRequest) returns (CommandDetailsChunkResponse) {} // Synchronous

  rpc ResultInfo(ResultInfoRequest) returns (EmptyResponse) {} // Synchronous
  rpc ResultDetailsChunk(ResultDetailsChunkRequest) returns (EmptyResponse) {} // Synchronous
}

message RegisterNewSessionRequest {
    string session_id = 1;
    string project_id = 2;
}

message SessionHeartbeatRequest {
    string session_id = 1;
}

message CommandInfoRequest {
    string session_id = 1;
}

message CommandDetailsChunkRequest {
  string command_id = 1;
  int64 offset = 2;
  int64 length = 3;
}

message ResultInfoRequest {
    string command_id = 1;
    int64 length = 2;
}

message ResultDetailsChunkRequest {
    string command_id = 1;
    int64 offset = 2;
    bytes data = 3;
}

message CommandInfoResponse {
    string command_id = 1;
    string type = 2;
    int64 command_length = 3;
}

message CommandDetailsChunkResponse {
    string data = 1;
    bool is_last_chunk = 2;
}

message EmptyResponse {}
```
