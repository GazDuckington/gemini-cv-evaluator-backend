package dto

import "mime/multipart"

type SubmitCvRequest struct {
	UserID  string                `form:"user_id"`
	Title   string                `form:"title" binding:"required,max=150"`
	File    *multipart.FileHeader `form:"file" binding:"required"`
	Summary string                `form:"summary,omitempty"`
}

type CVResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	FilePath  string    `json:"file_path"`
	Summary   string    `json:"summary"`
	Embedding []float32 `json:"embedding,omitempty"`
}

type WorkerStatusResponse struct {
	ID     string                `json:"id"`
	Status string                `json:"status"`
	Eval   *CVEvaluationResponse `json:"evaluation"`
}

type CVEvaluationResponse struct {
	CVMatchRate     float64 `json:"cv_match_rate"`
	CVFeedback      string  `json:"cv_feedback"`
	ProjectScore    float64 `json:"project_score"`
	ProjectFeedback string  `json:"project_feedback"`
	OverallSummary  string  `json:"overall_summary"`
}
