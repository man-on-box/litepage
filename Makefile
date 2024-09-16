dev:
	air

build:
	@echo "Building static site"
	go run cmd/litepage/main.go

devtools:
	@echo "Installing recommended devtools"
	@go install github.com/air-verse/air@latest
	@echo "Devtools installed"
