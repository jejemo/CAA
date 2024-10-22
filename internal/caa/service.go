package caa

import (
	"fmt"
	"go-caa/internal/agent"
	"strconv"
)

type Service struct {
	repo     *QueueRepository
	agentSvc *agent.Service
}

func NewService(repo *QueueRepository, agentSvc *agent.Service) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ProcessAllocateAgent(caa QismoCAAWebhook) error {
	roomID, err := strconv.Atoi(caa.RoomID)
	if err != nil {
		return fmt.Errorf("failed to parse room ID: %w", err)
	}

	if err := s.repo.CreateQueue(roomID); err != nil {
		return fmt.Errorf("failed to create queue for room %d: %w", roomID, err)
	}

	if err := s.AllocateAgentsToRooms(); err != nil {
		return fmt.Errorf("failed to allocate agent into room: %w", err)
	}

	return nil
}

func (s *Service) AllocateAgentsToRooms() error {
	rooms, err := s.repo.GetPendingQueues()
	if err != nil {
		return fmt.Errorf("failed to get pending queues: %w", err)
	}

	agents, err := s.agentSvc.GetAvailableAgents()
	if err != nil {
		return fmt.Errorf("failed to get available agents: %w", err)
	}

	if len(agents) == 0 {
		return fmt.Errorf("no agents available for allocation")
	}

	allocations := s.roundRobinAllocation(rooms, agents)

	if err := s.agentSvc.AssignAgentsToRooms(allocations); err != nil {
		return fmt.Errorf("failed to assign agents to rooms: %w", err)
	}

	return nil
}

// RoundRobinAllocation performs the round-robin allocation of agents to rooms
func (s *Service) roundRobinAllocation(rooms []Queue, agents []agent.Agent) map[int]int {
	allocations := make(map[int]int)
	agentQueue := agents

	for i := 0; i < len(rooms) && len(agentQueue) > 0; i++ {
		room := rooms[i]
		agent := agentQueue[0]

		allocations[room.RoomID] = agent.ID
		agent.CustomerCount++

		agentQueue = agentQueue[1:]

		if agent.CustomerCount < 2 {
			agentQueue = append(agentQueue, agent)
		}
	}

	return allocations
}

func (s *Service) ResolveRoom(room QismoMarkAsResolvedWebhook) error {
	roomID, err := strconv.Atoi(room.Service.RoomID)
	if err != nil {
		return fmt.Errorf("failed to parse room ID: %w", err)
	}
	
	return s.repo.ResolveQueue(roomID)
}
