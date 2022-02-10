// The MIT License
//
// Copyright (c) 2022 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package rpc

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/temporalio/ui-server/server/config"
)

var _ CertProvider = (*localStoreCertProvider)(nil)
var _ CertExpirationChecker = (*localStoreCertProvider)(nil)

type certCache struct {
	cert    *tls.Certificate
	cAPool  *x509.CertPool
	cACerts []*x509.Certificate // copies of certs in the cAPool CertPool for expiration checks
}

type localStoreCertProvider struct {
	sync.RWMutex

	tlsSettings *config.TLS

	certs           *certCache
	refreshInterval time.Duration

	ticker *time.Ticker
	stop   chan bool
}

type loadOrDecodeDataFunc func(item string) ([]byte, error)

type tlsCertFetcher func() (*tls.Certificate, error)

func (s *localStoreCertProvider) initialize() {

	if s.refreshInterval != 0 {
		s.stop = make(chan bool)
		s.ticker = time.NewTicker(s.refreshInterval)
		go s.refreshCerts()
	}
}

func NewLocalStoreCertProvider(
	tlsSettings *config.TLS,
	refreshInterval time.Duration) CertProvider {

	provider := &localStoreCertProvider{
		tlsSettings:     tlsSettings,
		refreshInterval: refreshInterval,
	}
	provider.initialize()
	return provider
}

func (s *localStoreCertProvider) Close() {

	if s.ticker != nil {
		s.ticker.Stop()
	}
	if s.stop != nil {
		s.stop <- true
		close(s.stop)
	}
}

func (s *localStoreCertProvider) FetchCA() (*x509.CertPool, error) {
	if s.tlsSettings == nil {
		return nil, nil
	}
	certs, err := s.getCerts()
	if err != nil {
		return nil, err
	}
	return certs.cAPool, nil
}

func (s *localStoreCertProvider) FetchCertificate() (*tls.Certificate, error) {
	if !s.isTLSEnabled() {
		return nil, nil
	}

	certs, err := s.getCerts()
	if err != nil {
		return nil, err
	}
	return certs.cert, nil
}

func (s *localStoreCertProvider) GetExpiringCerts(timeWindow time.Duration,
) (CertExpirationMap, CertExpirationMap, error) {

	expiring := make(CertExpirationMap)
	expired := make(CertExpirationMap)
	when := time.Now().UTC().Add(timeWindow)

	certs, err := s.getCerts()
	if err != nil {
		return nil, nil, err
	}

	checkError := checkTLSCertForExpiration(certs.cert, when, expiring, expired)
	err = appendError(err, checkError)

	checkCertsForExpiration(certs.cACerts, when, expiring, expired)

	return expiring, expired, err
}

func (s *localStoreCertProvider) getCerts() (*certCache, error) {
	s.RLock()
	if s.certs != nil {
		defer s.RUnlock()
		return s.certs, nil
	}
	s.RUnlock()
	s.Lock()
	defer s.Unlock()

	if s.certs != nil {
		return s.certs, nil
	}

	newCerts, err := s.loadCerts()
	if err != nil {
		return nil, err
	}

	if newCerts == nil {
		s.certs = &certCache{}
	} else {
		s.certs = newCerts
	}
	return s.certs, nil
}

func (s *localStoreCertProvider) loadCerts() (*certCache, error) {
	if !s.isTLSEnabled() {
		return nil, nil
	}

	newCerts := certCache{}
	var err error

	if s.tlsSettings != nil {
		newCerts.cert, err = s.fetchCertificate(s.tlsSettings.CertFile, s.tlsSettings.CertData,
			s.tlsSettings.KeyFile, s.tlsSettings.KeyData)
		if err != nil {
			return nil, err
		}

		certPool, certs, err := s.fetchCA(s.tlsSettings.CaFile, s.tlsSettings.CaData,
			"cannot specify both clientCAFiles and clientCAData properties")
		if err != nil {
			return nil, err
		}
		newCerts.cAPool = certPool
		newCerts.cACerts = certs
	}

	return &newCerts, nil
}

func (s *localStoreCertProvider) fetchCertificate(
	certFile string, certData string,
	keyFile string, keyData string) (*tls.Certificate, error) {
	if certFile == "" && certData == "" {
		return nil, nil
	}

	if certFile != "" && certData != "" {
		return nil, errors.New("unable to specify both certFile and certData properties")
	}

	var certBytes []byte
	var keyBytes []byte
	var err error

	if certFile != "" {
		log.Printf("loading certificate from file: %v", certFile)
		certBytes, err = os.ReadFile(certFile)
		if err != nil {
			return nil, err
		}
	} else if certData != "" {
		certBytes, err = base64.StdEncoding.DecodeString(certData)
		if err != nil {
			return nil, fmt.Errorf("TLS public certificate could not be decoded: %w", err)
		}
	}

	if keyFile != "" {
		log.Printf("loading private key from file: %v", keyFile)
		keyBytes, err = os.ReadFile(keyFile)
		if err != nil {
			return nil, err
		}
	} else if keyData != "" {
		keyBytes, err = base64.StdEncoding.DecodeString(keyData)
		if err != nil {
			return nil, fmt.Errorf("TLS private key could not be decoded: %w", err)
		}
	}

	cert, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return nil, fmt.Errorf("loading tls certificate failed: %v", err)
	}

	return &cert, nil
}

