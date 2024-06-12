package main

import (
	"demo-ui/internal/docker"
	"demo-ui/internal/key_reader"
	"fmt"
	"time"
)

func main() {
	//win := ui.NewWindow()
	//win.ShowAndRun()

	dockerClient := docker.NewDockerClient()
	if err := dockerClient.Connect(); err != nil {
		panic(err)
	}

	// sokolko/export-cryptopro-cert
	//out, err := dockerClient.PullImage("dimajolkin/export-cryptopro-cert:1.0")
	//if err != nil {
	//	panic(err)
	//}
	//stdcopy.StdCopy(os.Stdout, os.Stderr, *out)

	containerName := "demo-ui-container"
	id, err := dockerClient.Run("dimajolkin/export-cryptopro-cert:1.0", containerName, []string{"sleep", "60"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s \n", id)

	reader := key_reader.CreateKeyReader()
	//pwd, _ := os.Getwd()
	//path := pwd + "/test/testdata/container.024"
	path := "/Users/dimajolkin/Project/Moydom/www/ansible/tmp/10561020.024"
	key, err := reader.OpenDir(path)
	if err != nil {
		panic(err)
	}

	err = dockerClient.Exec(containerName, []string{"mkdir", "/data/container.000"})
	if err != nil {
		panic(err)
	}

	err = dockerClient.Copy(containerName, "/data/container.000", key.Container)
	if err != nil {
		panic(err)
	}

	//err = dockerClient.Exec(id, []string{"get-cpcert", "container.000", ">", "/tmp/output"})
	//if err != nil {
	//	panic(err)
	//}

	//for name, content := range key.Container.Files {
	//	fmt.Printf("%s: %T  %s \n", name, content, id)
	//
	//
	//
	//	//err = dockerClient.Exec("test-container", []string{"touch", "/tmp/" + name})
	//	//if err != nil {
	//	//	panic(err)
	//	//}
	//
	//	r
	//
	//	err := dockerClient.Copy("test-container", "/tmp", bytes.NewReader(content))
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	//fmt.Printf("Container ID: %s", id)
	//stdcopy.StdCopy(os.Stdout, os.Stderr, *out)
	time.Sleep(10)

	println("Success")
}
