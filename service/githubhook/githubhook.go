// Package githubhook implements handling and verification of github webhooks.
package githubhook

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

	// Id specifies the Id of a github webhook request.
	Id string
	// Event contains unmarshaled webhook.
	Event Event
	// Event specifies the event name of a github webhook request.
	EventName string
	// Signature specifies the signature of a github webhook request.
	Signature string
	// Payload contains the raw contents of the webhook request.
	Payload []byte
}

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}

// SignedBy checks that the provided secret matches the hook Signature.
//
// Implements validation described in github's documentation:
// https://developer.github.com/webhooks/securing/
func (h *Hook) SignedBy(secret []byte) bool {
	if len(h.Signature) != signatureLength || !strings.HasPrefix(h.Signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(h.Signature[5:]))

	return hmac.Equal(signBody(secret, h.Payload), actual)
}

// New reads a Hook from an incoming HTTP Request.
func New(req *http.Request) (hook *Hook, err error) {
	hook = new(Hook)
	if !strings.EqualFold(req.Method, "POST") {
		return nil, microerror.Maskf(invalidHookError, fmt.Sprintf("%#q requests are not supported", req.Method))
	}

	if hook.Signature = req.Header.Get("x-hub-signature"); len(hook.Signature) == 0 {
		return nil, microerror.Maskf(invalidHookError, "no signature found")
	}

	if hook.EventName = req.Header.Get("x-github-event"); len(hook.EventName) == 0 {
		return nil, microerror.Maskf(invalidHookError, "no event found")
	}

	if hook.Id = req.Header.Get("x-github-delivery"); len(hook.Id) == 0 {
		return nil, microerror.Maskf(invalidHookError, "no event id found")
	}

	hook.Payload, err = ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	err = json.Unmarshal(hook.Payload, &hook.Event)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return hook, nil
}

// Parse reads and verifies the hook in an inbound request.
func Parse(secret []byte, req *http.Request) (hook *Hook, err error) {
	hook, err = New(req)
	if err != nil {
		return hook, microerror.Mask(err)
	}
	if !hook.SignedBy(secret) {
		err = microerror.Maskf(invalidHookError, "invalid signature found")
	}
	return
}
