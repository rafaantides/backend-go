package hooks

import (
	"context"
	"fmt"

	"backend-go/pkg/ent"
	"backend-go/pkg/ent/paymentstatus"
)

func SetDefaultStatusHook(client *ent.Client) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			dm, ok := m.(*ent.DebtMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type: %T", m)
			}

			if _, exists := dm.StatusID(); exists {
				return next.Mutate(ctx, m)
			}

			status, err := client.PaymentStatus.
				Query().
				Where(paymentstatus.NameEQ("pending")).
				Only(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to find 'pending' status: %w", err)
			}

			dm.SetStatusID(status.ID)
			return next.Mutate(ctx, dm)
		})
	}
}
