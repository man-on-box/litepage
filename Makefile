dev:
	air

serve:
	@echo "Serving site on http://localhost:3001"
	@go run cmd/litepage/main.go -dev

build:
	@echo "Building static site"
	go run cmd/litepage/main.go

devtools:
	@echo "Installing recommended devtools"
	@go install github.com/air-verse/air@latest
	@echo "Devtools installed"
