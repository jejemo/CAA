package agent

import (
	"fmt"
	"go-caa/internal/api/client"
	"net/http"

	"gorm.io/gorm"
)

type AgentRepository struct {
	db    *gorm.DB
	qismo *client.QismoClient
}

func NewRepository(db *gorm.DB, qismo *client.QismoClient) *AgentRepository {
	return &AgentRepository{
		db:    db,
		qismo: qismo,
	}
}

func (r *AgentRepository) AssignAgentToRoom(roomID int, agentID int) error {
	formData := map[string]string{
		"room_id":  string(roomID),
		"agent_id": string(agentID),
	}

	var response AgentQismoResponse
	err := r.qismo.CallAPI(http.MethodPost, "/api/v1/admin/service/assign_agent", client.URLEncoded, formData, nil, &response)

	if err != nil {
		return err
	}

	return nil
}

func (r *AgentRepository) FetchAgentsFromOmnichannel() ([]Agent, error) {
	var response AgentQismoResponse
	err := r.qismo.CallAPI(http.MethodGet, "/api/v2/admin/agents", client.AplicationJSON, nil, nil, &response)

	if err != nil {
		return nil, err
	}

	agents := make([]Agent, 0, len(response.Data.Agents))
	for _, remoteAgent := range response.Data.Agents {
		agent := Agent{
			ID:       remoteAgent.ID,
			Name:     remoteAgent.Name,
			Email:    remoteAgent.Email,
			IsOnline: remoteAgent.IsAvailable,
		}

		agents = append(agents, agent)
	}

	return agents, nil
}

func (r *AgentRepository) FetchAgentsFromPostgres() ([]Agent, error) {
	var agentDBs []Agent
	result := r.db.Table("agent_dbs").
		Select("agent_dbs.id, agent_dbs.name, agent_dbs.email, COUNT(queue_dbs.room_id) as customer_count").
		Joins("LEFT JOIN queue_dbs ON agent_dbs.id = queue_dbs.agent_id").
		Group("agent_dbs.id, agent_dbs.name, agent_dbs.email").
		Having("COUNT(queue_dbs.room_id) < ?", 2).
		Order("agent_dbs.id").
		Find(&agentDBs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch agents from database: %w", result.Error)
	}

	return agentDBs, nil

}

func (r *AgentRepository) CreateOrUpdateAgentPostgres(agent Agent) error {
	var existingAgent AgentDB
	result := r.db.First(&existingAgent, agent.ID)

	if result.Error == nil {
		// Agent exists, update it
		return r.UpdateAgentPostgres(agent)
	} else if result.Error == gorm.ErrRecordNotFound {
		// Agent doesn't exist, save it
		return r.CreateAgentPostgres(agent)
	} else {
		// Other database error
		return fmt.Errorf("error checking agent existence: %w", result.Error)
	}
}

func (r *AgentRepository) CreateAgentPostgres(agent Agent) error {
	agentDB := &AgentDB{
		ID:    agent.ID,
		Name:  agent.Name,
		Email: agent.Email,
	}

	result := r.db.Create(agentDB)
	if result.Error != nil {
		return fmt.Errorf("failed to create agent: %w", result.Error)
	}

	return nil
}

func (r *AgentRepository) UpdateAgentPostgres(agent Agent) error {
	result := r.db.Model(&AgentDB{}).Where("id = ?", agent.ID).Updates(AgentDB{
		Name:  agent.Name,
		Email: agent.Email,
	})

	if result.Error != nil {
		return fmt.Errorf("failed to update agent: %w", result.Error)
	}

	return nil
}
