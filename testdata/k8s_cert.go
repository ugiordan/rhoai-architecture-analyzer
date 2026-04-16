package testdata

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

func generateServerCert() *x509.Certificate {
	serialNumber, _ := rand.Int(rand.Reader, big.NewInt(time.Now().UnixNano()))
	return &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "test-server",
		},
		IsCA:     true,
		KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		DNSNames: []string{"localhost"},
	}
}
