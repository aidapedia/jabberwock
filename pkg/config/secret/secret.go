package secret

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/aidapedia/gdk/environment"
	"github.com/aidapedia/jabberwock/pkg/config/model"
)

// projectIDMap is the map of project id
var projectIDMap = map[string]string{
	environment.Development: "",
	environment.Staging:     "",
	environment.Production:  "",
}

func GetSecret() (model.SecretConfig, error) {
	var (
		secretCfg model.SecretConfig
		byteCfg   []byte
	)
	ctx := context.Background()
	if environment.GetAppEnvironment() == environment.Development {
		// 	// Update read local json
		// } else {
		client, err := secretmanager.NewClient(ctx)
		if err != nil {
			return secretCfg, err
		}
		defer client.Close()

		projectID, ok := projectIDMap[environment.GetAppEnvironment()]
		if !ok || projectID == "" {
			return secretCfg, errors.New("project id not found")
		}

		resp, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
			Name: fmt.Sprintf("projects/%s/secrets/domea/versions/latest", projectID),
		})
		if err != nil {
			return secretCfg, err
		}
		byteCfg = resp.GetPayload().GetData()
	}

	err := json.Unmarshal(byteCfg, &secretCfg)
	if err != nil {
		return secretCfg, err
	}

	return secretCfg, nil
}
