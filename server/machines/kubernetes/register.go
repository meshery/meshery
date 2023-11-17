package kubernetes

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshkit/models/events"
)

type RegisterAction struct {}

// Execute On Entry and Exit should not return next eventtype i suppose, look again.
func(ra *RegisterAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	
	eventBuilder := events.NewEvent().ActedUpon(uuid.Nil).WithCategory("connection").WithAction("register").FromSystem(uuid.Nil).FromUser(uuid.Nil) // pass userID and systemID in acted upon first pass user id if we can get context then update with connection Id
	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		return machines.NoOp, eventBuilder.Build(), err
	}
	
	err = AssignClientSetToContext(machinectx, eventBuilder)
	if err != nil {
		return machines.NotFound, eventBuilder.Build(), err
	}
	
	return machines.NoOp, eventBuilder.Build(), nil
}

func(ra *RegisterAction) Execute(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {

	eventBuilder := events.NewEvent().ActedUpon(uuid.Nil).WithCategory("connection").WithAction("register").FromSystem(uuid.Nil).FromUser(uuid.Nil) // pass userID and systemID in acted upon first pass user id if we can get context then update with connection Id

	
	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		return machines.NoOp, eventBuilder.Build(), err
	}
	
	fmt.Println("inside register line 43---------------", machineCtx)
	err = machinectx.K8sContext.PingTest()
	fmt.Println("inside register line 45---------------", err)
	if err != nil {
		// machinectx.log.Error(err)
		fmt.Println("inside register line 48---------------", err)
		
		// peform error handling and event publishing
		return machines.NotFound, nil, err
	}
	token, _ := ctx.Value(models.TokenCtxKey).(string)
	
	connection, statusCode, err := machinectx.Provider.UpdateConnectionStatusByID(token, uuid.FromStringOrNil(machinectx.K8sContext.ConnectionID), connections.REGISTERED)

	// peform error handling and event publishing
	if err != nil {
		return machines.NoOp, eventBuilder.Build(), err
	}

	fmt.Println("connection inside register.go", connection, statusCode)

	return machines.Connect, nil, nil
}

func(ra *RegisterAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
