# Specifications

NRPC translates [protocol buffers services](https://developers.google.com/protocol-buffers/docs/proto3#services) into [NATS](https://nats.io/) patterns.

This page describes the rules that applies, and is meant to be a reference for other languages ports (current status is draft).

Most of the time, a rpc call over nats corresponds to a REQ/REP pattern on a specific subject. See also 'Special request/replies' below.

# Subjects

A given method has a specific subject that may optionally contains 'parameters', ie some parts that can carry arbitrary strings and are not used to determinate which method is being called.

The general syntax of a subject is the following:

```shell
[<pkg-subject>.][<pkg-param>.[<pkg-param>.[...]]]<svc-subject>.[<svc-param>.[<svc-param>.[...]]].<mt-subject>[.<mt-param>[.<mt-param>[...]]][.<encoding>]
```

where:

- `<pkg-subject>` is the package-related part of the subject. The default is the package name (nothing if it is unnamed), and can be modified with the `nrpc.packageSubject` file option. It can contain dots.
- `<pkg-param>` is a package-related parameter. The default is nothing, parameters can be added with the repeatable `nrpc.packageSubjectParams` file option.
- `<svc-subject>` is the service-related part of the subject. The default is the service name. If the file option `nrpc.serviceSubjectRule` is `LOWER`, the service name is lowered. In any case, it can be overridden by the `nrpc.serviceSubject` service option, in which dots are allowed.
- `<svc-param>` is a service-related parameter. The default is nothing, parameters can be added with the repeatable `nrpc.serviceSubjectParams` service option.
- `<mt-subject>` is the method-related part of the subject. The default is the method name on which the file option `nrpc.methodSubjectRule` is applied: `COPY` keeps the name intact, and `LOWER` put all letter to lowercase. In any case, it can be overridden by the `nrpc.methodSubject` method option, in which dots are *not* allowed.
- `<mt-param>` is a method-related parameter. The default is nothing, parameters can be added with the repeatable `nrpc.methodSubjectParams` method option.
- `<encoding>` is an optional encoding indication. If omitted, the message are encoded with protobuf. If possible, it is recommended that implementations support the `json` encoding.

# Message encoding

The primary encoding of the message is protobuf.

In some case however, protobuf may not be an option for a given client. To address this, an alternative json encoding is possible, in which case a '.json' suffix is added to the subject.

## Request

A request encoding is the direct encoding of the request message in the wanted encoding.

## Reply

A successful reply is the direct encoding of the reply message in protobuf and json.

The format of an error differs depending in the target encoding:

- protobuf: A '0x00' byte followed by a protobuf-encoded nrpc.Error
- json: A json object with a single '**error**' attribute. The value of this attribute is a json-encoded nrpc.Error.

In both case, an error message is easy to identify, so the decoding can be done in an efficient manner.

# Special request/replies

## Void

`nrpc.Void` can be used as a input or output type of a method as a hint to the code-generator that no particular data is transported in the request or the reply.

The generated code should use this hint to remove unnecessary arguments to function signatures.

Note: returning `nrpc.Void` changes nothing to the call pattern, and allow errors to be returned.

## NoReply

A method that `returns (nrpc.NoReply)` will not return anything from the server side, not even an error.

This special case corresponds to a simple `PUB` operation to the method subject instead of a `REQ/REP`.

The generated code will use the nats 'Publish' method instead of 'Request' for such calls.

## NoRequest

A method that take a `NoRequest` input is not a method that can be called but a stream to which a client can subscribe. It is the exact opposite of a method that returns `nrpc.NoReply`.

The generated code should provide helpers to easily subscribe to the method subject.

## Streamed replies

A method on which the option `nrpc.streamd_reply` is set to `true` will send back a stream of replies instead of just one. It also sends keep-alive messages and expects a heart-beat from the client.

The server-side behavior is the following:

- send messages to the request inbox, until the call completions
- at least every second, a message should be sent. If no reply can be sent, a keep-alive message is sent.
- a keep-alive message is a single '0' byte string
- If the client receives no message for more than a few seconds, the call is considered canceled.
- the stream is close when a nrpc.Error is sent to the inbox
- a special error type `EOS` identifies a normal end of stream
- expects to receive a `nrpc.Heartbeat` message on the `<inbox>.heartbeat` subject, where `<inbox>` is the inbox of the request
- If a heartbeat message with `lastbeat` set to `true` is received, the method call is canceled.
- If no heartbeat is received for more than 5 seconds, the method call is canceled