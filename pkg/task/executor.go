package task

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"github.com/zhangzh-pku/software-engineering/pkg/api/entity"
)

func GetScriptsAndDockerImage(doid string) (string, string, error) {
	manager := entity.GetManagerInstance()
	pin, err := manager.GetPin(doid)
	if err != nil {
		return "", "", err
	}
	if pin.Type != entity.ReproductionPinType {
		return "", "", fmt.Errorf("got unexcepted pin type %s", pin.Type)
	}
	links, err := manager.GetAllLinks(doid)
	if err != nil {
		return "", "", err
	}
	var scripts string
	var image string
	fmt.Printf("got len links %d", len(links))
	for i := range links {
		pin, err := manager.GetPin(links[i].To)
		if err != nil {
			return "", "", err
		}
		switch pin.Type {
		case entity.DockerPinType:
			if _image, ok := pin.Metadata[entity.DockerKey]; !ok {
				return "", "", fmt.Errorf("image pin format error. no image found in pin.metadata")
			} else {
				if image, ok = _image.(string); !ok {
					return "", "", fmt.Errorf("image pin format error. cant convert image to string")
				}
			}

		case entity.ScriptPinType:
			if _scripts, ok := pin.Metadata[entity.ScriptKey]; !ok {
				return "", "", fmt.Errorf("script pin format error. no scripts found in pin.metadata")
			} else {
				if scripts, ok = _scripts.(string); !ok {
					fmt.Println(_scripts)
					return "", "", fmt.Errorf("script pin format error. can't convert script to string")
				}
			}
		default:
			fmt.Printf("got type %s", pin.Type)
		}
	}
	if len(scripts) == 0 {
		return "", "", fmt.Errorf("script not found")
	}
	if image == "" {
		return "", "", fmt.Errorf("image not found")
	}
	return scripts, image, nil
}

func RunScriptsInDocker(scripts string, dockerImage string) (string, error) {
	// 创建Docker客户端
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}

	// 检查本地是否存在所需的Docker映像，如果不存在则拉取
	if err := checkAndPullImage(cli, dockerImage); err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	// 在Docker映像中运行脚本，并返回任务ID
	taskID, err := executeScripts(cli, scripts, dockerImage)

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

func generateTaskID() string {
	return uuid.New().String()
}
func executeScripts(cli *client.Client, scripts string, dockerImage string) (string, error) {
	taskID := generateTaskID()

	// 创建容器
	containerConfig := &container.Config{
		Image: dockerImage,
		Cmd:   []string{"tail", "-f", "/dev/null"},
	}
	workdir := "/root/workdir/" + taskID
	outputDir := filepath.Join("/output", taskID)
	pwd, _ := os.Getwd()
	fmt.Println(pwd)
	outputDir = pwd + outputDir
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}
	outDirIn := "/root/" + taskID
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/root/test_rg/abacus_input/000_I3Pb1C1N2H5/lcao",
				Target: workdir,
			},
			{
				Type:   mount.TypeBind,
				Source: outputDir,
				Target: outDirIn,
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

	// 等待容器启动
	for {
		containerInfo, err := cli.ContainerInspect(context.Background(), resp.ID)
		if err != nil {
			log.Fatalf("Failed to inspect container: %v", err)
		}

		if containerInfo.State.Running {
			break
		}

		if containerInfo.State.Status == "exited" || containerInfo.State.Status == "dead" {
			log.Fatalf("Container exited unexpectedly with status: %s", containerInfo.State.Status)
		}

		time.Sleep(1 * time.Second)
	}
	script := scripts
	// 在容器中顺序执行脚本
	execConfig := &types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"bash", "-c", "cd " + workdir + " && " + script},
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

	mvCmd := fmt.Sprintf("mv %s/* %s/", workdir, outDirIn)
	execConfigMv := &types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"bash", "-c", mvCmd},
	}

	execIDMv, err := cli.ContainerExecCreate(context.Background(), resp.ID, *execConfigMv)
	if err != nil {
		log.Fatalf("Failed to create exec configuration for mv command: %v", err)
	}

	responseMv, err := cli.ContainerExecAttach(context.Background(), execIDMv.ID, types.ExecStartCheck{})
	if err != nil {
		log.Fatalf("Failed to attach to exec for mv command: %v", err)
	}
	defer responseMv.Close()

	// 输出mv命令执行结果
	stdcopy.StdCopy(os.Stdout, os.Stderr, responseMv.Reader)

	execInspectMv, err := cli.ContainerExecInspect(context.Background(), execIDMv.ID)
	if err != nil {
		log.Fatalf("Failed to inspect exec for mv command: %v", err)
	}

	if execInspectMv.ExitCode != 0 {
		log.Printf("mv command exited with code %d\n", execInspectMv.ExitCode)
		// 报告错误，根据实际需求自定义错误信息
		err = fmt.Errorf("mv command execution failed with exit code %d", execInspectMv.ExitCode)
		return "", err
	}
	if execInspect.ExitCode != 0 {
		log.Printf("Script %s exited with code %d\n", script, execInspect.ExitCode)
		// 报告错误，根据实际需求自定义错误信息
		err = fmt.Errorf("script %s execution failed with exit code %d", script, execInspect.ExitCode)
		return "", err
	}

	// 停止并删除容器
	stopOptions := container.StopOptions{}
	if err := cli.ContainerStop(context.Background(), resp.ID, stopOptions); err != nil {
		log.Printf("Failed to stop container: %v", err)
	}
	if err := cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{}); err != nil {
		log.Printf("Failed to remove container: %v", err)
	}

	return taskID, nil
}
