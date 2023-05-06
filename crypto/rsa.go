package crypto

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/tls"
    "crypto/x509"
    "encoding/pem"
    "math/big"
)

type (
    rsx  struct{}
    iRsx interface {
        ToTls() *tls.Config
    }
)

func Irsx() iRsx {
    return &rsx{}
}

//
// generate
//  @Description:
//  @return keyPEM
//  @return certPEM
//
func generate() (keyPEM, certPEM []byte) {

    key, err := rsa.GenerateKey(rand.Reader, 1024)
    if err != nil {
        panic(err)
    }
    template := x509.Certificate{SerialNumber: big.NewInt(1)}
    certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
    if err != nil {
        panic(err)
    }
    keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
    certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

    return keyPEM, certPEM

}

//
// ToTls
//  @Description:
//  @receiver this
//  @return *tls.Config
//
func (this rsx) ToTls() *tls.Config {

    keyPEM, certPEM := generate()
    tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
    if err != nil {
        panic(err)
    }
    return &tls.Config{
        Certificates: []tls.Certificate{tlsCert},
        NextProtos:   []string{"mixin-quic-peer"},
    }

}
