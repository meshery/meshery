package kubernetes

import (
	"context"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/core"
)

func TestStateMachines_PanicFixes(t *testing.T) {
	// actionsWithProvider are actions that validate user, sysID, AND provider
	actionsWithProvider := []struct {
		name    string
		execute func(context.Context) error
	}{
		{
			name: "ConnectAction",
			execute: func(ctx context.Context) error {
				ca := &ConnectAction{}
				_, _, err := ca.Execute(ctx, nil, nil)
				return err
			},
		},
		{
			name: "RegisterAction",
			execute: func(ctx context.Context) error {
				ra := &RegisterAction{}
				_, _, err := ra.Execute(ctx, nil, nil)
				return err
			},
		},
		{
			name: "DeleteAction",
			execute: func(ctx context.Context) error {
				da := &DeleteAction{}
				_, _, err := da.Execute(ctx, nil, nil)
				return err
			},
		},
		{
			name: "DiscoverAction",
			execute: func(ctx context.Context) error {
				da := &DiscoverAction{}
				_, _, err := da.Execute(ctx, nil, nil)
				return err
			},
		},
	}

	// DisconnectAction only validates user and sysID (no provider check)
	actionsWithoutProvider := []struct {
		name    string
		execute func(context.Context) error
	}{
		{
			name: "DisconnectAction",
			execute: func(ctx context.Context) error {
				da := &DisconnectAction{}
				_, _, err := da.Execute(ctx, nil, nil)
				return err
			},
		},
	}

	for _, action := range actionsWithProvider {
		t.Run(action.name, func(t *testing.T) {
			ctx := context.Background()

			// Should return "user missing from context" without panicking
			err := action.execute(ctx)
			if err == nil {
				t.Errorf("expected error due to missing context values, got nil")
			} else if err.Error() != "user missing from context" {
				t.Errorf("unexpected error: %v", err)
			}

			// Add user, but no sysID
			user := &models.User{ID: core.Uuid{}}
			ctxWithUser := context.WithValue(ctx, models.UserCtxKey, user)

			err = action.execute(ctxWithUser)
			if err == nil {
				t.Errorf("expected error due to missing context values, got nil")
			} else if err.Error() != "system ID missing from context" {
				t.Errorf("unexpected error: %v", err)
			}

			// Add user and sysID, but no provider
			sysID := &core.Uuid{}
			ctxWithUserAndSysID := context.WithValue(ctxWithUser, models.SystemIDKey, sysID)

			err = action.execute(ctxWithUserAndSysID)
			if err == nil {
				t.Errorf("expected error due to missing provider, got nil")
			} else if err.Error() != "provider missing from context" {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}

	for _, action := range actionsWithoutProvider {
		t.Run(action.name, func(t *testing.T) {
			ctx := context.Background()

			// Should return "user missing from context" without panicking
			err := action.execute(ctx)
			if err == nil {
				t.Errorf("expected error due to missing context values, got nil")
			} else if err.Error() != "user missing from context" {
				t.Errorf("unexpected error: %v", err)
			}

			// Add user, but no sysID
			user := &models.User{ID: core.Uuid{}}
			ctxWithUser := context.WithValue(ctx, models.UserCtxKey, user)

			err = action.execute(ctxWithUser)
			if err == nil {
				t.Errorf("expected error due to missing context values, got nil")
			} else if err.Error() != "system ID missing from context" {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestAssignInitialCtx_PanicFix(t *testing.T) {
	ctx := context.Background()

	// Should return "user missing from context" without panicking
	_, _, err := AssignInitialCtx(ctx, nil, nil)
	if err == nil {
		t.Errorf("expected error due to missing context values, got nil")
	} else if err.Error() != "user missing from context" {
		t.Errorf("unexpected error: %v", err)
	}

	// Add user and sysID, but no provider
	user := &models.User{ID: core.Uuid{}}
	sysID := &core.Uuid{}
	ctxWithUserAndSysID := context.WithValue(
		context.WithValue(ctx, models.UserCtxKey, user),
		models.SystemIDKey, sysID,
	)

	_, _, err = AssignInitialCtx(ctxWithUserAndSysID, nil, nil)
	if err == nil {
		t.Errorf("expected error due to missing provider, got nil")
	} else if err.Error() != "provider missing from context" {
		t.Errorf("unexpected error: %v", err)
	}
}
