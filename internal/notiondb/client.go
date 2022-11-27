package notiondb

import (
	"context"
	"log"

	"github.com/dstotijn/go-notion"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func CreateClient(token string) *notion.Client {
	client := notion.NewClient(token)
	return client
}

func UpdateDatabase(client *notion.Client, databaseID string, properties *notion.DatabasePageProperties) {
	_, err := client.CreatePage(context.Background(), notion.CreatePageParams{
		ParentID:               databaseID,
		ParentType:             notion.ParentTypeDatabase,
		DatabasePageProperties: properties,
	})
	if err != nil {
		log.Default().Println(err)
	}

}

func CreateDatabase(client *notion.Client) notion.Database {
	database, err := client.CreateDatabase(context.Background(), notion.CreateDatabaseParams{
		Title: []notion.RichText{
			{
				PlainText: "Kubernetes Cluster Pods",
				Text: &notion.Text{
					Content: "Kubernetes Cluster Pods",
				},
			},
		},
		ParentPageID: "06336be68e0b4278999eb22c6b461a26",
		Properties: notion.DatabaseProperties{
			"Name": notion.DatabaseProperty{
				Type:  "title",
				Name:  "Name",
				Title: &notion.EmptyMetadata{},
			},

			"namespace": notion.DatabaseProperty{
				Type: "select",
				Name: "namespace",
				Select: &notion.SelectMetadata{
					Options: []notion.SelectOptions{},
				},
			},
			"node": notion.DatabaseProperty{
				Type: "select",
				Name: "node",
				Select: &notion.SelectMetadata{
					Options: []notion.SelectOptions{},
				},
			},
			"status": notion.DatabaseProperty{
				Type:     "rich_text",
				Name:     "status",
				RichText: &notion.EmptyMetadata{},
			},
			"id": notion.DatabaseProperty{
				Type:     "rich_text",
				Name:     "id",
				RichText: &notion.EmptyMetadata{},
			},

			"date": notion.DatabaseProperty{
				Type:     "rich_text",
				Name:     "date",
				RichText: &notion.EmptyMetadata{},
			},
		},
	})
	if err != nil {
		log.Default().Println(err)
	}

	return database
}

func DeleteBlock(client *notion.Client, databaseID string, pod *v1.Pod) {

	blocks, err := client.QueryDatabase(context.Background(), databaseID, &notion.DatabaseQuery{})

	if err != nil {
		log.Default().Println(err)
	}

	for _, block := range blocks.Results {
		blockProperties := block.Properties.(notion.DatabasePageProperties)
		if blockProperties["id"].RichText[0].Text.Content == string(pod.UID) {
			client.DeleteBlock(context.Background(), block.ID)
		}

	}

}

func UpdateBlock(client *notion.Client, databaseID string, pod *v1.Pod, properties *notion.DatabasePageProperties) {

	blocks, err := client.QueryDatabase(context.Background(), databaseID, &notion.DatabaseQuery{})

	if err != nil {
		log.Default().Println(err)
	}

	for _, block := range blocks.Results {
		blockProperties := block.Properties.(notion.DatabasePageProperties)
		if blockProperties["id"].RichText[0].Text.Content == string(pod.UID) {
			client.UpdatePage(context.Background(), block.ID, notion.UpdatePageParams{
				DatabasePageProperties: *properties,
			})
		}

	}
}

func CreatePodsBlock(event watch.Event, pod *v1.Pod) *notion.DatabasePageProperties {
	blockPods := &notion.DatabasePageProperties{
		"Name": notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: pod.Name,
					},
				},
			},
		},
		"id": notion.DatabasePageProperty{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: string(pod.UID),
					},
				},
			},
		},
		"namespace": notion.DatabasePageProperty{
			Select: &notion.SelectOptions{
				Name: pod.Namespace,
			},
		},
		"status": notion.DatabasePageProperty{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: string(pod.Status.Phase),
					},
				},
			},
		},
		"date": notion.DatabasePageProperty{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: string(pod.CreationTimestamp.String()),
					},
				},
			},
		},
	}
	//append to blockPods
	newMap := make(map[string]notion.DatabasePageProperty)
	for k, v := range *blockPods {
		newMap[k] = v
	}

	if pod.Spec.NodeName != "" {
		newMap["node"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{
				Name: pod.Spec.NodeName,
			},
		}
	}

	properties := notion.DatabasePageProperties(newMap)
	return &properties

}
