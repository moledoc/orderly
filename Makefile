
all: clean dir compile

dir: clean
	mkdir -p bin

build:
	go build -o ./bin/mgmtuser ./cmd/mgmtuser/main.go
	go build -o ./bin/mgmtorder ./cmd/mgmtorder/main.go
	go build -o ./bin/orderly ./cmd/orderly/main.go

run-user:
	go run ./cmd/mgmtuser/main.go

run-order:
	go run ./cmd/mgmtorder/main.go

run:
	go run ./cmd/orderly/main.go

up-user: build
	./bin/mgmtuser

up-order: build
	./bin/mgmtorder

up: build
	./bin/orderly

clean:
	rm -rf bin
	lsof -ti tcp:8080

# FUNCTIONAL TESTING

tests-user-svc: 
	go test -v -test.count=1 -test.run=TestUserSvcSuite ./tests/user/...

tests-user-http: 
	go test -v -test.count=1 -test.run=TestUserHTTPTestSuite ./tests/user/...

tests-user-http-manual: 
	go test -v -test.count=1 -test.run=TestUserReqSuite ./tests/user/...

tests-user:
	go test -v -test.count=1 -test.run="TestUser[SvcHTTPTestReq]*Suite" ./tests/user/...

tests-order-svc:
	go test -v -test.count=1 -test.run=TestOrderSvcSuite ./tests/order/...

tests-order-http:
	go test -v -test.count=1 -test.run=TestOrderHTTPTestSuite ./tests/order/...

tests-order-http-manual:
	go test -v -test.count=1 -test.run=TestOrderReqSuite ./tests/order/...

tests-order: 
	go test -v -test.count=1 -test.run="TestOrder[SvcHTTPTestReq]*Suite" ./tests/order/...

tests-all:
	go test -v -test.count=1 -test.skip=".*[SvcHTTPTestReq]*PerformanceSuite" ./tests/...


# PERFORMANCE TESTING

perf-test-user-svc: 
	go test -v -test.count=1 -test.timeout=30m -test.run=TestUserSvcPerformanceSuite ./tests/user/...

perf-test-user-http: 
	go test -v -test.count=1 -test.timeout=30m -test.run=TestUserHTTPTestPerformanceSuite ./tests/user/...

perf-test-user-http-manual: 
	go test -v -test.count=1 -test.timeout=30m -test.run=TestUserReqPerformanceSuite ./tests/user/...

perf-test-user:
	go test -v -test.count=1 -test.timeout=30m -test.run="TestUser[SvcHTTPTestReq]*PerformanceSuite" ./tests/user/...

perf-test-order-svc:
	go test -v -test.count=1 -test.timeout=30m -test.run=TestOrderSvcPerformanceSuite ./tests/order/...

perf-test-order-http:
	go test -v -test.count=1 -test.timeout=30m -test.run=TestOrderHTTPTestPerformanceSuite ./tests/order/...

perf-test-order-http-manual:
	go test -v -test.count=1 -test.timeout=30m -test.run=TestOrderReqPerformanceSuite ./tests/order/...

perf-test-order: 
	go test -v -test.count=1 -test.timeout=30m -test.run="TestOrder[SvcHTTPTestReq]*PerformanceSuite" ./tests/order/...

perf-test-all:
	go test -v -test.timeout=60m -test.count=1 -test.run=".*[SvcHTTPTestReq]*PerformanceSuite" ./tests/...
