package vikunja

// User is a struct that represents a user in Vikunja
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Created  string `json:"created"`
	Updated  string `json:"updated"`
}

// Project is a struct that represents a project in Vikunja
type Project struct {
	ID                    int         `json:"id"`
	Title                 string      `json:"title"`
	HexColor              string      `json:"hex_color"`
	BackgroundBlurHash    string      `json:"background_blur_hash"`
	Position              float64     `json:"position"`
	Created               string      `json:"created"`
	Updated               string      `json:"updated"`
	Description           string      `json:"description"`
	DefaultBucketID       int         `json:"default_bucket_id"`
	DoneBucketID          int         `json:"done_bucket_id"`
	Identifier            string      `json:"identifier"`
	IsArchived            bool        `json:"is_archived"`
	IsFavorite            bool        `json:"is_favorite"`
	Owner                 User        `json:"owner"`
	ParentProjectID       int         `json:"parent_project_id"`
	BackgroundInformation interface{} `json:"background_information"`
}

// Task represents a task in Vikunja
type Task struct {
	ID                     int                      `json:"id"`
	Title                  string                   `json:"title"`
	Description            string                   `json:"description"`
	Done                   bool                     `json:"done"`
	DoneAt                 string                   `json:"done_at"`
	DueDate                string                   `json:"due_date"`
	Reminders              interface{}              `json:"reminders"`
	ProjectID              int                      `json:"project_id"`
	RepeatAfter            int                      `json:"repeat_after"`
	RepeatMode             int                      `json:"repeat_mode"`
	Priority               int                      `json:"priority"`
	StartDate              string                   `json:"start_date"`
	EndDate                string                   `json:"end_date"`
	Assignees              []User                   `json:"assignees"`
	Labels                 []Label                  `json:"labels"`
	HexColor               string                   `json:"hex_color"`
	PercentDone            float32                  `json:"percent_done"`
	Identifier             string                   `json:"identifier"`
	Index                  int                      `json:"index"`
	RelatedTasks           map[string][]interface{} `json:"related_tasks,omitempty"` // Marking this with []Task instead of interface{} causes parsing errors...
	Attachments            interface{}              `json:"attachments"`
	CoverImageAttachmentID int                      `json:"cover_image_attachment_id"`
	IsFavorite             bool                     `json:"is_favorite"`
	Created                string                   `json:"created"`
	Updated                string                   `json:"updated"`
	BucketID               int                      `json:"bucket_id"`
	Position               float64                  `json:"position"`
	KanbanPosition         float64                  `json:"kanban_position"`
	CreatedBy              User                     `json:"created_by"`
}

// Webhook represents a webhook in Vikunja
type Webhook struct {
	ID        int      `json:"id"`
	TargetURL string   `json:"target_url"`
	Events    []string `json:"events"`
	ProjectID int      `json:"project_id"`
	Secret    *string  `json:"secret,omitempty"`
	Created   *string  `json:"created,omitempty"`
	Updated   *string  `json:"updated,omitempty"`
	CreatedBy *User    `json:"created_by,omitempty"`
}

// WebhookCallbackData represents the data that is sent to a webhook callback
type WebhookCallbackData struct {
	Doer User `json:"doer"`
	Task Task `json:"task"`
}

// WebhookCallback represents a webhook callback
type WebhookCallback struct {
	EventName VikunjaWebhookEventType `json:"event_name"`
	Time      string                  `json:"time"`
	Data      WebhookCallbackData     `json:"data"`
}

type Comment struct {
	ID      int    `json:"id"`
	Comment string `json:"comment"`
	Author  User   `json:"author"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}

type Label struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	HexColor    string `json:"hex_color"`
	CreatedBy   User   `json:"created_by"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

// LabelID is a struct used to communicate to vikunja which label you are after. It is not the ID field in Label though.
type LabelID struct {
	ID      int    `json:"label_id"`
	Created string `json:"created,omitempty"`
}

type VikunjaWebhookEventType string

const (
	ProjectDeleted        VikunjaWebhookEventType = "project.deleted"
	TaskAssigneeCreated   VikunjaWebhookEventType = "task.assignee.created"
	TaskCommentCreated    VikunjaWebhookEventType = "task.comment.created"
	TaskDeleted           VikunjaWebhookEventType = "task.deleted"
	TaskRelationCreated   VikunjaWebhookEventType = "task.relation.created"
	TaskCommentDeleted    VikunjaWebhookEventType = "task.comment.deleted"
	TaskAssigneeDeleted   VikunjaWebhookEventType = "task.assignee.deleted"
	ProjectSharedTeam     VikunjaWebhookEventType = "project.shared.team"
	ProjectSharedUser     VikunjaWebhookEventType = "project.shared.user"
	TaskAttachmentCreated VikunjaWebhookEventType = "task.attachment.created"
	TaskCommentEdited     VikunjaWebhookEventType = "task.comment.edited"
	TaskRelationDeleted   VikunjaWebhookEventType = "task.relation.deleted"
	TaskUpdated           VikunjaWebhookEventType = "task.updated"
	TaskCreated           VikunjaWebhookEventType = "task.created"
	TaskAttachmentDeleted VikunjaWebhookEventType = "task.attachment.deleted"
	ProjectUpdated        VikunjaWebhookEventType = "project.updated"
)
