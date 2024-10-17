dev:
	@LP_MODE=serve LP_PORT=3001 air

serve:
	@LP_MODE=serve go run ./... -lp-serve

build:
	@go run ./...

install-devtools:
	@echo "Installing recommended devtools"
	@go install github.com/air-verse/air@latest
	@echo "Devtools installed"
