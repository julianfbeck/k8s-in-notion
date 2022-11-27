package main

import (
	"log"
	"os"

	"github.com/julianfbeck/k8s-in-notion/internal/kubernetes"
	"github.com/julianfbeck/k8s-in-notion/internal/notiondb"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func main() {
	notion_secret := os.Getenv("NOTION_SECRET")
	notionClient := notiondb.CreateClient(notion_secret)
	database := notiondb.CreateDatabase(notionClient)
	k8sClient := kubernetes.CreateClient()

	kubernetes.WatchForPods(k8sClient, func(p *v1.Pod, e watch.Event) {
		if e.Type == "ADDED" {
			property := notiondb.CreatePodsBlock(e, p)
			notiondb.UpdateDatabase(notionClient, database.ID, property)
			log.Default().Printf("Added pod %s to database %s", p.Name, database.ID)
		} else if e.Type == "DELETED" {
			notiondb.DeleteBlock(notionClient, database.ID, p)
			log.Default().Printf("Deleted pod %s from database %s", p.Name, database.ID)
		} else if e.Type == "MODIFIED" {
			property := notiondb.CreatePodsBlock(e, p)
			notiondb.UpdateBlock(notionClient, database.ID, p, property)
			log.Default().Printf("Updated pod %s in database %s", p.Name, database.ID)
		}

	})

}
