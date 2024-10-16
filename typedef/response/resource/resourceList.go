package resource

type ResourceListRes struct {
	FileUrl  string `json:"fileUrl"`  //文件url
	FileName string `json:"fileName"` //文件名
	FileType string `json:"fileType"` //文件类型
}
