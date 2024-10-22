package caa

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type QueueRepository struct {
	db *gorm.DB
}

func NewQueueRepository(db *gorm.DB) *QueueRepository {
	return &QueueRepository{
		db: db,
	}
}

func (r *QueueRepository) GetPendingQueues() ([]Queue, error) {
	var queueDBs []QueueDB
	result := r.db.Where("agent_id IS NULL AND resolve_at IS NULL").
		Order("create_at").
		Find(&queueDBs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch pending queues: %w", result.Error)
	}

	queues := make([]Queue, len(queueDBs))
	for i, queueDB := range queueDBs {
		queues[i] = Queue{
			RoomID:   queueDB.RoomID,
			AgentID:  queueDB.AgentID,
			CreateAt: queueDB.CreateAt,
		}
	}

	return queues, nil
}

func (r *QueueRepository) CreateQueue(roomID int) error {
	log.Info().Msgf("roomID %v", roomID)
	queue := &QueueDB{
		RoomID:   roomID,
		AgentID:  0,
		CreateAt: time.Now(),
	}
	result := r.db.Create(queue)
	if result.Error != nil {
		return fmt.Errorf("failed to create queue: %w", result.Error)
	}

	return nil
}

func (r *QueueRepository) AssignAgent(roomID int, agentID int) error {
	result := r.db.Model(&QueueDB{}).Where("room_id = ?", roomID).Updates(QueueDB{
		AgentID: agentID,
	})

	if result.Error != nil {
		return fmt.Errorf("failed to assign agent to queue: %w", result.Error)
	}

	return nil
}

func (r *QueueRepository) ResolveQueue(roomID int) error {
	result := r.db.Model(&QueueDB{}).Where("room_id = ?", roomID).Updates(QueueDB{
		ResolveAt: time.Now(),
	})

	if result.Error != nil {
		return fmt.Errorf("failed to resolve queue: %w", result.Error)
	}

	return nil
}
