package mns

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"sort"
	"strings"
	"time"
)

// Credential mns credential
type Credential interface {
	Sign() (string, error)
	SetAccessSecret(string)
	SetHeader(http.Header)
	SetResource(string)
	SetMethod(string)
}

type credential struct {
	accessKeySecret string
	h               http.Header
	resource        string
	method          string
}

// NewCredential new credential
func NewCredential(accessSecret string) Credential {
	c := new(credential)
	c.accessKeySecret = accessSecret
	return c
}

// Sign calc sign
func (c *credential) Sign() (string, error) {
	var date = time.Now().UTC().Format(http.TimeFormat)
	if v := c.h.Get("Date"); v != "" {
		date = v
	}

	var contentMD5 = c.h.Get("Content-MD5")
	var contentType = c.h.Get("Content-Type")

	var mnsHeaders []string
	for k := range c.h {
		if strings.HasPrefix(strings.ToLower(k), "x-mns-") {
			mnsHeaders = append(mnsHeaders, strings.ToLower(k)+":"+strings.TrimSpace(c.h.Get(k)))
		}
	}

	sort.Sort(sort.StringSlice(mnsHeaders))

	stringToSign := c.method + "\n" +
		contentMD5 + "\n" +
		contentType + "\n" +
		date + "\n" +
		strings.Join(mnsHeaders, "\n") + "\n" +
		"/" + c.resource

	hmaced := hmac.New(sha1.New, []byte(c.accessKeySecret))
	if _, err := hmaced.Write([]byte(stringToSign)); err != nil {
		return "", err
	}

	s := base64.StdEncoding.EncodeToString(hmaced.Sum(nil))
	return s, nil
}

// SetAccessSecret set mns access secret
func (c *credential) SetAccessSecret(accessSecret string) {
	c.accessKeySecret = accessSecret
}

// SetHeader set header
func (c *credential) SetHeader(h http.Header) {
	c.h = h
}

// SetResource set resource
func (c *credential) SetResource(resource string) {
	c.resource = resource
}

// SetMethod set method
func (c *credential) SetMethod(method string) {
	c.method = method
}
