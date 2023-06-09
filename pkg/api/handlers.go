package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zhangzh-pku/software-engineering/pkg/api/entity"
	"github.com/zhangzh-pku/software-engineering/pkg/file"
	"github.com/zhangzh-pku/software-engineering/pkg/task"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/run", runHandler)
	router.GET("/task_status/:taskID", taskStatusHandler)
	router.GET("/tasks/:taskID/files", listGeneratedFiles)
	router.GET("/tasks/:taskID/files/:fileName", downloadGeneratedFile)
	router.POST("/task/create", generateReproduction)
	router.GET("/reproductions", getAllReproductions)
	router.POST("/upload", uploadFiles)
	router.GET("/reproduction/*reproductionID", getReproduction)
	router.GET("/files/:taskId", getTaskFilesHandler)
	router.GET("/file/:taskId/:filename", getSingleFileHandler)
}

func getSingleFileHandler(c *gin.Context) {
	taskId := c.Param("taskId")
	filename := c.Param("filename")
	parentDir := "/data/zzh/workspace/"
	c.File(parentDir + taskId + "/" + filename)
}

func getTaskFilesHandler(c *gin.Context) {
	taskId := c.Param("taskId")

	// 替换为实际的父目录
	parentDir := "/data/zzh/workspace/"

	files, err := ioutil.ReadDir(parentDir + taskId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileNames := make([]string, 0, len(files))
	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}
	c.JSON(http.StatusOK, gin.H{"files": fileNames})
}

