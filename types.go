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
	success      bool
	file_path    string
	file_exists  bool
	file_content string
	access_error string
}

// keyRotationResponse contains the result of a key rotation operation
type keyRotationResponse struct {
	success        bool
	new_key        string
	rotation_error string
}