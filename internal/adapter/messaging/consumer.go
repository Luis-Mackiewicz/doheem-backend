package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"doheem-backend/internal/domain"

	"github.com/segmentio/kafka-go"
)

type memberLister interface {
	ListMembers(ctx context.Context, groupID string) ([]domain.GroupMemberWithUser, error)
}

type notificationCreator interface {
	Create(ctx context.Context, params domain.CreateNotificationParams) (domain.Notification, error)
}

type ConsumerDeps struct {
	NotifSvc  notificationCreator
	MemberSvc memberLister
}

func StartConsumer(ctx context.Context, brokers string, deps ConsumerDeps) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{brokers},
		Topic:     "doheem.events",
		GroupID:   "doheem-notifications",
		MinBytes:  1,
		MaxBytes:  10e6,
	})

	go func() {
		defer reader.Close()
		slog.Info("kafka consumer started")

		for {
			select {
			case <-ctx.Done():
				slog.Info("kafka consumer stopped")
				return
			default:
			}

			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				slog.Error("kafka consumer read error", "error", err)
				continue
			}

			var event domain.DomainEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				slog.Error("kafka consumer unmarshal error", "error", err)
				continue
			}

			processEvent(ctx, deps, event)
		}
	}()
}

func processEvent(ctx context.Context, deps ConsumerDeps, event domain.DomainEvent) {
	switch event.Type {
	case "expense.created":
		members, err := deps.MemberSvc.ListMembers(ctx, event.GroupID)
		if err != nil {
			slog.Error("failed to list group members for notification", "error", err)
			return
		}
		for _, m := range members {
			deps.NotifSvc.Create(ctx, domain.CreateNotificationParams{
				UserID:  m.UserID,
				GroupID: &event.GroupID,
				Type:    "expense_added",
				Title:   "Nova despesa",
				Message: "Uma nova despesa foi adicionada ao grupo.",
			})
		}

	case "payment.confirmed":
		receiverID, _ := event.Payload["receiver_id"].(string)
		payerName, _ := event.Payload["payer_name"].(string)
		amount, _ := event.Payload["amount"].(float64)

		if receiverID != "" {
			msg := payerName + " confirmou um pagamento de R$ " + formatAmount(amount)
			deps.NotifSvc.Create(ctx, domain.CreateNotificationParams{
				UserID:  receiverID,
				GroupID: &event.GroupID,
				Type:    "payment_confirmed",
				Title:   "Pagamento confirmado",
				Message: msg,
			})
		}

	case "task.occurrence.completed":
		assignedTo, _ := event.Payload["assigned_to"].(string)
		taskTitle, _ := event.Payload["task_title"].(string)

		if assignedTo != "" {
			deps.NotifSvc.Create(ctx, domain.CreateNotificationParams{
				UserID:  assignedTo,
				GroupID: &event.GroupID,
				Type:    "task_completed",
				Title:   "Tarefa concluída",
				Message: "A tarefa \"" + taskTitle + "\" foi concluída.",
			})
		}

	default:
		slog.Debug("unknown event type", "type", event.Type)
	}
}

func formatAmount(v float64) string {
	return fmt.Sprintf("%.2f", v)
}
