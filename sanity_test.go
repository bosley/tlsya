package tlsya

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestTLSYAWithGin(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "tlsya_sanity_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	keyPath := filepath.Join(tempDir, "server.key")
	certPath := filepath.Join(tempDir, "server.crt")

	config := TLSConfig{
		IPAddresses: []string{"127.0.0.1"},
		KeyPath:     keyPath,
		CertPath:    certPath,
	}

	// Generate TLS certificate
	err = GenerateTLS(config)
	if err != nil {
		t.Fatalf("Failed to generate TLS certificate: %v", err)
	}

	// Set up Gin server
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Start the server in a goroutine
	go func() {
		err := r.RunTLS(":8443", config.CertPath, config.KeyPath)
		if err != nil {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Wait for the server to start
	time.Sleep(2 * time.Second)

	// Create a custom HTTP client that skips certificate verification
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: customTransport}

	// Make a request to the server
	resp, err := client.Get("https://localhost:8443/ping")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expected := `{"message":"pong"}`
	if string(body) != expected {
		t.Errorf("Expected body %q, got %q", expected, string(body))
	}

	// Verify that the connection was over HTTPS
	if resp.TLS == nil {
		t.Error("Expected HTTPS connection, but TLS was nil")
	}

	fmt.Println("HTTPS server is working correctly!")
}
