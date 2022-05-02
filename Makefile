.PHONY: default test test-cover dev generate hooks lint-web doc

# for dev
dev:
	air -c .air.toml	

lint:
	golangci-lint run

build:
	go build -o image-pipeline-server 