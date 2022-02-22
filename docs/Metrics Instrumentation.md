# Metrics Instrumentation

The nRPC protoc plugin can generate code that has metrics instrumentation. Currently, only [Prometheus](https://prometheus.io/) is supported. To generate code with Prometheus metrics built-in, use:

```
$ protoc --go_out=. --nrpc_out=plugins=prometheus:. helloworld.proto
```

Here is the generated file: [helloworld.nrpc.go](https://github.com/rapidloop/nrpc/blob/master/examples/metrics_helloworld/helloworld/helloworld.nrpc.go)

The metrics must be hosted via the `promhttp` handler so that it can be scraped by Prometheus. Here is an example [server](https://github.com/rapidloop/nrpc/blob/master/examples/metrics_helloworld/metrics_greeter_server/main.go) and an example [client](https://github.com/rapidloop/nrpc/blob/master/examples/metrics_helloworld/metrics_greeter_client/main.go).

The following metrics are generated:

## Server-side

- **nrpc_server_handler_execution_time_seconds** is a φ-quantile with 0.90, 0.95 and 0.99 quantiles. It is the time taken to execute the user-supplied handler code at server-side. The label `service` contains the name of the service and `method` contains the name of the service method.
- **nrpc_server_requests_count** is a counter, incremented for each invocation of the handler. The label `result_type` contains one of the values `success`, `handler_fail`, `protobuf_fail` or `name_fail`. The counter with the appropriate label value for result_type is incremented depending on the result of the invocation. This metric also has the `service` and `method` labels. "handler_fail" is when the user-supplied handler returns an error, "protobuf_fail" is when (un)marshalling fails, and "name_fail" is when an unknown method name is received.

Here is a sample output:

```shell
# HELP nrpc_server_handler_execution_time_seconds The handler execution time for Greeter calls, measured server-side.
# TYPE nrpc_server_handler_execution_time_seconds summary
nrpc_server_handler_execution_time_seconds{method="SayHello",service="Greeter",quantile="0.9"} 0.073982319
nrpc_server_handler_execution_time_seconds{method="SayHello",service="Greeter",quantile="0.95"} 0.097965047
nrpc_server_handler_execution_time_seconds{method="SayHello",service="Greeter",quantile="0.99"} 0.097965047
nrpc_server_handler_execution_time_seconds_sum{method="SayHello",service="Greeter"} 0.38668777499999996
nrpc_server_handler_execution_time_seconds_count{method="SayHello",service="Greeter"} 13
# HELP nrpc_server_requests_count The count of requests handled by the server.
# TYPE nrpc_server_requests_count counter
nrpc_server_requests_count{method="SayHello",result_type="handler_fail",service="Greeter"} 10
nrpc_server_requests_count{method="SayHello",result_type="success",service="Greeter"} 3
```

## Client-side

- **nrpc_client_request_completion_time_seconds** is a φ-quantile with 0.90, 0.95 and 0.99 quantiles. It is the time taken for the generated client code to call a remote method and unmarshal the result. The label `service` contains the name of the service and `method` contains the name of the service method.
- **nrpc_client_calls_count** is a counter, incremented for each invocation of any call. The label `result_type` contains one of the values `success`, `call_fail` or `protobuf_fail`. The counter with the appropriate label value for result_type is incremented depending on the result of the invocation. This metric also has the `service` and `method` labels. "call_fail" is when NATS invocation fails and "protobuf_fail" is when (un)marshalling fails.

Here is a sample output:

```shell
# HELP nrpc_client_calls_count The count of calls made by the client.
# TYPE nrpc_client_calls_count counter
nrpc_client_calls_count{method="SayHello",result_type="success",service="Greeter"} 1
# HELP nrpc_client_request_completion_time_seconds The request completion time for Greeter calls, measured client-side.
# TYPE nrpc_client_request_completion_time_seconds summary
nrpc_client_request_completion_time_seconds{method="SayHello",service="Greeter",quantile="0.9"} 0.000455521
nrpc_client_request_completion_time_seconds{method="SayHello",service="Greeter",quantile="0.95"} 0.000455521
nrpc_client_request_completion_time_seconds{method="SayHello",service="Greeter",quantile="0.99"} 0.000455521
nrpc_client_request_completion_time_seconds_sum{method="SayHello",service="Greeter"} 0.000455521
nrpc_client_request_completion_time_seconds_count{method="SayHello",service="Greeter"} 1
```



(TODO: client-side protobuf marshal errors and NATS timeout errors also get classified as "call_fail". Fix this.)