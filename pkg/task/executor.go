package task

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
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

func GetScriptsAndDockerImageAndPath(doid string) (string, string, string, error) {
	manager := entity.GetManagerInstance()
	pin, err := manager.GetPin(doid)
	if err != nil {
		return "", "", "", err
	}
	if pin.Type != entity.ReproductionPinType {
		return "", "", "", fmt.Errorf("got unexcepted pin type %s", pin.Type)
	}
	links, err := manager.GetAllLinks(doid)
	if err != nil {
		return "", "", "", err
	}
	var scripts string
	var image string
	var path string
	fmt.Printf("got len links %d", len(links))
	// todo merge marshal
	for i := range links {
		pin, err := manager.GetPin(links[i].To)
		if err != nil {
			return "", "", "", err
		}
		switch pin.Type {
		case entity.DockerPinType:
			if _image, ok := pin.Metadata[entity.DockerKey]; !ok {
				return "", "", "", fmt.Errorf("image pin format error. no image found in pin.metadata")
			} else {
				if image, ok = _image.(string); !ok {
					return "", "", "", fmt.Errorf("image pin format error. cant convert image to string")
				}
			}

		case entity.ScriptPinType:
			if _scripts, ok := pin.Metadata[entity.ScriptKey]; !ok {
				return "", "", "", fmt.Errorf("script pin format error. no scripts found in pin.metadata")
			} else {
				if scripts, ok = _scripts.(string); !ok {
					fmt.Println(_scripts)
					return "", "", "", fmt.Errorf("script pin format error. can't convert script to string")
				}
			}
		case entity.DatasetPinType:
			if _path, ok := pin.Metadata[entity.DataPathKey]; !ok {
				return "", "", "", fmt.Errorf("dataset pin format error. no dataset found in pin.metadata")
			} else {
				if path, ok = _path.(string); !ok {
					fmt.Println(_path)
					return "", "", "", fmt.Errorf("dataset pin format error. can't convert dataset path to string")
				}
			}
		default:
			fmt.Printf("got type %s", pin.Type)
		}
	}
	if len(path) == 0 {
		return "", "", "", fmt.Errorf("path not found")
	}
	if len(scripts) == 0 {
		return "", "", "", fmt.Errorf("script not found")
	}
	if len(image) == 0 {
		return "", "", "", fmt.Errorf("image not found")
	}
	return scripts, image, path, nil
}

func RunScriptsInDocker(scripts string, dockerImage string, mntPath string, taskID string) error {
	// 创建Docker客户端
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	// 检查本地是否存在所需的Docker映像，如果不存在则拉取
	if err := checkAndPullImage(cli, dockerImage); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	// 在Docker映像中运行脚本
	return executeScripts(cli, scripts, dockerImage, mntPath, taskID)
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

func GenerateTaskID() string {
	return uuid.New().String()
}

func executeScripts(cli *client.Client, scripts string, dockerImage string, mntPath string, taskID string) error {
	tm := GetTaskManager()
	tm.ChangeTaskStatus(taskID, RUNNING)
	// 创建容器
	containerConfig := &container.Config{
		Image: dockerImage,
		Cmd:   []string{"tail", "-f", "/dev/null"},
	}
	workdir := "/root/workdir/" + taskID
	pwd, _ := os.Getwd()
	fmt.Println(pwd)
	fmt.Println("mntPath", mntPath)
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: mntPath,
				Target: workdir,
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

	// fmt.Println("mvcmd:" + mvCmd)
	// execConfigMv := &types.ExecConfig{
	// 	AttachStdout: true,
	// 	AttachStderr: true,
	// 	Cmd:          []string{"bash", "-c", mvCmd},
	// }

	// execIDMv, err := cli.ContainerExecCreate(context.Background(), resp.ID, *execConfigMv)
	// if err != nil {
	// 	log.Fatalf("Failed to create exec configuration for mv command: %v", err)
	// }

	// responseMv, err := cli.ContainerExecAttach(context.Background(), execIDMv.ID, types.ExecStartCheck{})
	// if err != nil {
	// 	log.Fatalf("Failed to attach to exec for mv command: %v", err)
	// }
	// defer responseMv.Close()

	// 输出mv命令执行结果
	// stdcopy.StdCopy(os.Stdout, os.Stderr, responseMv.Reader)

	// execInspectMv, err := cli.ContainerExecInspect(context.Background(), execIDMv.ID)
	// if err != nil {
	// 	log.Fatalf("Failed to inspect exec for mv command: %v", err)
	// }

	// if execInspectMv.ExitCode != 0 {
	// 	log.Printf("mv command exited with code %d\n", execInspectMv.ExitCode)
	// 	// 报告错误，根据实际需求自定义错误信息
	// 	err = fmt.Errorf("mv command execution failed with exit code %d", execInspectMv.ExitCode)
	// 	return err
	// }
	if execInspect.ExitCode != 0 {
		log.Printf("Script %s exited with code %d\n", script, execInspect.ExitCode)
		// 报告错误，根据实际需求自定义错误信息
		err = fmt.Errorf("script %s execution failed with exit code %d", script, execInspect.ExitCode)
		return err
	}

	// 停止并删除容器
	stopOptions := container.StopOptions{}
	if err := cli.ContainerStop(context.Background(), resp.ID, stopOptions); err != nil {
		log.Printf("Failed to stop container: %v", err)
	}
	if err := cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{}); err != nil {
		log.Printf("Failed to remove container: %v", err)
	}

	return nil
}
