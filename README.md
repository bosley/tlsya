# tlsya - TLS Certificate Generator for Go

`tlsya` is a simple Go library for generating self-signed TLS certificates.
It's designed to be easy to use while also providing flexibility for more advanced use cases.

## Features

- Generate self-signed TLS certificates with a single function call
- Customize certificate details using your own x509.Certificate template
- Specify IP addresses to be included in the certificate
- Easy integration with Go web servers, including Gin
- Command-line interface (CLI) for quick certificate generation

## Installation

To use tlsya as a library in your Go project:

```
go get github.com/bosley/tlsya
```

To install the CLI tool:

```
go install github.com/bosley/tlsya/cmd/cli@latest
```

## Usage

### CLI Usage

After installing the CLI tool, you can generate TLS certificates from the command line:

```
tlsya-cli -ips=127.0.0.1,192.168.1.1 -key=server.key -cert=server.crt
```

Options:
- `-ips`: Comma-separated list of IP addresses (required)
- `-key`: Path to save the private key (default: "key.pem")
- `-cert`: Path to save the certificate (default: "cert.pem")

### Library Usage

#### Basic Usage

```go
import "github.com/bosley/tlsya"

config := tlsya.TLSConfig{
    IPAddresses: []string{"127.0.0.1", "192.168.1.1"},
    KeyPath:     "server.key",
    CertPath:    "server.crt",
}

err := tlsya.GenerateTLS(config)
if err != nil {
    // Handle error
}
```

#### Advanced Usage

```go
import (
    "crypto/x509"
    "crypto/x509/pkix"
    "math/big"
    "time"
    "github.com/bosley/tlsya"
)

config := tlsya.TLSConfig{
    IPAddresses: []string{"127.0.0.1", "192.168.1.1"},
    KeyPath:     "server.key",
    CertPath:    "server.crt",
}

template := &x509.Certificate{
    SerialNumber: big.NewInt(1),
    Subject: pkix.Name{
        Organization: []string{"My Company"},
        Country:      []string{"US"},
    },
    NotBefore:             time.Now(),
    NotAfter:              time.Now().AddDate(1, 0, 0),
    KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
    ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
    BasicConstraintsValid: true,
}

err := tlsya.GenerateTLSFrom(config, template)
if err != nil {
    // Handle error
}
```

### Using with Gin

Here's an example of how to use `tlsya` with a Gin web server:

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/bosley/tlsya"
)

func main() {
    config := tlsya.TLSConfig{
        IPAddresses: []string{"127.0.0.1"},
        KeyPath:     "server.key",
        CertPath:    "server.crt",
    }

    err := tlsya.GenerateTLS(config)
    if err != nil {
        // Handle error
    }

    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

    r.RunTLS(":8080", config.CertPath, config.KeyPath)
}
```

This will start a Gin server with HTTPS enabled on port 8080.
