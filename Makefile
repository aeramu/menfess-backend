generate:
	go mod init github.com/aeramu/menfess-backend
	gocto generate
	mockery --all
	go mod tidy
mock:
	mockery --all
test:
	go test ./... --cover
