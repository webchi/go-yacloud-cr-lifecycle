package main

/*
yc container repository list
yc container repository lifecycle-policy list --repository-name crpbcv0kq3k8f2813aha/admin
*/

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	durationpb "google.golang.org/protobuf/types/known/durationpb"
	"time"

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
		log.Fatal(err)
	}
	getListDocker(ctx, sdk)
	//dockerList, _ := getListDocker(ctx, sdk)

	//list , err := createPolicy(ctx,sdk)
	//fmt.Println(dockerList)
	//fmt.Println(list,err)
}

func checkEnv(key string) {
	_, ok := os.LookupEnv(key)
	if !ok {
		fmt.Printf("%s not set\n", key)
	}
}

func getListDocker(ctx context.Context, sdk *ycsdk.SDK) (string, error) {
	listDocker, err := sdk.ContainerRegistry().Repository().List(ctx, &containerregistry.ListRepositoriesRequest{
		FolderId: os.Getenv("YANDEX_FOLDER_ID"),
	})
	for _, v := range listDocker.GetRepositories() {
		//fmt.Printf("Repo %v ID %v\n", i, v)
		parts := strings.Split(v.String(), ":")
		//fmt.Println(parts[2])
		data := strings.Replace(parts[2], "\"", "", 2)
		//fmt.Println(data)
		createPolicy(ctx, sdk, data)
	}
	//fmt.Println(fix)
	return listDocker.String(), err
}

func createPolicy(ctx context.Context, sdk *ycsdk.SDK, RepositoryId_In string) (string, error) {
	fmt.Println(RepositoryId_In)

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
				Description:  "clear feature-* images",
				TagRegexp:    "feature-.*",
				Untagged:     false,
				ExpirePeriod: durationpb.New(720 * time.Hour),
			},
			{
				Description:  "clear release-* images",
				TagRegexp:    "release-.*",
				Untagged:     false,
				ExpirePeriod: durationpb.New(720 * time.Hour),
			},
		},
	})
	fmt.Println(err)
	return createPolicy.String(), err
}