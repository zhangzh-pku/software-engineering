package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/zhangzh-pku/software-engineering/pkg/api/entity"
	"github.com/zhangzh-pku/software-engineering/pkg/task"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/run", runHandler)
	router.GET("/task_status/:taskID", taskStatusHandler)
	router.GET("/tasks/:taskID/files", listGeneratedFiles)
	router.GET("/tasks/:taskID/files/:fileName", downloadGeneratedFile)
	router.POST("/task/create", generateReproduction)
	router.GET("/reproductions", getAllReproductions)
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

	scripts, dockerImage, err := task.GetScriptsAndDockerImage(requestBody.DOID)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	taskID, err := task.RunScriptsInDocker(scripts, dockerImage)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"taskID": taskID,
	})
}

func taskStatusHandler(c *gin.Context) {
	taskID := c.Param("taskID")

	taskStatus := task.GetTaskStatus(taskID)

	c.JSON(http.StatusOK, gin.H{
		"status": taskStatus,
	})
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
	}
	if err := c.BindJSON(&requestBody); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	manager := entity.GetManagerInstance()
	err := manager.GenerateReproduction(requestBody.DockerImage, requestBody.RunScript, requestBody.DissertationDOI)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
