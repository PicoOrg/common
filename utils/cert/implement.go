package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

func New() Cert {
	return &implement{}
}

type implement struct {
}

func (m *implement) GenCA(commonName string) (caCertPem, caKeyPem []byte, err error) {
	var serialNumber *big.Int
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	if serialNumber, err = rand.Int(rand.Reader, serialNumberLimit); err != nil {
		return
	}

	caCert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:         commonName,
			Organization:       []string{commonName},
			OrganizationalUnit: []string{"ca"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, 800),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	var caKey *rsa.PrivateKey
	if caKey, err = rsa.GenerateKey(rand.Reader, 4096); err != nil {
		return
	}
	return m.GenCert(caCert, caKey, caCert, caKey)
}

func (m *implement) GenHttpCert(caCertPem, caKeyPem []byte, commonName string) (childCertPem, childKeyPem []byte, err error) {
	var caCert *x509.Certificate
	if caCert, err = m.CertFromPem(caCertPem); err != nil {
		return
	}

	var caKey *rsa.PrivateKey
	if caKey, err = m.KeyFromPem(caKeyPem); err != nil {
		return
	}

	var serialNumber *big.Int
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	if serialNumber, err = rand.Int(rand.Reader, serialNumberLimit); err != nil {
		return
	}

	childCert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:         fmt.Sprintf("*.%s", commonName),
			Organization:       []string{commonName},
			OrganizationalUnit: []string{"server"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 0, 800),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		DNSNames:     []string{commonName, fmt.Sprintf("*.%s", commonName)},
	}

	var childKey *rsa.PrivateKey
	if childKey, err = rsa.GenerateKey(rand.Reader, 4096); err != nil {
		return
	}
	return m.GenCert(caCert, caKey, childCert, childKey)
}

func (m *implement) GenCert(
	caCert *x509.Certificate, caKey *rsa.PrivateKey,
	childCert *x509.Certificate, childKey *rsa.PrivateKey) (childCertPem, childKeyPem []byte, err error) {

	var childCertBytes []byte
	if childCertBytes, err = x509.CreateCertificate(rand.Reader, childCert, caCert, &childKey.PublicKey, caKey); err != nil {
		return
	}

	childCertPemBuffer := new(bytes.Buffer)
	if err = pem.Encode(childCertPemBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: childCertBytes}); err != nil {
		return
	}
	childCertPem = childCertPemBuffer.Bytes()

	childKeyPemBuffer := new(bytes.Buffer)
	if err = pem.Encode(childKeyPemBuffer, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(childKey)}); err != nil {
		return
	}
	childKeyPem = childKeyPemBuffer.Bytes()

	return
}

func (m *implement) CertToPem(cert *x509.Certificate) (certPem []byte, err error) {
	certPemBuffer := new(bytes.Buffer)
	if err = pem.Encode(certPemBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}); err != nil {
		return
	}
	certPem = certPemBuffer.Bytes()
	return
}

func (m *implement) CertFromPem(certPem []byte) (cert *x509.Certificate, err error) {
	certBlock, _ := pem.Decode(certPem)
	if cert, err = x509.ParseCertificate(certBlock.Bytes); err != nil {
		return
	}
	return
}

func (m *implement) KeyToPem(key *rsa.PrivateKey) (keyPem []byte, err error) {
	keyPemBuffer := new(bytes.Buffer)
	if err = pem.Encode(keyPemBuffer, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
		return
	}
	keyPem = keyPemBuffer.Bytes()
	return
}

func (m *implement) KeyFromPem(keyPem []byte) (key *rsa.PrivateKey, err error) {
	keyBlock, _ := pem.Decode(keyPem)
	if key, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes); err != nil {
		return
	}
	return
}
