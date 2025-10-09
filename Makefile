
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

tests-user-svc: 
	go test -v -test.count=1 -test.run=TestUserSvcSuite/UserAPISvc ./tests/user/...

tests-user-http: 
	go test -v -test.count=1  -test.run=TestUserHTTPTestSuite/UserAPIHTTPTest ./tests/user/...

tests-user: tests-user-svc tests-user-http

tests-order-svc:
	go test -v -test.count=1  -test.run=TestOrderSvcSuite/OrderAPISvc ./tests/order/...

tests-order-http:
	go test -v -test.count=1  -test.run=TestOrderHTTPTestSuite/OrderAPIHTTPTest ./tests/order/...

tests-order: tests-order-svc tests-order-http


# NOTE: needs mgmtuser-service running
tests-user-http-manual: 
	go test -v -test.count=1  -test.run=TestUserReqSuite/UserAPIReq ./tests/user/...

# NOTE: needs mgmtorder-service running
tests-order-http-manual:
	go test -v -test.count=1  -test.run=TestOrderReqSuite/OrderAPIReq ./tests/order/...

clean:
	rm -rf bin
	lsof -ti tcp:8080