package test_utils

import (
  "fmt"
  "strconv"
  _ "embed"
  "crypto/rsa"
  "crypto/x509"
  "encoding/json"
  "encoding/pem"
  "github.com/lestrrat-go/jwx/jwa"
  "github.com/lestrrat-go/jwx/jws"
  "github.com/lestrrat-go/jwx/jwt"
)

var (
  initialized bool
  rsaPrivateKey *rsa.PrivateKey
)

//go:embed private-key.pem
var privateKeyFile []byte

const (
	issuer = "authority"
	kid    = "aa7c6287-c45d-4966-84b4-a1633e4e3a64"
)

func InitKey() {
  block, _ := pem.Decode(privateKeyFile)
  key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
  if err != nil {
    panic(fmt.Sprintf("failed to parse private key: %s", err))
  }

  rsaPrivateKey = key
}

func CreateAccessToken(sub string) ([]byte, error) {
  token := jwt.New()

  if err := token.Set(jwt.IssuerKey, issuer); err != nil {
    return nil, fmt.Errorf("failed to set the issuer key to the token: %w", err)
  }

  if err := token.Set(jwt.SubjectKey, sub); err != nil {
    return nil, fmt.Errorf("failed to set the subject key to the token: %w", err)
  }

  headers := jws.NewHeaders()
  if err := headers.Set(jws.KeyIDKey, kid); err != nil {
    return nil, fmt.Errorf("failed to create jws headers: %w", err)
  }

  if err := headers.Set(jws.AlgorithmKey, jwa.RS256); err != nil {
    return nil, fmt.Errorf("failed to set the alg key to the token: %w", err)
  }

  if err := headers.Set(jws.TypeKey, "JWT"); err != nil {
    return nil, fmt.Errorf("failed to set the typ key to the token: %w", err)
  }

  b, err := json.Marshal(token)
  if err != nil {
    return nil, fmt.Errorf("failed to marshal the token: %w", err)
  }

  signedToken, err := jws.Sign(b, jwa.RS256, rsaPrivateKey, jws.WithHeaders(headers))
  if err != nil {
    return nil, fmt.Errorf("failed to sign the token: %w", err)
  }
  return signedToken, nil
}

func GetAccessToken( customer_id int ) string {
  if !initialized {
    InitKey()
    initialized = true
  }
  customer_id_in_str := strconv.Itoa( customer_id )
  token, e := CreateAccessToken( customer_id_in_str )
  if e != nil {
    panic( e )
  }
  return string(token)
}

