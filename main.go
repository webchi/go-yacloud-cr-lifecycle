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
	d := 168 * time.Hour
	createPolicy, err := sdk.ContainerRegistry().LifecyclePolicy().Create(ctx, &containerregistry.CreateLifecyclePolicyRequest{
		RepositoryId: RepositoryId_In,
		Name:         "delete",
		Description:  "delete",
		Status:       1,
		Rules: []*containerregistry.LifecycleRule{
			// &containerregistry.LifecycleRule{
			// 	Description:  "delete after week",
			// 	Untagged:     false,
			// 	ExpirePeriod: durationpb.New(d),
			// },
			&containerregistry.LifecycleRule{
				Description:  "delete after week",
				Untagged:     true,
				ExpirePeriod: durationpb.New(d),
			},
		},
	})
	fmt.Println(err)
	return createPolicy.String(), err
}

func main() {
	token := os.Getenv("YANDEX_QAUTH_TOCKEN")
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