func (s *localStoreCertProvider) fetchCA(
	file string,
	data string,
	duplicateErrorMessage string) (*x509.CertPool, []*x509.Certificate, error) {
	caPoolFromFiles, caCertsFromFiles, err := s.buildCAPoolFromFiles([]string{file})
	if err != nil {
		return nil, nil, err
	}

	caPoolFromData, caCertsFromData, err := buildCAPoolFromData([]string{data})
	if err != nil {
		return nil, nil, err
	}

	if caPoolFromFiles != nil && caPoolFromData != nil {
		return nil, nil, errors.New(duplicateErrorMessage)
	}

	var certPool *x509.CertPool
	var certs []*x509.Certificate

	if caPoolFromData != nil {
		certPool = caPoolFromData
		certs = caCertsFromData
	} else {
		certPool = caPoolFromFiles
		certs = caCertsFromFiles
	}

	return certPool, certs, nil
}

func checkTLSCertForExpiration(
	cert *tls.Certificate,
	when time.Time,
	expiring CertExpirationMap,
	expired CertExpirationMap,
) error {

	if cert == nil {
		return nil
	}

	x509cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return err
	}
	checkCertForExpiration(x509cert, when, expiring, expired)
	return nil
}

func checkCertsForExpiration(
	certs []*x509.Certificate,
	time time.Time,
	expiring CertExpirationMap,
	expired CertExpirationMap,
) {

	for _, cert := range certs {
		checkCertForExpiration(cert, time, expiring, expired)
	}
}

func checkCertForExpiration(
	cert *x509.Certificate,
	pointInTime time.Time,
	expiring CertExpirationMap,
	expired CertExpirationMap,
) {
	if cert != nil && expiresBefore(cert, pointInTime) {
		record := CertExpirationData{
			Thumbprint: md5.Sum(cert.Raw),
			IsCA:       cert.IsCA,
			DNSNames:   cert.DNSNames,
			Expiration: cert.NotAfter,
		}
		if record.Expiration.Before(time.Now().UTC()) {
			expired[record.Thumbprint] = record
		} else {
			expiring[record.Thumbprint] = record
		}
	}
}

func expiresBefore(cert *x509.Certificate, pointInTime time.Time) bool {
	return cert.NotAfter.Before(pointInTime)
}

func buildCAPoolFromData(caData []string) (*x509.CertPool, []*x509.Certificate, error) {

	return buildCAPool(caData, base64.StdEncoding.DecodeString)
}

func (s *localStoreCertProvider) buildCAPoolFromFiles(caFiles []string) (*x509.CertPool, []*x509.Certificate, error) {
	log.Printf("loading CA certs from: %v", caFiles)
	return buildCAPool(caFiles, os.ReadFile)
}

func buildCAPool(cas []string, getBytes loadOrDecodeDataFunc) (*x509.CertPool, []*x509.Certificate, error) {

	var caPool *x509.CertPool
	var certs []*x509.Certificate

	for _, ca := range cas {
		if ca == "" {
			continue
		}

		caBytes, err := getBytes(ca)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to decode ca cert: %w", err)
		}

		if caPool == nil {
			caPool = x509.NewCertPool()
		}
		if !caPool.AppendCertsFromPEM(caBytes) {
			return nil, nil, errors.New("unknown failure constructing cert pool for ca")
		}

		cert, err := parseCert(caBytes)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse x509 certificate: %w", err)
		}
		certs = append(certs, cert)
	}
	return caPool, certs, nil
}

// logic borrowed from tls.X509KeyPair()
func parseCert(bytes []byte) (*x509.Certificate, error) {

	var certBytes [][]byte
	for {
		var certDERBlock *pem.Block
		certDERBlock, bytes = pem.Decode(bytes)
		if certDERBlock == nil {
			break
		}
		if certDERBlock.Type == "CERTIFICATE" {
			certBytes = append(certBytes, certDERBlock.Bytes)
		}
	}

	if len(certBytes) == 0 || len(certBytes[0]) == 0 {
		return nil, fmt.Errorf("failed to decode PEM certificate data")
	}
	return x509.ParseCertificate(certBytes[0])
}

func appendError(aggregatedErr error, err error) error {
	if aggregatedErr == nil {
		return err
	}
	if err == nil {
		return aggregatedErr
	}
	return fmt.Errorf("%v, %w", aggregatedErr, err)
}

func (s *localStoreCertProvider) refreshCerts() {
	for {
		select {
		case <-s.stop:
			break
		case <-s.ticker.C:
		}

		newCerts, err := s.loadCerts()
		if err != nil {
			log.Printf("unable to load certificates: %s", err)
			continue
		}

		s.RLock()
		currentCerts := s.certs
		s.RUnlock()
		if currentCerts.isEqual(newCerts) {
			continue
		}

		log.Printf("loaded new TLS certificates")
		s.Lock()
		s.certs = newCerts
		s.Unlock()
	}
}

func (s *localStoreCertProvider) isTLSEnabled() bool {
	return s.tlsSettings != nil
}

func (c *certCache) isEqual(other *certCache) bool {
	if c == other {
		return true
	}
	if c == nil || other == nil {
		return false
	}

	if !equalTLSCerts(c.cert, other.cert) ||
		!equalX509(c.cACerts, other.cACerts) {
		return false
	}
	return true
}

func equal(a, b [][]byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !bytes.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

func equalX509(a, b []*x509.Certificate) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Equal(b[i]) {
			return false
		}
	}
	return true
}

func equalTLSCerts(a, b *tls.Certificate) bool {
	if a != nil {
		if b == nil || !equal(a.Certificate, b.Certificate) {
			return false
		}
	} else {
		if b != nil {
			return false
		}
	}
	return true
}
