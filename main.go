package main

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	bucket, err = GetBucket()
)

func main() {
	e := echo.New()

	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	e.GET("/", hello)
	e.POST("/buildApp", buildApp)
	//e.Logger.Fatal(e.Start(":80"))
	e.Logger.Fatal(e.StartAutoTLS(":443"))
	//e.Logger.Fatal(e.StartTLS(":443", "server.crt", "server.key"))
	//s := http.Server{
	//	Addr:    ":443",
	//	Handler: e, // set Echo as handler
	//	TLSConfig: &tls.Config{
	//		MinVersion: tls.VersionTLS11, // customize TLS configuration
	//		MaxVersion: tls.VersionTLS12,
	//	},
	//	//ReadTimeout: 30 * time.Second, // use custom timeouts
	//}
	//e.Logger.Fatal(s.ListenAndServeTLS("server.crt", "server.key"))
}

// e.GET("/", hello)
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func buildApp(c echo.Context) error {
	props := new(BuildAppProps)
	if err := c.Bind(props); err != nil {
		return err
	}

	cacheID := uuid.NewV4().String()
	fmt.Printf("Generated UUID: %s\n", cacheID)

	cacheDir := "/tmp/" + cacheID
	if err := os.MkdirAll(cacheDir, 0644); err != nil {
		log.Println(err)
	}

	binName := "app"
	//if props.Bin != "" {
	//	binName = props.Bin
	//}

	//文件名
	filename := cacheDir + "/" + binName + ".py"
	//要写入的内容
	content := []byte(props.Code)

	//创建或打开文件（如果已存在则追加）
	err := os.WriteFile(filename, content, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}

	cmd := fmt.Sprintf("cd %s && pigar generate --auto-select --question-answer=yes --index-url=http://mirrors.cloud.aliyuncs.com/pypi/simple/ && pip install -r requirements.txt && pyinstaller --onefile %s.py", cacheDir, binName)

	command := exec.Command("cmd", "/c", cmd)
	//command := exec.Command("bash", "-c", cmd)
	output, err := command.CombinedOutput()
	//output, err := exec.Command("bash", "-lc", cmd).CombinedOutput()
	if err != nil {
		log.Println("编译失败")
		log.Println(string(output))
		//return err
		return c.JSON(http.StatusUnprocessableEntity, BuildAppResult{
			Success:      false,
			ErrorMessage: string(output),
		})
	}
	log.Println(string(output))

	ext := ".exe"

	binFilename := cacheDir + "/dist/" + binName + ext

	if _, err := os.Stat(binFilename); err != nil && !os.IsExist(err) {
		log.Println("新程序不存在")
		return err
	}

	//md5Hash := md5.Sum(content)
	//md5Hex := hex.EncodeToString(md5Hash[:])
	//objectKey := "autoapp/" + md5Hex + ext
	//
	//err = bucket.UploadFile(objectKey, binFilename, 1024*1024)
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	//
	//err = os.RemoveAll(cacheDir)
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}

	return c.File(binFilename)
	//return c.JSON(http.StatusCreated, props)
	//return c.JSON(http.StatusCreated, BuildAppResult{
	//	Success: true,
	//	Data:    objectKey,
	//})
}

type BuildAppProps struct {
	Code string `json:"code" xml:"code" form:"code" query:"code"`
	Bin  string `json:"bin" xml:"bin" form:"bin" query:"bin"`
	ID   string `json:"id" xml:"id" form:"id" query:"id"`
}

type BuildAppResult struct {
	Success      bool   `json:"success" xml:"success" form:"success" query:"success"`
	ErrorMessage string `json:"errorMessage" xml:"errorMessage" form:"errorMessage" query:"errorMessage"`
	Data         string `json:"data" xml:"data" form:"data" query:"data"`
}

// GetBucket get bucket
func GetBucket() (*oss.Bucket, error) {
	// New client
	client, err := oss.New(endpoint, accessID, accessKey)
	if err != nil {
		return nil, err
	}

	// Get bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}
