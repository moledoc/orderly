
all: clean dir compile

dir: clean
	mkdir -p bin

compile:
	go build -o ./bin/mgmtuser ./cmd/mgmtuser/main.go
	go build -o ./bin/mgmtorder ./cmd/mgmtorder/main.go
	go build -o ./bin/orderly ./cmd/orderly/main.go

run-user:
	go run ./cmd/mgmtuser/main.go

run-order:
	go run ./cmd/mgmtorder/main.go

run:
	go run ./cmd/orderly/main.go

tests-user:
	echo "TODO: test mgmtuser"

tests-order:
	echo "TODO: test mgmtorder"

tests:
	echo "TODO: test orderly"

clean:
	rm -rf bin