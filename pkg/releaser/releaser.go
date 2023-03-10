package releaser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alex-held/devctl-release-bot/pkg/devctl"
	"github.com/alex-held/devctl-release-bot/pkg/source/actions"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

//Releaser is what opens PR
type Releaser struct {
	Token                           string
	TokenEmail                      string
	TokenUserHandle                 string
	TokenUsername                   string
	UpstreamDevctlIndexRepo         string
	UpstreamDevctlIndexRepoOwner    string
	UpstreamDevctlIndexRepoCloneURL string
	LocalDevctlIndexRepo            string
	LocalDevctlIndexRepoOwner       string
	LocalDevctlIndexRepoCloneURL    string
}

func getCloneURL(owner, repo string) string {
	return fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
}

//TODO: get email, userhandle, name from token
func getUserDetails(token string) (string, string, string) {
	return "devctl-release-bot", "Devctl Release Bot", "devctlpluginreleasebot@gmail.com"
}

//New returns new releaser object
func New(ghToken string) *Releaser {
	tokenUserHandle, tokenUsername, tokenEmail := getUserDetails(ghToken)

	return &Releaser{
		Token:                           ghToken,
		TokenEmail:                      tokenEmail,
		TokenUserHandle:                 tokenUserHandle,
		TokenUsername:                   tokenUsername,
		UpstreamDevctlIndexRepo:         devctl.GetDevctlIndexRepoName(),
		UpstreamDevctlIndexRepoOwner:    devctl.GetDevctlIndexRepoOwner(),
		UpstreamDevctlIndexRepoCloneURL: getCloneURL(devctl.GetDevctlIndexRepoOwner(), devctl.GetDevctlIndexRepoName()),
		LocalDevctlIndexRepo:            devctl.GetDevctlIndexRepoName(),
		LocalDevctlIndexRepoOwner:       tokenUserHandle,
		LocalDevctlIndexRepoCloneURL:    "https://github.com/alex-held/devctl-index.git",
	}
}

//HandleActionLambdaWebhook handles requests from github actions
func (releaser *Releaser) HandleActionLambdaWebhook(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	hook, err := actions.NewGithubActions()
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       errors.Wrap(err, "creating instance of action handler").Error(),
		}, nil
	}

	releaseRequest, err := hook.ParseLambdaRequest(request)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       errors.Wrap(err, "getting release request").Error(),
		}, nil
	}

	pr, err := releaser.Release(releaseRequest)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       errors.Wrap(err, "opening pr").Error(),
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("PR %q submitted successfully", pr),
	}, nil
}

//HandleActionWebhook handles requests from github actions
func (releaser *Releaser) HandleActionWebhook(w http.ResponseWriter, r *http.Request) {
	hook, err := actions.NewGithubActions()
	if err != nil {
		http.Error(w, errors.Wrap(err, "creating instance of action handler").Error(), http.StatusInternalServerError)
		return
	}

	releaseRequest, err := hook.Parse(r)
	if err != nil {
		http.Error(w, errors.Wrap(err, "getting release request").Error(), http.StatusInternalServerError)
		return
	}

	pr, err := releaser.Release(releaseRequest)
	if err != nil {
		http.Error(w, errors.Wrap(err, "opening pr").Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("PR %q submitted successfully", pr)))
}
