package team

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"gitlab.com/libs-artifex/wrapper/v2"

	"github.com/djalben/epn-killer-mvp/internal/entity"
	// добавляем зависимость от user repo
)

// CreateTeam создаёт новую команду и добавляет владельца как owner
func (r *Repository) CreateTeam(ctx context.Context, ownerID uuid.UUID, name string) (*entity.Team, error) {
	const query = `
		INSERT INTO teams (name, owner_id) 
		VALUES ($1, $2) 
		RETURNING id, name, owner_id, created_at, updated_at
	`

	var team entity.Team
	err := r.Client.QueryRowContext(ctx, query, name, ownerID).
		Scan(&team.ID, &team.Name, &team.OwnerID, &team.CreatedAt, &team.UpdatedAt)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to create team",
			slog.String("owner_id", ownerID.String()),
			slog.String("name", name),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}

	// Добавляем владельца в team_members как owner
	const addOwnerQuery = `
		INSERT INTO team_members (team_id, user_id, role, invited_by) 
		VALUES ($1, $2, 'owner', $2)
	`
	_, err = r.Client.ExecContext(ctx, addOwnerQuery, team.ID, ownerID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to add owner to team_members",
			slog.String("team_id", team.ID),
			slog.String("owner_id", ownerID.String()),
			slog.Any("error", err))

		// Откатываем создание команды
		r.Client.ExecContext(ctx, "DELETE FROM teams WHERE id = $1", team.ID)
		return nil, wrapper.Wrap(err)
	}

	r.Logger.InfoContext(ctx, "team created successfully",
		slog.String("team_id", team.ID),
		slog.String("owner_id", ownerID.String()),
		slog.String("name", name))

	return &team, nil
}

