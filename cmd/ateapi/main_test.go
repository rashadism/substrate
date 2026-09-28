// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestConnectStoreRequiresPostgresConnectionString(t *testing.T) {
	oldDSN := *postgresConnectionString
	t.Cleanup(func() {
		*postgresConnectionString = oldDSN
	})
	*postgresConnectionString = ""

	_, err := connectStore(context.Background())
	if err == nil || !strings.Contains(err.Error(), "--postgres-connection-string is required") {
		t.Fatalf("connectStore() error = %v, want missing-connection-string error", err)
	}
}

func TestBuildServerTLSConfigWithoutCACertsAllowsCertlessClients(t *testing.T) {
	cfg, err := buildServerTLSConfig(context.Background(), "/nonexistent-cred-bundle.pem", "")
	if err != nil {
		t.Fatalf("buildServerTLSConfig() error = %v", err)
	}
	if cfg.GetConfigForClient != nil {
		t.Fatalf("buildServerTLSConfig() with no CA path set GetConfigForClient, want nil (no client-cert verification configured)")
	}
}

func TestBuildServerTLSConfigRejectsUnreadableCACerts(t *testing.T) {
	_, err := buildServerTLSConfig(context.Background(), "/nonexistent-cred-bundle.pem", filepath.Join(t.TempDir(), "absent.pem"))
	if err == nil {
		t.Fatalf("buildServerTLSConfig() error = nil, want an error for a missing CA file")
	}
}

// TestBuildServerTLSConfigReloadsCACertsWithoutRestart proves the specific
// regression this fix closes: a pod-identity CA rotation on disk must be
// picked up by the next handshake, not frozen at the config's construction.
func TestBuildServerTLSConfigReloadsCACertsWithoutRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "trust-bundle.pem")
	writeCA(t, path, "ca-one")

	cfg, err := buildServerTLSConfig(context.Background(), "/nonexistent-cred-bundle.pem", path)
	if err != nil {
		t.Fatalf("buildServerTLSConfig() error = %v", err)
	}
	if cfg.GetConfigForClient == nil {
		t.Fatalf("buildServerTLSConfig() with a CA path did not set GetConfigForClient")
	}

	before, err := cfg.GetConfigForClient(nil)
	if err != nil {
		t.Fatalf("GetConfigForClient() first call error = %v", err)
	}

	writeCA(t, path, "ca-two")

	after, err := cfg.GetConfigForClient(nil)
	if err != nil {
		t.Fatalf("GetConfigForClient() second call error = %v", err)
	}

	if before.ClientCAs.Equal(after.ClientCAs) {
		t.Fatalf("GetConfigForClient() returned the same trust pool after the CA file changed, want the rotated one")
	}
}

// writeCA writes a fresh self-signed certificate (distinguished by cn) to
// path, suitable for AppendCertsFromPEM.
func writeCA(t *testing.T, path, cn string) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: cn},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		IsCA:         true,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("CreateCertificate() error = %v", err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	if err := os.WriteFile(path, pemBytes, 0o600); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", path, err)
	}
}
