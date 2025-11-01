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
	rm -rf output/*
co:
	rm -rf output/obstacleLUT-*.json
cn:
	rm -rf output/nodeLUT-*.json
cl:
	find . -maxdepth 2 -name "*.json" -print0 | while IFS= read -r -d $$'\0' file; \
	do size_bytes=$$(wc -c < "$$file") size_mb=$$(awk "BEGIN {printf \"%.2f\", $$size_bytes / (1024*1024)}"); \
 	echo "$$size_mb MB $$file"; \
	done | sort -h
u:
	go get -u ./...; \
	go mod tidy
