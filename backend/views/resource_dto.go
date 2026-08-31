package views

import "medicalBooking/models"

type CreateResourceRequest struct {
	Name        string              `json:"name" binding:"required,min=2,max=150"`
	Type        models.ResourceType `json:"type" binding:"required"`
	Description string              `json:"description"`
	Location    string              `json:"location"`
	Capacity    int                 `json:"capacity"`
}

type UpdateResourceRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=150"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Capacity    int    `json:"capacity"`
	Active      *bool  `json:"active"`
}

type ResourceResponse struct {
	ID          uint                `json:"id"`
	Name        string              `json:"name"`
	Type        models.ResourceType `json:"type"`
	Description string              `json:"description"`
	Location    string              `json:"location"`
	Capacity    int                 `json:"capacity"`
	Active      bool                `json:"active"`
}

func ToResourceResponse(r models.Resource) ResourceResponse {
	return ResourceResponse{
		ID:          r.ID,
		Name:        r.Name,
		Type:        r.Type,
		Description: r.Description,
		Location:    r.Location,
		Capacity:    r.Capacity,
		Active:      r.Active,
	}
}

func ToResourceResponseList(resources []models.Resource) []ResourceResponse {
	result := make([]ResourceResponse, 0, len(resources))
	for _, r := range resources {
		result = append(result, ToResourceResponse(r))
	}
	return result
}
