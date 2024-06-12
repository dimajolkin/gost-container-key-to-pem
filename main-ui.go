package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"

	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	// Получение списка запуцщенных контейнеров(docker ps)
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		panic(err)
	}

	// Вывод всех идентификаторов контейнеров
	for _, container := range containers {
		fmt.Println(container.ID)
	}
}

//for _, img := range imgs {
//	fmt.Println("ID: ", img.ID)
//	fmt.Println("RepoTags: ", img.RepoTags)
//	fmt.Println("Created: ", img.Created)
//	fmt.Println("Size: ", img.Size)
//	fmt.Println("VirtualSize: ", img.VirtualSize)
//	fmt.Println("ParentId: ", img.ParentID)
//}

//}