func runHandler(c *gin.Context) {
	var requestBody struct {
		DOID string `json:"doid"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	taskID := task.GenerateTaskID()
	fmt.Println("got taskID: " + taskID)
	tm := task.GetTaskManager()
	tm.AddTask(taskID)
	go func() {
		//认为路径应当是在/data/zzh这个路径下的 并且需要把它copy到/data/zzh/workspace下面去
		scripts, dockerImage, path, err := task.GetScriptsAndDockerImageAndPath(requestBody.DOID)
		if err != nil {
			fmt.Println(err.Error())
			tm.ChangeTaskStatus(taskID, task.FAILED)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		workPath := filepath.Join("/data/zzh/workspace", taskID)
		fmt.Println("workPath: " + workPath)
		err = file.CopyDir(path, workPath)

		if err != nil {
			fmt.Println("Failed to copy files.")
			tm.ChangeTaskStatus(taskID, task.FAILED)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = task.RunScriptsInDocker(scripts, dockerImage, workPath, taskID)
		if err != nil {
			fmt.Println(err.Error())
			tm.ChangeTaskStatus(taskID, task.FAILED)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tm.ChangeTaskStatus(taskID, task.SUCCESS)
	}()

	c.JSON(http.StatusOK, gin.H{
		"taskID": taskID,
	})
}

func taskStatusHandler(c *gin.Context) {
	taskID := c.Param("taskID")

	tm := task.GetTaskManager()
	status := tm.GetTaskStatus(taskID)
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "task not found",
		})
		return
	}
	c.JSON(http.StatusOK, status)
}

func listGeneratedFiles(c *gin.Context) {
	taskID := c.Param("taskID")
	outputDir := filepath.Join("output", taskID)

	files, err := ioutil.ReadDir(outputDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileInfos := make([]map[string]interface{}, 0)
	for _, file := range files {
		fileInfo := map[string]interface{}{
			"name": file.Name(),
			"size": file.Size(),
			"mode": file.Mode().String(),
			"time": file.ModTime(),
		}
		fileInfos = append(fileInfos, fileInfo)
	}

	c.JSON(http.StatusOK, fileInfos)
}

func downloadGeneratedFile(c *gin.Context) {
	taskID := c.Param("taskID")
	fileName := c.Param("fileName")
	outputDir := filepath.Join("output", taskID)
	filePath := filepath.Join(outputDir, fileName)

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(filePath)
}

func generateReproduction(c *gin.Context) {
	var requestBody struct {
		DockerImage     string `json:"docker_image"`
		RunScript       string `json:"run_script"`
		DissertationDOI string `json:"dissertation_doi"`
		Path            string `json:"data_path"`
		Charged         bool   `json:"data_charged"`
	}
	var err error
	if err = c.BindJSON(&requestBody); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(requestBody.DockerImage)
	m := entity.GetManagerInstance()
	taskID := task.GenerateTaskID()
	fmt.Println(
		"taskID:" + taskID,
	)
	tm := task.GetTaskManager()
	tm.AddTask(taskID)
	fmt.Println("task added with id:", taskID)
	go func() {

		err = task.RunScriptsInDocker(requestBody.RunScript, requestBody.DockerImage, requestBody.Path, taskID)

		if err != nil {
			fmt.Println(err)
			tm.ChangeTaskStatus(taskID, task.FAILED)
			return
		}
		// fmt.Println("m.GenerateReproduction(requestBody.DockerImage, requestBody.RunScript, requestBody.DissertationDOI, requestBody.Path)")
		// fmt.Println(requestBody.DockerImage)
		// fmt.Println(requestBody.RunScript)
		// fmt.Println(requestBody.DissertationDOI)
		// fmt.Println(requestBody.Path)
		reproductionPath := "/data/zzh/" + requestBody.DissertationDOI
		err = file.CopyDir(requestBody.Path, reproductionPath)
		if err != nil {
			fmt.Println(err)
			tm.ChangeTaskStatus(taskID, task.FAILED)
			return
		}
		// 跑完之后给复制到一个指定文件夹
		err = m.GenerateReproduction(requestBody.DockerImage, requestBody.RunScript, requestBody.DissertationDOI, reproductionPath, requestBody.Charged)
		if err != nil {
			fmt.Println(err)
			tm.ChangeTaskStatus(taskID, task.FAILED)
			return
		}
		tm.ChangeTaskStatus(taskID, task.SUCCESS)
	}()
	fmt.Println("taskID:" + taskID)
	c.JSON(http.StatusOK, taskID)
}

func getAllReproductions(c *gin.Context) {
	manager := entity.GetManagerInstance()
	reproductions, err := manager.GetAllReproductions()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reproductions)
}

func uploadFiles(c *gin.Context) {
	_file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有上传的文件或文件字段名错误。字段名应为 'file'"})
		return
	}
	// 验证文件类型
	if filepath.Ext(_file.Filename) != ".zip" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅接受.zip文件类型"})
		return
	}

	// 生成一个新的 UUID 作为目录名称
	dirUUID := uuid.New().String()

	// 创建一个新的文件保存路径
	savePath := path.Join("/data/zzh/tmp", dirUUID, _file.Filename)

	// 创建存储目录
	if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败"})
		return
	}

	// 保存上传的文件到指定路径
	if err := c.SaveUploadedFile(_file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 解压文件
	if err := file.UnzipFile(savePath, filepath.Dir(savePath)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解压文件失败"})
		return
	}

	// 返回成功的消息
	c.JSON(http.StatusOK, gin.H{"message": path.Join("/data/zzh/tmp", dirUUID)})
}

func getReproduction(c *gin.Context) {
	reproductionID := strings.TrimPrefix(c.Param("reproductionID"), "/")
	m := entity.GetManagerInstance()
	links, err := m.GetAllLinks(reproductionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var scripts string
	var image string
	var path string
	var doi string
	// todo move marshal
	fmt.Printf("len links%d", len(links))
	for i := range links {
		pin, err := m.GetPin(links[i].To)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		switch pin.Type {
		case entity.DockerPinType:
			if _image, ok := pin.Metadata[entity.DockerKey]; !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			} else {
				if image, ok = _image.(string); !ok {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}

		case entity.ScriptPinType:
			if _scripts, ok := pin.Metadata[entity.ScriptKey]; !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			} else {
				if scripts, ok = _scripts.(string); !ok {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		case entity.DatasetPinType:
			if _path, ok := pin.Metadata[entity.DataPathKey]; !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			} else {
				if path, ok = _path.(string); !ok {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		case entity.DissertationPinType:
			if _doi, ok := pin.Metadata[entity.DissertationKey]; !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			} else {
				if doi, ok = _doi.(string); !ok {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		default:
			fmt.Printf("got type %s", pin.Type)
		}
	}
	// if len(path) == 0 {
	// 	fmt.Println("path not found")
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "path not found"})
	// 	return
	// }
	// if len(scripts) == 0 {
	// 	fmt.Println("script not found")
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "script not found"})
	// }
	// if len(image) == 0 {
	// 	fmt.Println("image not found")
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "image not found"})
	// }
	// if len(doi) == 0 {
	// 	fmt.Println("doi not found")
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "doi not found"})
	// }
	pin, err := m.GetPin(reproductionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	charged := false

	if _charged, ok := pin.Metadata[entity.ChargedKey]; ok {
		if __charged, ok := _charged.(bool); ok {
			charged = __charged
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"script":  scripts,
		"doi":     doi,
		"image":   image,
		"path":    path,
		"charged": charged,
	})
}
