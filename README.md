A small scale experiment to demonstrate that a Go HTTP handler will start as soon as the mux parses the request path, before the server has received the entire request.

# How it works

The client sends a very large payload to the server, but we can choose to use a normal io.Reader or a throttled io.Reader. The throttled reader breaks the payload into 100 chunks and reads a one chunk at a time, sleeping a few tens of ms between each chunk.

## Normal mode

### Client logs:

```
➜  golang_clientserver go run client/main.go
2025/05/02 10:11:22 Request started at second: 22, ns: 342310000
2025/05/02 10:11:22 Using normal mode
2025/05/02 10:11:22 Response status: 200 OK
2025/05/02 10:11:22 Sent payload size: 10.00 MB
2025/05/02 10:11:22 Request took: 35.761791ms
```


### Server logs:

```
2025/05/02 10:11:22 request started at second: 22, ns: 345352000
2025/05/02 10:11:22


2025/05/02 10:11:22 payload sent at 2025-05-02T10:11:22-07:00
2025/05/02 10:11:22 payload received at second: 22, ns: 377848000
2025/05/02 10:11:22 duration to decode payload: 32.49675ms
2025/05/02 10:11:22


2025/05/02 10:11:22 request completed in 32.557041ms
2025/05/02 10:11:22
```

As one would expect, the client initiated the request and sent over a 10MB payload.
The server started reading the payload, spent 32.4ms decoding it and deserializing it into a struct,
and then returned a 200 OK response. Let's try it with a throttled reader.

## Throttled mode

### Client logs:

```
➜  golang_clientserver git:(main) ✗ go run client/main.go --throttled
2025/05/02 10:11:27 Request started at second: 27, ns: 705022000
2025/05/02 10:11:27 Using throttled mode
2025/05/02 10:11:31 Response status: 200 OK
2025/05/02 10:11:31 Sent payload size: 10.00 MB
2025/05/02 10:11:31 Request took: 3.589888375s
```

### Server logs

```
2025/05/02 10:11:27 request started at second: 27, ns: 707670000
2025/05/02 10:11:27


2025/05/02 10:11:31 payload sent at 2025-05-02T10:11:27-07:00
2025/05/02 10:11:31 payload received at second: 31, ns: 294295000
2025/05/02 10:11:31 duration to decode payload: 3.586703875s
2025/05/02 10:11:31


2025/05/02 10:11:31 request completed in 3.586766959s
2025/05/02 10:11:31
```

The server started handling the request at 10:11:27, but spent a full 3.5 seconds between
lines 30 and 37, just receiving and deserializing the payload.
The only difference is the client running in throttled mode vs normal mode, taking extra
time to send the payload. 

This mimics behavior we saw during a period of bandwidth congestion where traces would start
and then show a long period of nothing before starting work.
The only thing the server could have been doing during this period was receiving and deserializing
the request payload.