package main

import (
        "fmt"
        "log"
        "context"
        "net/http"
	"os"
        "path/filepath"
	"strings"
        "github.com/gin-gonic/gin"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	toml "github.com/pelletier/go-toml"
)

type MediaUploader interface{
  Initialize()
  StoreFile()
}

type MediaStorage struct{
  url string
  container string
  client azblob.Client
  host string
  port int64
}

type Body struct {
  // json tag to de-serialize json body
  Filepath string `json:"win_filepath" binding:"required"`
}

func handleFatalError(err error){
  if err != nil {
    log.Fatal(err.Error())
  }
}

func handleNonFatalError(err error){
  if err != nil {
    log.Println(err.Error())
  }
}

func (ms *MediaStorage) Initialize(ctx context.Context){
  log.Println("Initializing MediaStorage...")
  config_path, err := filepath.Abs("storage.toml")
  handleFatalError(err)
  config, err := toml.LoadFile(config_path)
  handleFatalError(err)
  os.Setenv("AZURE_CLIENT_ID", config.Get("auth.azure_client_id").(string))
  os.Setenv("AZURE_TENANT_ID", config.Get("auth.azure_tenant_id").(string))
  os.Setenv("AZURE_CLIENT_SECRET", config.Get("auth.azure_client_secret").(string))
  log.Println("Using CLIENT_ID=",config.Get("auth.azure_client_id").(string),"and TENANT_ID=",config.Get("auth.azure_tenant_id").(string))
  ms.url = config.Get("storage.url").(string)
  ms.container = config.Get("storage.container").(string)
  log.Println("Will be writing data to container:", ms.url)

  credential, err := azidentity.NewAzureCLICredential(nil)
  handleFatalError(err)

  client, err := azblob.NewClient(ms.url, credential, nil)
  handleFatalError(err)

  containerName := ms.container
  log.Println("Creating a container named:", containerName)
  _, err = client.CreateContainer(ctx, containerName, nil)
  handleNonFatalError(err)



  ms.host = config.Get("web_server.host").(string)
  ms.port = config.Get("web_server.port").(int64)
}

func linuxizeFilename (fileName string) string{
  tempString := strings.Replace(fileName, ":", "", -1)
  newFilename := "/mnt/" + tempString
  fmt.Println("New filename is:",newFilename)
  return newFilename
}


func (ms *MediaStorage) StoreFile(ctx context.Context, fileName string){
  newFilename := linuxizeFilename(fileName)
  if _,err := os.Stat(newFilename); err != nil{
    log.Printf("File '%s' does not exist!", fileName)
  } else {
    data,err := os.ReadFile(newFilename)
    handleNonFatalError(err)
    credential, err := azidentity.NewAzureCLICredential(nil)
    handleFatalError(err)
    client, err := azblob.NewClient(ms.url, credential, nil)
    handleFatalError(err)
    _, blobName := filepath.Split(newFilename)
    log.Println("Uploading a blob named:", blobName)
    _, err = client.UploadBuffer(ctx, ms.container, blobName, data, &azblob.UploadBufferOptions{})
    handleNonFatalError(err)
    log.Println("Upload completed successfully!")
}
}

func main(){
  ms := &MediaStorage{}
  ctx := context.Background()
  ms.Initialize(ctx)

  r := gin.New()
  r.MaxMultipartMemory = 8<<20
  r.GET("/version", func(c *gin.Context){
    c.JSON(http.StatusOK, gin.H{
      "app": "Web-Based Media Objects Uploader",
      "version": "0.1",
    })
  })
  r.POST("/upload", func(context *gin.Context) {
    body:=Body{}
    // using BindJson method to serialize body with struct
    if err:=context.BindJSON(&body);err!=nil{
      context.AbortWithError(http.StatusBadRequest,err)
      return
    }
    context.JSON(http.StatusAccepted,&body)
    fmt.Println("Will upload file: ",body.Filepath)
    ms.StoreFile(ctx, body.Filepath)
  })
  server_url := fmt.Sprintf(":%d", ms.port)
  r.Run(server_url)
}

