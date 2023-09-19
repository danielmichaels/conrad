package providers

import (
	"crypto/tls"
	"github.com/xanzy/go-gitlab"
	"net/http"
)

type Gitlab struct {
	Client *gitlab.Client
}

func NewGitlab(token string, url string, insecure bool) (*Gitlab, error) {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}}
	hc := &http.Client{Transport: tr}
	cl, err := gitlab.NewClient(token, gitlab.WithBaseURL(url), gitlab.WithHTTPClient(hc))
	if err != nil {
		return nil, err
	}
	return &Gitlab{Client: cl}, nil
}
