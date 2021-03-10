build-local:
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
	@echo "> Done compiling!"
	@echo

run:
	CGO_ENABLED=0 go run main.go	

test:
	CGO_ENABLED=0 go test -v -cover ./...

test-coverage:
	CGO_ENABLED=0 go test -v -cover -coverprofile=coverage.txt ./...

# Builds the Golang executables and creates a Pip package out of it.
# This is mainly meant for locally testing whether the generated package is any good.
# The real packages for distribution on PyPI are made on the CI by cibuildwheel, see '.github/workflows/build-publish.yml'
package: build-all pre-package package-pip

# Cleans the build folder and copies ReadMe into build dir. Used on CI.
pre-package:
	@echo "> Cleaning build folder"
	rm -r build/build build/dist build/mllint.egg-info || true
	@echo
	@echo "> Copying ReadMe.md"
	cp ReadMe.md build/
	@echo

package-pip:
	@echo "> Creating 'mllint' Python package..."
	pip install build
	cd build && python -m build
	@echo
	@echo "> Done! Find your package and wheels in  ./build/dist "
	@echo

publish:
	@echo "> Publishing Python package..."
	@echo "Not yet implemented"
	@echo

	