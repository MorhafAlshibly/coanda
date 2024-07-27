package matchmaking

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
)

type CreateMatchmakingTicketCommand struct {
	service *Service
	In      *api.CreateMatchmakingTicketRequest
	Out     *api.CreateMatchmakingTicketResponse
}

/*

dont create ticket if
 - a ticket exists that hasnt passed expireAt
 - a ticket has a match which hasnt passed endedAt
 - side note: a match that hasnt started in x time should update the endedAt to now


*/

func NewCreateMatchmakingTicketCommand(service *Service, in *api.CreateMatchmakingTicketRequest) *CreateMatchmakingTicketCommand {
	return &CreateMatchmakingTicketCommand{
		service: service,
		In:      in,
	}
}

func (c *CreateMatchmakingTicketCommand) Execute(ctx context.Context) error {
	panic("implement me")
}
