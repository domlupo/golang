lines:
	find . -type f -name "*.go" -exec wc -l {} +
cpu:
	go tool pprof -http=:6060 cpu.prof
mem:
	go tool pprof -http=:6060 mem.prof
cover:
	go test -v -cover -coverprofile=cover.out; \
        go tool cover -html=cover.out
