package agent

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

type Service struct {
	repo *AgentRepository
}

func NewService(repo *AgentRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAvailableAgents() ([]Agent, error) {
	remoteAgents, err := s.repo.FetchAgentsFromOmnichannel()
	if err != nil {
		return nil, fmt.Errorf("fetch agents from remote: %w", err)
	}

	onlineAgents := make(map[int]struct{}, len(remoteAgents))
	for _, remoteAgent := range remoteAgents {
		agent := Agent{
			ID:       remoteAgent.ID,
			Name:     remoteAgent.Name,
			Email:    remoteAgent.Email,
			IsOnline: remoteAgent.IsOnline,
		}
		if err := s.repo.CreateOrUpdateAgentPostgres(agent); err != nil {
			log.Error().Err(err).Int("agent_id", agent.ID).Msg("Failed to save or update agent")
			continue
		}
		if agent.IsOnline {
			onlineAgents[agent.ID] = struct{}{}
		}
	}

	localAgents, err := s.repo.FetchAgentsFromPostgres()
	if err != nil {
		return nil, fmt.Errorf("fetch agents from DB: %w", err)
	}

	availableAgents := make([]Agent, 0, len(localAgents))
	for _, localAgent := range localAgents {
		if _, isOnline := onlineAgents[localAgent.ID]; isOnline && localAgent.CustomerCount <= 2 {
			availableAgents = append(availableAgents, Agent{
				ID:            localAgent.ID,
				Name:          localAgent.Name,
				Email:         localAgent.Email,
				IsOnline:      true,
				CustomerCount: localAgent.CustomerCount,
			})
		}
	}

	// for _, k := range remoteAgents {
	// 	log.Info().Msgf("%v", k)
	// }

	// for _, k := range localAgents {
	// 	log.Debug().Msgf("%v", k)
	// }

	for _, k := range availableAgents {
		log.Info().Msgf("%v", k)
	}

	return availableAgents, nil
}

func (s *Service) AssignAgentsToRooms(allocations map[int]int) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(allocations))

	for roomID, agentID := range allocations {
		wg.Add(1)
		go func(roomID, agentID int) {
			defer wg.Done()
			if err := s.repo.AssignAgentToRoom(roomID, agentID); err != nil {
				errChan <- fmt.Errorf("failed to assign agent %d to room %d: %w", agentID, roomID, err)
			}
		}(roomID, agentID)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred during agent assignment: %v", errs)
	}

	return nil
}

func (s *Service) AssignAgentToRoom(roomID int, agentID int) error {
	return s.repo.AssignAgentToRoom(roomID, agentID)
}
