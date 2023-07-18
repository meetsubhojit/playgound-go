# golang api testing with coverage

We generally write unit tests, and run `go test` with cover flags to visualise the coverage our test code has. 
Then we run automation suites in frameworks like karate which hits our service's endpoints. Here we care more about the correctness of the response received. However using `go test -c` (an inbuilt utility to build a go binary in "test" mode) we can additionally make our service to emit `unit test` like coverage reports. 

### how to use this module
run` make test/api` which will build the binary in "test" mode, and then call a http enpoint that our service exposes. Ultimtely showing the coverage report.

there are other make commands which builds, runs the binary as a normal server.

### internals
we have written a simple service which exposes a `/time` http endpoint which gives you the current time. On calling `/time` it hits the controller->api->utils pkgs. `internal` package has all our application code.

`cmd` has two packages, `cmd/service` pkg contains the actual `main.go` which starts the service's gin server under `8080`. Consider this as our main binary which will be built during CI and deployed to various environments. 

`cmd/test-service` pkg's `main_test.go` also does what the above `cmd/service/main.go` do, expect it does this under a "test wrapper". This gives the hook to the `go test` runtime to generate coverage report.

you fill notice that all functions have 100% coverage expect `NotUsedFunction` because its internationnaly not used in the request path for `/time`.

#### what is happening in cmd/test-service/main_test.go
1. we start the main server in a go routine
2. we run our api "tests". this happens under function `runApiTest`. This is nothing but a http call. This http call is to `localhost:8080/time` which hits the server created in step 1. Consider this as equivalent to running some "karate" like tests from outside.
3. test ends, closing the gin server created in step 1.

PS. compare the file under `cmd/service` and `cmd/test-service` to get more understanding on the test instrumentation.
