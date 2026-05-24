# Example for http3 transporter

## Install all dependencies
```
go mod tidy
```

## Install msgp if not already installed for code generation
```
go install github.com/tinylib/msgp@latest
```

## Use code generation for Response and Request structs
```
go generate ./...
```

## Start server
```
go run ./...
```