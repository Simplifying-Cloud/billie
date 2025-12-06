.PHONY: dev build templ css run clean

# Development with live reload (run all in parallel)
dev:
	@echo "Starting development servers..."
	@make -j3 templ-watch css-watch run-watch

# Watch templ files
templ-watch:
	~/go/bin/templ generate --watch --proxy="http://localhost:8080" --open-browser=false

# Watch CSS changes
css-watch:
	./tailwindcss -i static/css/input.css -o static/css/output.css --watch

# Run server with air or go run
run-watch:
	go run ./cmd/server

# One-time generation
templ:
	~/go/bin/templ generate

# One-time CSS build
css:
	./tailwindcss -i static/css/input.css -o static/css/output.css

# Production build
build: templ
	./tailwindcss -i static/css/input.css -o static/css/output.css --minify
	go build -o bin/server ./cmd/server

# Run the server (after building)
run: templ css
	go run ./cmd/server

# Clean generated files
clean:
	rm -rf bin/
	rm -f static/css/output.css
	find . -name "*_templ.go" -delete
