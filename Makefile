build:
	CGO_ENABLED=0 go build -o build/bin/mllint main.go

build-all:
	@echo "> Compiling all executables..."
	@echo
	# 32-bit
	GOOS=freebsd GOARCH=386 CGO_ENABLED=0 go build -o build/bin/mllint-freebsd-386 main.go
	# GOOS=darwin GOARCH=386 CGO_ENABLED=0 go build -o build/bin/mllint-darwin-386 main.go
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o build/bin/mllint-linux-386 main.go
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -o build/bin/mllint-windows-386 main.go

	# 64-bit
	GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 go build -o build/bin/mllint-freebsd-amd64 main.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o build/bin/mllint-darwin-amd64 main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/bin/mllint-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o build/bin/mllint-windows-amd64 main.go
	@echo
	@echo "> Done!"
	@echo

run:
	CGO_ENABLED=0 go run main.go	

test:
	CGO_ENABLED=0 go test -v -cover ./...

test-coverage:
	CGO_ENABLED=0 go test -v -cover -coverprofile=coverage.txt ./...

# TODO: make a command to create an mllint Python package.
# - Build executable
# - Create platform-specific Python wheels that include the binary.
# - Publish to PyPI
package: build-all package-pip publish

package-pip:
	@echo "> Creating 'mllint' Python package..."
	cd build
	rm -r build dist mllint.egg-info/
	python -m build
	@echo

publish:
	@echo "> Publishing Python package wheels..."
	@echo "Not yet implemented"
	@echo

	