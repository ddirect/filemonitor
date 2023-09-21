package common

const FILE_MONITOR_DEFAULT_COLLECTION = "filemonitor"
const FILE_INFO_ID_FIELD_NAME = "Id"
const FILE_INFO_FILE_FIELD_NAME = "FileName"

type FileInfo struct {
	Id       string `json:",omitempty"`
	FileName string
	ModTime  int64 `json:",string"`
	Hash     []byte
}
