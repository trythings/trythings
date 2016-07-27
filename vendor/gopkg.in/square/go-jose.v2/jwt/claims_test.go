/*-
 * Copyright 2016 Zbigniew Mandziejewicz
 * Copyright 2016 Square, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/square/go-jose.v2"
)

type testClaims struct {
	Subject string `json:"sub"`
}

func TestCustomClaimsNonPointer(t *testing.T) {
	signingKey := jose.SigningKey{Algorithm: jose.RS256, Key: testPrivRSAKey1}
	signer, err := jose.NewSigner(signingKey, nil)
	if err != nil {
		t.Fatalf("error generating token: %v", err)
	}
	jwt, err := Signed(signer).Claims(testClaims{"foo"}).CompactSerialize()
	if err != nil {
		t.Fatalf("error creating jwt: %v", err)
	}
	parsed, err := ParseSigned(jwt)
	if err != nil {
		t.Fatalf("error parsing jwt: %v", err)
	}

	out := &testClaims{}
	if assert.NoError(t, parsed.Claims(out, &testPrivRSAKey1.PublicKey)) {
		assert.Equal(t, out.Subject, "foo")
	}
}

func TestCustomClaimsPointer(t *testing.T) {
	signingKey := jose.SigningKey{Algorithm: jose.RS256, Key: testPrivRSAKey1}
	signer, err := jose.NewSigner(signingKey, nil)
	if err != nil {
		t.Fatalf("error generating token: %v", err)
	}
	jwt, err := Signed(signer).Claims(&testClaims{"foo"}).CompactSerialize()
	if err != nil {
		t.Fatalf("error creating jwt: %v", err)
	}
	parsed, err := ParseSigned(jwt)
	if err != nil {
		t.Fatalf("error parsing jwt: %v", err)
	}

	out := &testClaims{}
	if assert.NoError(t, parsed.Claims(out, &testPrivRSAKey1.PublicKey)) {
		assert.Equal(t, out.Subject, "foo")
	}
}

func TestEncodeClaims(t *testing.T) {
	now := time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)

	c := Claims{
		Issuer:   "issuer",
		Subject:  "subject",
		Audience: []string{"a1", "a2"},
		IssuedAt: now,
		Expiry:   now.Add(1 * time.Hour),
	}

	b, err := c.marshalJSON()
	assert.NoError(t, err)

	expected := `{"iss":"issuer","sub":"subject","aud":["a1","a2"],"exp":1451610000,"iat":1451606400}`
	assert.Equal(t, expected, string(b))
}

func TestDecodeClaims(t *testing.T) {
	s := []byte(`{"iss":"issuer","sub":"subject","aud":["a1","a2"],"exp":1451610000,"iat":1451606400}`)
	now := time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)

	c := Claims{}
	err := c.unmarshalJSON(s)
	assert.NoError(t, err)

	assert.Equal(t, "issuer", c.Issuer)
	assert.Equal(t, "subject", c.Subject)
	assert.Equal(t, []string{"a1", "a2"}, c.Audience)
	assert.True(t, now.Equal(c.IssuedAt))
	assert.True(t, now.Add(1*time.Hour).Equal(c.Expiry))
}
