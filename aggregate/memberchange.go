package aggregate

import (
	"fmt"

	"github.com/sorintlab/sircles/command/commands"
	"github.com/sorintlab/sircles/common"
	ep "github.com/sorintlab/sircles/events"
	"github.com/sorintlab/sircles/eventstore"
	"github.com/sorintlab/sircles/models"
	"github.com/sorintlab/sircles/util"
)

type MemberChangeRepository struct {
	es           *eventstore.EventStore
	uidGenerator common.UIDGenerator
}

func NewMemberChangeRepository(es *eventstore.EventStore, uidGenerator common.UIDGenerator) *MemberChangeRepository {
	return &MemberChangeRepository{es: es, uidGenerator: uidGenerator}
}

func NewMemberChangeRepositoryDisable(es *eventstore.EventStore, uidGenerator common.UIDGenerator) *MemberChangeRepository {
	return &MemberChangeRepository{es: es, uidGenerator: uidGenerator}
}

func NewMemberChangeRepositoryActivate(es *eventstore.EventStore, uidGenerator common.UIDGenerator) *MemberChangeRepository {
	return &MemberChangeRepository{es: es, uidGenerator: uidGenerator}
}

func (mr *MemberChangeRepository) Load(id util.ID) (*MemberChange, error) {
	log.Debugf("Load id: %s", id)
	m, err := NewMemberChange(mr.es, mr.uidGenerator, id)
	if err != nil {
		return nil, err
	}

	if err := batchLoader(mr.es, id.String(), m); err != nil {
		return nil, err
	}

	return m, nil
}

type MemberChange struct {
	id      util.ID
	version int64

	completed bool

	es           *eventstore.EventStore
	uidGenerator common.UIDGenerator
}

func NewMemberChange(es *eventstore.EventStore, uidGenerator common.UIDGenerator, id util.ID) (*MemberChange, error) {
	return &MemberChange{
		id:           id,
		es:           es,
		uidGenerator: uidGenerator,
	}, nil
}

func (m *MemberChange) Version() int64 {
	return m.version
}

func (m *MemberChange) ID() string {
	return m.id.String()
}

func (m *MemberChange) AggregateType() AggregateType {
	return MemberChangeAggregate
}

func (m *MemberChange) HandleCommand(command *commands.Command) ([]ep.Event, error) {
	var events []ep.Event
	var err error

	// skip if already completed
	if m.completed {
		return events, nil
	}

	switch command.CommandType {
	case commands.CommandTypeRequestCreateMember:
		events, err = m.HandleRequestCreateMemberCommand(command)
	case commands.CommandTypeRequestUpdateMember:
		events, err = m.HandleRequestUpdateMemberCommand(command)
	case commands.CommandTypeRequestUpdateMemberDisable:
		events, err = m.HandleRequestUpdateMemberCommandDisable(command)
	case commands.CommandTypeRequestUpdateActivateMember:
		events, err = m.HandleRequestUpdateActivateMemberCommand(command)
	case commands.CommandTypeRequestSetMemberMatchUID:
		events, err = m.HandleRequestSetMemberMatchUIDCommand(command)
	case commands.CommandTypeCompleteRequest:
		events, err = m.HandleCompleteRequestCommand(command)

	default:
		err = fmt.Errorf("unhandled command: %#v", command)
	}

	return events, err
}

func (m *MemberChange) HandleRequestCreateMemberCommand(command *commands.Command) ([]ep.Event, error) {
	events := []ep.Event{}

	c := command.Data.(*commands.RequestCreateMember)

	member := &models.Member{
		UserName: c.UserName,
		FullName: c.FullName,
		Email:    c.Email,
		IsAdmin:  c.IsAdmin,
	}
	member.ID = c.MemberID

	events = append(events, ep.NewEventMemberChangeCreateRequested(m.id, member, c.MatchUID, c.PasswordHash, c.Avatar))

	return events, nil
}

func (m *MemberChange) HandleRequestUpdateMemberCommand(command *commands.Command) ([]ep.Event, error) {
	events := []ep.Event{}

	c := command.Data.(*commands.RequestUpdateMember)

	member := &models.Member{
		UserName: c.UserName,
		FullName: c.FullName,
		Email:    c.Email,
		IsAdmin:  c.IsAdmin,
	}
	member.ID = c.MemberID

	events = append(events, ep.NewEventMemberChangeUpdateRequested(m.id, member, c.Avatar, c.PrevUserName, c.PrevEmail))

	return events, nil
}

func (m *MemberChange) HandleRequestUpdateActivateMemberCommand(command *commands.Command) ([]ep.Event, error) {
	events := []ep.Event{}

	c := command.Data.(*commands.RequestUpdateActivateMember)

	member := &models.Member{
	}
	member.ID = c.MemberID

	events = append(events, ep.NewEventMemberChangeUpdateActivateRequested(m.id, member))

	return events, nil
}

func (m *MemberChange) HandleRequestUpdateMemberCommandDisable(command *commands.Command) ([]ep.Event, error) {
	events := []ep.Event{}

	c := command.Data.(*commands.RequestUpdateMemberDisable)

	member := &models.Member{
	}
	member.ID = c.MemberID

	events = append(events, ep.NewEventMemberChangeUpdateRequestedDisable(m.id, member))

	return events, nil
}

func (m *MemberChange) HandleRequestSetMemberMatchUIDCommand(command *commands.Command) ([]ep.Event, error) {
	events := []ep.Event{}

	c := command.Data.(*commands.RequestSetMemberMatchUID)

	events = append(events, ep.NewEventMemberChangeSetMatchUIDRequested(m.id, c.MemberID, c.MatchUID))

	return events, nil
}

func (m *MemberChange) HandleCompleteRequestCommand(command *commands.Command) ([]ep.Event, error) {
	events := []ep.Event{}

	c := command.Data.(*commands.CompleteRequest)

	events = append(events, ep.NewEventMemberChangeCompleted(m.id, c.Error, c.Reason))

	return events, nil
}

func (m *MemberChange) ApplyEvents(events []*eventstore.StoredEvent) error {
	for _, e := range events {
		if err := m.ApplyEvent(e); err != nil {
			return err
		}
	}
	return nil
}

func (m *MemberChange) ApplyEvent(event *eventstore.StoredEvent) error {
	log.Debugf("event: %v", event)

	m.version = event.Version

	switch ep.EventType(event.EventType) {
	case ep.EventTypeMemberChangeCompleted:
		m.completed = true

	}

	return nil
}
