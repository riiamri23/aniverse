package types

type UploadStatus struct {
	Success bool       `json:"success"`
	Result  ResultData `json:"result"`
}

type ResultData map[string]FileData

type FileData struct {
	FileRef  string  `json:"fileref"`
	Title    string  `json:"title"`
	Size     string  `json:"size"`
	Duration *string `json:"duration"`
	Subtitle bool    `json:"subtitle"`
	IsVideo  bool    `json:"isvideo"`
	IsAudio  bool    `json:"isaudio"`
	Added    string  `json:"added"`
	Status   string  `json:"status"`
	Deleted  bool    `json:"deleted"`
	Thumb    *string `json:"thumb"`
	URL      string  `json:"url"`
	YourFile bool    `json:"yourfile"`
}
