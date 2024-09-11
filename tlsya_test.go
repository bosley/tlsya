package tlsya

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGenerateTLS(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "tlsya_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	keyPath := filepath.Join(tempDir, "test.key")
	certPath := filepath.Join(tempDir, "test.crt")

	config := TLSConfig{
		IPAddresses: []string{"127.0.0.1", "192.168.1.1"},
		KeyPath:     keyPath,
		CertPath:    certPath,
	}

	// Test GenerateTLS function
	err = GenerateTLS(config)
	if err != nil {
		t.Fatalf("GenerateTLS failed: %v", err)
	}

	// Check if files were created
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		t.Errorf("Key file was not created at %s", keyPath)
	}

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		t.Errorf("Certificate file was not created at %s", certPath)
	}

	// Read and parse the certificate
	certPEM, err := ioutil.ReadFile(certPath)
	if err != nil {
		t.Fatalf("Failed to read certificate file: %v", err)
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		t.Fatalf("Failed to parse certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse certificate: %v", err)
	}

	// Check IP addresses
	if len(cert.IPAddresses) != 2 {
		t.Errorf("Expected 2 IP addresses, got %d", len(cert.IPAddresses))
	}

	expectedIPs := map[string]bool{
		"127.0.0.1":   true,
		"192.168.1.1": true,
	}

	for _, ip := range cert.IPAddresses {
		if !expectedIPs[ip.String()] {
			t.Errorf("Unexpected IP address in certificate: %s", ip.String())
		}
		delete(expectedIPs, ip.String())
	}

	if len(expectedIPs) > 0 {
		t.Errorf("Not all expected IP addresses were in the certificate")
	}

	// Check validity period
	expectedValidity := 365 * 24 * time.Hour
	actualValidity := cert.NotAfter.Sub(cert.NotBefore)
	if actualValidity != expectedValidity {
		t.Errorf("Expected validity period of %v, got %v", expectedValidity, actualValidity)
	}
}
