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
	success      bool `json:"success"`
	file_path    string `json:"file_path"`
	file_exists  bool `json:"file_exists"`
	file_content string `json:"file_content"`
	access_error string `json:"access_error"`
}

// keyRotationResponse contains the result of a key rotation operation
type keyRotationResponse struct {
	success        bool `json:"success"`
	new_key        string `json:"new_key"`
	rotation_error string `json:"rotation_error"`
}