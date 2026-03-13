package domain

import "time"

type TicketStatus string

const (
	TicketNew        TicketStatus = "NEW"
	TicketInProgress TicketStatus = "IN_PROGRESS"
	TicketDone       TicketStatus = "DONE"
)

type Ticket struct {
	ID          UUID         `json:"id"`
	UserID      UUID         `json:"userId"`
	AdminID     *UUID        `json:"adminId,omitempty"`
	TGChatID    *int64       `json:"tgChatId,omitempty"`
	Subject     string       `json:"subject"`
	Status      TicketStatus `json:"status"`
	UserMessage string       `json:"userMessage"`
	AdminReply  string       `json:"adminReply,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
	ClosedAt    *time.Time   `json:"closedAt,omitempty"`
}

func NewTicket(userID UUID, subject, message string, tgChatID *int64) *Ticket {
	return &Ticket{
		ID:          NewUUID(),
		UserID:      userID,
		Subject:     subject,
		Status:      TicketNew,
		UserMessage: message,
		TGChatID:    tgChatID,
		CreatedAt:   time.Now().UTC(),
	}
}

func (t *Ticket) Take(adminID UUID) {
	t.AdminID = &adminID
	t.Status = TicketInProgress
}

func (t *Ticket) Close(reply string) {
	t.AdminReply = reply
	t.Status = TicketDone
	now := time.Now().UTC()
	t.ClosedAt = &now
}
