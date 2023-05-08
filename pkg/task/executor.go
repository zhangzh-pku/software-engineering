package task

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
)

func GetScriptsAndDockerImage(doip string) ([]string, string) {
	// 替换为实际的接口调用，返回脚本列表和Docker映像名称
	return []string{"script1.sh", "script2.sh"}, "your-docker-image"
}

func RunScriptsInDocker(scripts []string, dockerImage string) (string, error) {
	// 创建Docker客户端
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}

	// 检查本地是否存在所需的Docker映像，如果不存在则拉取
	if err := checkAndPullImage(cli, dockerImage); err != nil {
		return "", err
	}

	// 调用API获取实际运行的脚本
	executableScripts, err := fetchExecutableScripts(scripts)
	if err != nil {
		return "", err
	}

	// 在Docker映像中运行脚本，并返回任务ID
	taskID, err := executeScripts(cli, executableScripts, dockerImage)

	return taskID, err
}

func GetTaskStatus(taskID string) string {
	// 在此处添加查询任务状态的逻辑，并返回任务状态
	return "running"
}

func checkAndPullImage(cli *client.Client, imageName string) error {
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return err
	}

	for _, image := range images {
		for _, repoTag := range image.RepoTags {
			if repoTag == imageName {
				return nil
			}
		}
	}

	fmt.Printf("Pulling Docker image: %s\n", imageName)
	reader, err := cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return err
	}

	return nil
}

func fetchExecutableScripts(scripts []string) ([]string, error) {
	// 在此处调用API以获取可执行脚本
	// 模拟从API获取可执行脚本的过程
	executableScripts := []string{}
	for _, script := range scripts {
		if strings.Contains(script, "executable") {
			executableScripts = append(executableScripts, script)
		}
	}

	if len(executableScripts) == 0 {
		return nil, errors.New("no executable scripts found")
	}

	return executableScripts, nil
}

func generateTaskID() string {
	return uuid.New().String()
}

func executeScripts(cli *client.Client, scripts []string, dockerImage string) (string, error) {
	// 创建一个包含所有脚本的临时文件夹
	tempDir, err := ioutil.TempDir("", "scripts")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}

	taskID := generateTaskID()
	// 将脚本写入临时文件夹
	for _, script := range scripts {
		err := ioutil.WriteFile(filepath.Join(tempDir, filepath.Base(script)), []byte(script), 0755)
		if err != nil {
			log.Fatalf("Failed to write script to temp directory: %v", err)
		}
	}

	outputDir := filepath.Join("output", taskID)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// 创建容器
	containerConfig := &container.Config{
		Image: dockerImage,
		Cmd:   []string{"/bin/bash"},
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: tempDir,
				Target: "/scripts",
			},
			{
				Type:   mount.TypeBind,
				Source: outputDir,
				Target: "/output",
			},
		},
	}

	networkConfig := &network.NetworkingConfig{}

	containerName := "my-container-" + taskID

	resp, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, networkConfig, nil, containerName)
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	// 启动容器
	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatalf("Failed to start container: %v", err)
	}

	// 在容器中顺序执行脚本
	for _, script := range scripts {
		execConfig := &types.ExecConfig{
			AttachStdout: true,
			AttachStderr: true,
			Cmd:          []string{"bash", "-c", filepath.Join("/scripts", filepath.Base(script))},
		}

		execID, err := cli.ContainerExecCreate(context.Background(), resp.ID, *execConfig)
		if err != nil {
			log.Fatalf("Failed to create exec configuration: %v", err)
		}

		response, err := cli.ContainerExecAttach(context.Background(), execID.ID, types.ExecStartCheck{})
		if err != nil {
			log.Fatalf("Failed to attach to exec: %v", err)
		}
		defer response.Close()

		// 输出脚本执行结果
		stdcopy.StdCopy(os.Stdout, os.Stderr, response.Reader)

		execInspect, err := cli.ContainerExecInspect(context.Background(), execID.ID)
		if err != nil {
			log.Fatalf("Failed to inspect exec: %v", err)
		}

		if execInspect.ExitCode != 0 {
			log.Printf("Script %s exited with code %d\n", script, execInspect.ExitCode)
			// 报告错误，根据实际需求自定义错误信息
			err = fmt.Errorf("script %s execution failed with exit code %d", script, execInspect.ExitCode)
			return "", err
		}
	}

	// 停止并删除容器
	// 10s
	stopOptions := container.StopOptions{}
	if err := cli.ContainerStop(context.Background(), resp.ID, stopOptions); err != nil {
		log.Printf("Failed to stop container: %v", err)
	}
	if err := cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{}); err != nil {
		log.Printf("Failed to remove container: %v", err)
	}

	return taskID, nil
}
