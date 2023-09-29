package providers

import (
	"crypto/tls"
	"github.com/xanzy/go-gitlab"
	"net/http"
	"time"
)

const GitlabClientDefaultTimeout = 10 * time.Second

type Gitlab struct {
	Client *gitlab.Client
}

func NewGitlab(token string, url string, insecure bool, timeout time.Duration) (*Gitlab, error) {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}}
	hc := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
	cl, err := gitlab.NewClient(
		token,
		gitlab.WithBaseURL(url),
		gitlab.WithHTTPClient(hc),
		gitlab.WithoutRetries(),
	)
	if err != nil {
		return nil, err
	}
	return &Gitlab{Client: cl}, err
}
