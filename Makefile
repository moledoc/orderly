
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

tests-user: 
	go test -v -test.count=1 -test.run=TestUserSvcSuite/UserAPISvc ./tests/user/...

# NOTE: needs mgmtuser-service running
tests-user-http: 
	go test -v -test.count=1  -test.run=TestUserReqSuite/UserAPIReq ./tests/user/...

tests-order:
	go test -v -test.count=1  -test.run=TestOrderSvcSuite/OrderAPISvc ./tests/order/...

# NOTE: needs mgmtorder-service running
tests-order-http:
	go test -v -test.count=1  -test.run=TestOrderReqSuite/OrderAPIReq ./tests/order/...

clean:
	rm -rf bin
	lsof -ti tcp:8080