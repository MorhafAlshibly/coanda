package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type PollMatchmakingTicketCommand struct {
	service *Service
	In      *api.MatchmakingTicketRequest
	Out     *api.MatchmakingTicketResponse
}

func NewPollMatchmakingTicketCommand(service *Service, in *api.MatchmakingTicketRequest) *PollMatchmakingTicketCommand {
	return &PollMatchmakingTicketCommand{
		service: service,
		In:      in,
	}
}

func (c *PollMatchmakingTicketCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
