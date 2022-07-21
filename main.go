package main

import (
	"context"
	"os"

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

		// Получаем список политик для регистри
		listPolicies, err := sdk.ContainerRegistry().LifecyclePolicy().List(ctx, &containerregistry.ListLifecyclePoliciesRequest{
			Id: &containerregistry.ListLifecyclePoliciesRequest_RepositoryId{
				RepositoryId: v.Id,
			},
		})

		policies := listPolicies.GetLifecyclePolicies()
		// Если политик нет то создаём нужную
		if len(policies) == 0 {
			_, err = createPolicy(ctx, sdk, v.Id)

			if err != nil {
				lumber.ErrorMsg(err)
			} else {
				lumber.Success("Lifecycle policy successfully created!")
			}
		} else {
			// Если есть ...
			for _, z := range listPolicies.GetLifecyclePolicies() {

				// То обновляем ту которая ACTIVE
				if z.Status == 1 {
					_, err = updatePolicy(ctx, sdk, z.Id)

					if err != nil {
						lumber.ErrorMsg(err)
					} else {
						lumber.Success("Lifecycle policy successfully updated!")
					}
				}

			}
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
		Rules:        PolicyRules,
	})

	return createPolicy.String(), err
}

func updatePolicy(ctx context.Context, sdk *ycsdk.SDK, LifecyclePolicyId string) (string, error) {
	updatePolicy, err := sdk.ContainerRegistry().LifecyclePolicy().Update(ctx, &containerregistry.UpdateLifecyclePolicyRequest{
		LifecyclePolicyId: LifecyclePolicyId,
		Name:              "delete",
		Description:       "delete",
		Status:            1,
		Rules:             PolicyRules,
	})

	return updatePolicy.String(), err
}

var PolicyRules = []*containerregistry.LifecycleRule{
	{
		Description:  "clear untagged images",
		Untagged:     true,
		ExpirePeriod: durationpb.New(168 * time.Hour),
	},
	{
		Description:  "remove all with tags",
		TagRegexp:    ".*",
		Untagged:     false,
		ExpirePeriod: durationpb.New(4320 * time.Hour),
		RetainedTop:  5,
	},
	{
		Description:  "clear feature-* images",
		TagRegexp:    "feature-.*",
		Untagged:     false,
		ExpirePeriod: durationpb.New(1440 * time.Hour),
	},
	{
		Description:  "clear hotfix-* images",
		TagRegexp:    "hotfix-.*",
		Untagged:     false,
		ExpirePeriod: durationpb.New(1440 * time.Hour),
	},
	{
		Description:  "clear release-* images",
		TagRegexp:    "release-.*",
		Untagged:     false,
		ExpirePeriod: durationpb.New(4320 * time.Hour),
	},
}
