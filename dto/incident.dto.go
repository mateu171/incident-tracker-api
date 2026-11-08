package dto

type CreateIncidentDTO struct {
	Title       string `json:"title" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"required"`
	Priority    string `json:"priority" binding:"required,oneof=low medium high critical"`
}

type UpdateIncidentStatusDTO struct {
	Status string `json:"status" binding:"required,oneof=open in_progress resolved closed"`
}
