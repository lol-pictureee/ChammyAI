package segmind

type Headers struct {
	APIKey string `json:"x-sigmind-key"`
}

type Data struct {
	SourceImg          string  `json:"source_img"`
	TargetImg          string  `json:"target_img"`
	InputFacesIndex    int     `json:"input_faces_index"`
	SourceFacesIndex   int     `json:"source_faces_index"`
	FaceRestore        string  `json:"face_restore"`
	Interpolation      string  `json:"interpolation"`
	DetectionFaceOrder string  `json:"detection_face_order"`
	FaceDetection      string  `json:"facedetection"`
	DetectGenderInput  string  `json:"detect_gender_input"`
	DetectGenderSource string  `json:"detect_gender_source"`
	FaceRestoreWeight  float64 `json:"face_restore_weight"`
	ImageFormat        string  `json:"image_format"`
	ImageQuality       int     `json:"image_quality"`
	Base64             bool    `json:"base64"`
}
