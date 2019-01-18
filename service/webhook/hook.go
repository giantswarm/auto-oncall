// Package githubhook implements handling and verification of github webhooks.
package github

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/giantswarm/microerror"
)

const (
	signaturePrefix = "sha1="
	signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))
)

// Hook is an inbound github webhook.
type Hook struct {
	// ID specifies the ID of a github webhook request.
	ID string
	// Event contains unmarshaled webhook.
	Event Event
	// Event specifies the event name of a github webhook request.
	EventName string
	// Signature specifies the signature of a github webhook request.
	Signature string
	// Payload contains the raw contents of the webhook request.
	Payload []byte
}

func signBody(body, secret []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}

// signedBy checks that the provided secret matches the hook Signature.
//
// Implements validation described in github's documentation:
// https://developer.github.com/webhooks/securing/
func signedBy(h Hook, secret []byte) bool {
	if len(h.Signature) != signatureLength || !strings.HasPrefix(h.Signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(h.Signature[5:]))

	return hmac.Equal(signBody(h.Payload, secret), actual)
}
