# k8s-in-notion
![example](https://raw.githubusercontent.com/julianfbeck/k8s-in-notion/main/media/example.gif)

> Custom k8s dashboard in Notion using [Notion API](https://developers.notion.com/)

This is a POC to show how to create a kubernetes dashboard in Notion using the Notion API and the kubernetes-client to display the current pods inside the cluster.

## Installation
Clone the repository and install the dependencies:
* ko
* kind

## Usage
Replace the `KO_DOCKER_REPO` in the Makefile with your docker repository.
Replace the NOTION_TOKEN and NOTION_PAGE_ID in the `k8s/deployment.yaml` with your Notion token and page id.

Deploy the custom controller to your notion page:
```bash
# Create a kind cluster
make cluster
# Deploy the controller
make deploy
```
