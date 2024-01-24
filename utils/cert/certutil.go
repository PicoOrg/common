package cert

type Cert interface {
	GenCA(commonName string) (caCertPem, caKeyPem []byte, err error)
	GenHttpCert(caCertPem, caKeyPem []byte, commonName string) (childCertPem, childKeyPem []byte, err error)
}
