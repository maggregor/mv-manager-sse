package broadcaster

type SseEvent struct {
	Event     string `json:"event"`
	ProjectId string `json:"projectId"`
}
