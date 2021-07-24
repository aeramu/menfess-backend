generate:
	go mod init github.com/aeramu/menfess-backend
	gocto generate
	mockery --all --dir service
	go mod tidy
mock:
	mockery --all --dir service
test:
	go test ./... --cover
