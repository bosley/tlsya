package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/bosley/tlsya"
)

func main() {
	var config tlsya.TLSConfig
	var ipAddresses string

	flag.StringVar(&ipAddresses, "ips", "", "Comma-separated list of IP addresses")
	flag.StringVar(&config.KeyPath, "key", "key.pem", "Path to save the private key")
	flag.StringVar(&config.CertPath, "cert", "cert.pem", "Path to save the certificate")
	flag.Parse()

	if ipAddresses == "" {
		log.Fatal("At least one IP address is required")
	}

	config.IPAddresses = strings.Split(ipAddresses, ",")

	err := tlsya.GenerateTLS(config)
	if err != nil {
		log.Fatalf("Failed to generate TLS certificate: %v", err)
	}

	fmt.Println("TLS certificate and key generated successfully")
	fmt.Printf("Private key saved to: %s\n", config.KeyPath)
	fmt.Printf("Certificate saved to: %s\n", config.CertPath)
}
