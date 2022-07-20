package main

import (
	"context"
	"os"
	"strings"

	durationpb "google.golang.org/protobuf/types/known/durationpb"
	"time"

	"github.com/gleich/lumber/v2"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/containerregistry/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
)

func main() {
	checkEnv("YANDEX_OAUTH_TOKEN")
	checkEnv("YANDEX_FOLDER_ID")

	token := os.Getenv("YANDEX_OAUTH_TOKEN")
	ctx := context.Background()
	sdk, err := ycsdk.Build(ctx, ycsdk.Config{
		Credentials: ycsdk.OAuthToken(token),
	})
	if err != nil {
		lumber.Fatal(err)
	}
	getListDocker(ctx, sdk)
}

func checkEnv(key string) {
	_, ok := os.LookupEnv(key)
	if !ok {
		lumber.FatalMsg("%s not set\n", key)
	}
}

func getListDocker(ctx context.Context, sdk *ycsdk.SDK) (string, error) {
	listDocker, err := sdk.ContainerRegistry().Repository().List(ctx, &containerregistry.ListRepositoriesRequest{
		FolderId: os.Getenv("YANDEX_FOLDER_ID"),
	})
	for _, v := range listDocker.GetRepositories() {
		lumber.Info("Repo", v)
		parts := strings.Split(v.String(), ":")
		data := strings.Replace(parts[2], "\"", "", 2)

		_, err := createPolicy(ctx, sdk, data)
		if err != nil {
			lumber.ErrorMsg(err)
		} else {

			lumber.Success("Lifecycle policy successfully applied!")
		}
	}
	return listDocker.String(), err
}

func createPolicy(ctx context.Context, sdk *ycsdk.SDK, RepositoryId_In string) (string, error) {
	createPolicy, err := sdk.ContainerRegistry().LifecyclePolicy().Create(ctx, &containerregistry.CreateLifecyclePolicyRequest{
		RepositoryId: RepositoryId_In,
		Name:         "delete",
		Description:  "delete",
		Status:       1,
		Rules: []*containerregistry.LifecycleRule{
			{
				Description:  "clear untagged images",
				Untagged:     true,
				ExpirePeriod: durationpb.New(168 * time.Hour),
			},
			{
				Description:  "remove all older than month",
				TagRegexp:    ".*",
				Untagged:     false,
				ExpirePeriod: durationpb.New(720 * time.Hour),
				RetainedTop:  5,
			},
		},
	})

	return createPolicy.String(), err
}
