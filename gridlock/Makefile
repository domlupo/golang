lines:
	find . -type f -name "*.go" -exec wc -l {} + | sort -nr
cpu:
	go tool pprof -http=:6060 output/cpu.prof
mem:
	go tool pprof -http=:6060 output/mem.prof
cover:
	go test -v -cover -coverprofile=cover.out; \
        go tool cover -html=cover.out
c:
	rm -rf output
co:
	rm -rf output/obstacleLUT.json
cn:
	rm -rf output/nodeLUT.json
