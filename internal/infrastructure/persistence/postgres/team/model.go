package team

import (
	"time"

	"github.com/google/uuid"
)

// team — команда
type teamModel struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	OwnerID   uuid.UUID `db:"owner_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// teamMember — участник команды
type teamMemberModel struct {
	ID        uuid.UUID  `db:"id"`
	TeamID    uuid.UUID  `db:"team_id"`
	UserID    uuid.UUID  `db:"user_id"`
	Role      string     `db:"role"` // owner, admin, member
	InvitedBy *uuid.UUID `db:"invited_by"`
	JoinedAt  time.Time  `db:"joined_at"`
}
