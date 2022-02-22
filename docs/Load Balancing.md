# Load Balancing

nRPC generates one handler for each service. This allows you to use the [queueing](http://nats.io/documentation/concepts/nats-queueing/) feature to do load balancing.

Here is a regular server:

```go
// create the nRPC handler
handler := helloworld.NewGreeterHandler(context.TODO(), natsConn, server)

// do a regular subscribe
subscription, err := natsConn.Subscribe(handler.Subject(), handler.Handler)
```

And here is a server that becomes a queue group member:

```go
// create the nRPC handler (same as before)
handler := helloworld.NewGreeterHandler(context.TODO(), natsConn, server)

// do a queue subscribe
subscription, err := natsConn.QueueSubscribe(handler.Subject(), "groupname", handler.Handler)
```

You can run multiple servers, all having the same group name ("groupname" above), and the incoming requests will be randomly distributed between them.

The client code is not aware of the group name concept.