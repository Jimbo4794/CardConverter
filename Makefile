all: select convert

select:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/select-x86_64 cmd/select/select.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/select-x86_64.exe cmd/select/select.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/select-darwin-x86_64 cmd/select/select.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/select-darwin-arm64 cmd/select/select.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/select-linux-arm64 cmd/select/select.go

convert:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/convert-x86_64 cmd/convert/convert.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/convert-x86_64.exe cmd/convert/convert.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/convert-darwin-x86_64 cmd/convert/convert.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/convert-darwin-arm64 cmd/convert/convert.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/convert-linux-arm64 cmd/convert/convert.go
	

clean:
	rm -fr build/*