// Package githubhook implements handling and verification of github webhooks.
package webhook

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
	DeploymentEvent DeploymentEvent
	// Signature specifies the signature of a github webhook request.
	Signature string
	// Payload contains the raw contents of the webhook request.
	Payload []byte
}

// NewHook returns a Hook from an incoming HTTP Request.
func (s *Service) NewHook(req *http.Request) (hook Hook, err error) {
	if !strings.EqualFold(req.Method, "POST") {
		return Hook{}, microerror.Maskf(executionFailedError, fmt.Sprintf("%#q requests are not supported", req.Method))
	}

	if hook.Signature = req.Header.Get("x-hub-signature"); len(hook.Signature) == 0 {
		return Hook{}, microerror.Maskf(executionFailedError, "no signature found")
	}

	if hook.ID = req.Header.Get("x-github-delivery"); len(hook.ID) == 0 {
		return Hook{}, microerror.Maskf(executionFailedError, "no event id found")
	}

	if signedBy(hook, s.webhookSecret) {
		return Hook{}, microerror.Maskf(executionFailedError, "invalid signature found")
	}

	hook.Payload, err = ioutil.ReadAll(req.Body)
	if err != nil {
		return Hook{}, microerror.Mask(err)
	}

	err = json.Unmarshal(hook.Payload, &hook.DeploymentEvent)
	if err != nil {
		return Hook{}, microerror.Mask(err)
	}

	return hook, nil
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
