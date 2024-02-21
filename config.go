package main

import "os"

var (
	// Sample code's env configuration. You need to specify them with the actual configuration if you want to run sample code
	endpoint   = os.Getenv("OSS_ENDPOINT")
	accessID   = os.Getenv("OSS_ACCESS_KEY_ID")
	accessKey  = os.Getenv("OSS_ACCESS_KEY_SECRET")
	bucketName = "chatqa-cloud"
)