// GetUserTeams возвращает все команды пользователя
func (r *Repository) GetUserTeams(ctx context.Context, userID uuid.UUID) ([]entity.Team, error) {
	const query = `
		SELECT t.id, t.name, t.owner_id, t.created_at, t.updated_at
		FROM teams t
		INNER JOIN team_members tm ON t.id = tm.team_id
		WHERE tm.user_id = $1
		ORDER BY t.created_at DESC
	`

	rows, err := r.Client.QueryContext(ctx, query, userID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to get user teams",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}
	defer rows.Close()

	var teams []entity.Team
	for rows.Next() {
		var t entity.Team
		if err := rows.Scan(&t.ID, &t.Name, &t.OwnerID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			r.Logger.ErrorContext(ctx, "failed to scan team row", slog.Any("error", err))
			continue
		}
		teams = append(teams, t)
	}

	return teams, nil
}

// GetTeam возвращает команду по ID
func (r *Repository) GetTeam(ctx context.Context, teamID uuid.UUID) (*entity.Team, error) {
	const query = `
		SELECT id, name, owner_id, created_at, updated_at 
		FROM teams 
		WHERE id = $1
	`

	var team entity.Team
	err := r.Client.GetContext(ctx, &team, query, teamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, wrapper.Wrap(errors.New("team not found"))
		}
		r.Logger.ErrorContext(ctx, "failed to get team",
			slog.String("team_id", teamID.String()),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}

	return &team, nil
}

// GetTeamMembers возвращает всех участников команды
func (r *Repository) GetTeamMembers(ctx context.Context, teamID uuid.UUID) ([]entity.TeamMember, error) {
	const query = `
		SELECT tm.id, tm.team_id, tm.user_id, tm.role, tm.invited_by, tm.joined_at,
		       u.id, u.email, u.balance, u.status
		FROM team_members tm
		INNER JOIN users u ON tm.user_id = u.id
		WHERE tm.team_id = $1
		ORDER BY tm.joined_at ASC
	`

	rows, err := r.Client.QueryContext(ctx, query, teamID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to get team members",
			slog.String("team_id", teamID.String()),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}
	defer rows.Close()

	var members []entity.TeamMember
	for rows.Next() {
		var m entity.TeamMember
		var u entity.User
		var invitedBy sql.NullString

		err := rows.Scan(
			&m.ID, &m.TeamID, &m.UserID, &m.Role,
			&invitedBy, &m.JoinedAt,
			&u.ID, &u.Email, &u.Balance, &u.Status,
		)
		if err != nil {
			r.Logger.ErrorContext(ctx, "failed to scan team member", slog.Any("error", err))
			continue
		}

		if invitedBy.Valid {
			m.InvitedBy = &invitedBy.String
		}

		m.User = &u
		members = append(members, m)
	}

	return members, nil
}

// InviteTeamMember приглашает пользователя в команду по email
func (r *Repository) InviteTeamMember(ctx context.Context, teamID, inviterID uuid.UUID, email, role string) error {
	// Проверка прав приглашающего
	hasAccess, inviterRole, err := r.CheckTeamAccess(ctx, teamID, inviterID)
	if err != nil || !hasAccess {
		return wrapper.Wrap(errors.New("access denied"))
	}
	if inviterRole != "owner" && inviterRole != "admin" {
		return wrapper.Wrap(errors.New("insufficient permissions: only owner or admin can invite"))
	}

	// Находим пользователя по email
	user, err := r.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return wrapper.Wrap(fmt.Errorf("user with email %s not found", email))
	}

	// Проверяем, не состоит ли уже в команде
	var existingID uuid.UUID
	err = r.Client.GetContext(ctx, &existingID,
		"SELECT id FROM team_members WHERE team_id = $1 AND user_id = $2 LIMIT 1",
		teamID, user.ID)
	if err == nil {
		return nil // уже участник — не ошибка
	}

	if role != "admin" && role != "member" {
		return wrapper.Wrap(errors.New("invalid role: must be 'admin' or 'member'"))
	}

	const query = `
		INSERT INTO team_members (team_id, user_id, role, invited_by)
		VALUES ($1, $2, $3, $4)
	`

	_, err = r.Client.ExecContext(ctx, query, teamID, user.ID, role, inviterID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to invite team member",
			slog.String("team_id", teamID.String()),
			slog.String("email", email),
			slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	r.Logger.InfoContext(ctx, "user invited to team",
		slog.String("team_id", teamID.String()),
		slog.String("email", email),
		slog.String("role", role))

	return nil
}

// RemoveTeamMember удаляет участника из команды
func (r *Repository) RemoveTeamMember(ctx context.Context, teamID, userID, removerID uuid.UUID) error {
	hasAccess, removerRole, err := r.CheckTeamAccess(ctx, teamID, removerID)
	if err != nil || !hasAccess {
		return wrapper.Wrap(errors.New("access denied"))
	}
	if removerRole != "owner" && removerRole != "admin" {
		return wrapper.Wrap(errors.New("insufficient permissions"))
	}

	// Нельзя удалить owner'а
	var memberRole string
	err = r.Client.GetContext(ctx, &memberRole,
		"SELECT role FROM team_members WHERE team_id = $1 AND user_id = $2",
		teamID, userID)
	if err != nil {
		return wrapper.Wrap(errors.New("member not found"))
	}
	if memberRole == "owner" {
		return wrapper.Wrap(errors.New("cannot remove team owner"))
	}

	_, err = r.Client.ExecContext(ctx,
		"DELETE FROM team_members WHERE team_id = $1 AND user_id = $2",
		teamID, userID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to remove team member",
			slog.String("team_id", teamID.String()),
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	r.Logger.InfoContext(ctx, "user removed from team",
		slog.String("team_id", teamID.String()),
		slog.String("user_id", userID.String()))

	return nil
}

// UpdateTeamMemberRole изменяет роль участника
func (r *Repository) UpdateTeamMemberRole(ctx context.Context, teamID, userID uuid.UUID, newRole string, updaterID uuid.UUID) error {
	hasAccess, updaterRole, err := r.CheckTeamAccess(ctx, teamID, updaterID)
	if err != nil || !hasAccess {
		return wrapper.Wrap(errors.New("access denied"))
	}
	if updaterRole != "owner" {
		return wrapper.Wrap(errors.New("only owner can change roles"))
	}

	if newRole != "admin" && newRole != "member" {
		return wrapper.Wrap(errors.New("invalid role"))
	}

	var currentRole string
	err = r.Client.GetContext(ctx, &currentRole,
		"SELECT role FROM team_members WHERE team_id = $1 AND user_id = $2",
		teamID, userID)
	if err != nil {
		return wrapper.Wrap(errors.New("member not found"))
	}
	if currentRole == "owner" {
		return wrapper.Wrap(errors.New("cannot change owner role"))
	}

	_, err = r.Client.ExecContext(ctx,
		"UPDATE team_members SET role = $1 WHERE team_id = $2 AND user_id = $3",
		newRole, teamID, userID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to update team member role",
			slog.String("team_id", teamID.String()),
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	return nil
}

// CheckTeamAccess проверяет доступ пользователя к команде
func (r *Repository) CheckTeamAccess(ctx context.Context, teamID, userID uuid.UUID) (bool, string, error) {
	const query = `
		SELECT role FROM team_members 
		WHERE team_id = $1 AND user_id = $2
	`

	var role string
	err := r.Client.GetContext(ctx, &role, query, teamID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, "", nil
		}
		return false, "", wrapper.Wrap(err)
	}

	return true, role, nil
}
