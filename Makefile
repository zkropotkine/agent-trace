.PHONY: test coverage cover-html clean

# Run tests
test:
	go test ./... -v

# Run tests with coverage and output a report
coverage:
	go test ./... -coverprofile=coverage.out

# Open the coverage report in your browser
cover-html: coverage
	go tool cover -html=coverage.out

# Clean generated files
clean:
	rm -f coverage.out