package main

// Tasks represents the root structure containing a list of tasks
type Tasks struct {
	Tasks []Task `json:"tasks"`
}

// Task represents a single task with its type and optional file path
type Task struct {
	TaskType string `json:"type"`
	FilePath string `json:"file_path,omitempty"`
}

// controlCheckResponse contains the result of a control file check
type controlCheckResponse struct {
	Success     bool   `json:"success"`
	FilePath    string `json:"file_path"`
	FileExists  bool   `json:"file_exists"`
	FileContent string `json:"file_content"`
	AccessError string `json:"access_error"`
}

// keyRotationResponse contains the result of a key rotation operation
type keyRotationResponse struct {
	Success       bool   `json:"success"`
	NewKey        string `json:"new_key"`
	RotationError string `json:"rotation_error"`
}
