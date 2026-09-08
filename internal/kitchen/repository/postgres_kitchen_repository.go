package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"pos-backend/internal/database"
	"pos-backend/internal/kitchen/domain"
	kitchendb "pos-backend/internal/kitchen/repository/generated"
	"pos-backend/internal/pkg/pgconv"
)

type PostgresKitchenRepository struct {
	queries   *kitchendb.Queries
	txManager *database.TransactionManager
}

func NewPostgresKitchenRepository(
	queries *kitchendb.Queries,
	txManager *database.TransactionManager,
) *PostgresKitchenRepository {
	return &PostgresKitchenRepository{
		queries:   queries,
		txManager: txManager,
	}
}

func (r *PostgresKitchenRepository) Create(
	ctx context.Context,
	ticket *domain.KitchenTicket,
) error {
	return r.txManager.WithinTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		txQueries := r.queries.WithTx(tx)

		createdTicketRow, err := txQueries.CreateKitchenTicket(
			ctx,
			kitchendb.CreateKitchenTicketParams{
				ID:          uuidToPg(ticket.ID),
				CompanyID:   uuidToPg(ticket.CompanyID),
				StoreID:     uuidToPg(ticket.StoreID),
				OrderID:     uuidToPg(ticket.OrderID),
				OrderNumber: ticket.OrderNumber,
				OrderType:   ticket.OrderType,
				TableID:     optUuidToPg(ticket.TableID),
				Status:      string(ticket.Status),
				Priority:    string(ticket.Priority),
				Notes:       textToPg(ticket.Notes),
				CreatedAt:   timeToPg(ticket.CreatedAt),
				UpdatedAt:   timeToPg(ticket.UpdatedAt),
			},
		)
		if err != nil {
			return err
		}

		for _, item := range ticket.Items {
			_, err := txQueries.CreateKitchenTicketItem(
				ctx,
				kitchendb.CreateKitchenTicketItemParams{
					ID:          uuidToPg(item.ID),
					TicketID:    uuidToPg(ticket.ID),
					OrderItemID: optUuidToPg(item.OrderItemID),
					MenuItemID:  uuidToPg(item.MenuItemID),
					ItemName:    item.ItemName,
					Sku:         item.SKU,
					Quantity:    int64ToNumeric(item.Quantity),
					Notes:       textToPg(item.Notes),
					Status:      string(item.Status),
					CreatedAt:   timeToPg(item.CreatedAt),
					UpdatedAt:   timeToPg(item.UpdatedAt),
				},
			)
			if err != nil {
				return err
			}
		}

		_ = createdTicketRow
		return nil
	})
}

func (r *PostgresKitchenRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.KitchenTicket, error) {
	row, err := r.queries.GetKitchenTicketByID(ctx, uuidToPg(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTicketNotFound
		}
		return nil, err
	}

	itemsRows, err := r.queries.ListKitchenTicketItemsByTicketID(ctx, uuidToPg(id))
	if err != nil {
		return nil, err
	}

	items := make([]domain.KitchenTicketItem, 0, len(itemsRows))
	for _, ir := range itemsRows {
		items = append(items, toDomainItem(ir))
	}

	ticket := toDomainTicketFromGetByID(row)
	ticket.Items = items
	return &ticket, nil
}

func (r *PostgresKitchenRepository) GetByOrderID(
	ctx context.Context,
	orderID uuid.UUID,
) (*domain.KitchenTicket, error) {
	row, err := r.queries.GetKitchenTicketByOrderID(ctx, uuidToPg(orderID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTicketNotFound
		}
		return nil, err
	}

	itemsRows, err := r.queries.ListKitchenTicketItemsByTicketID(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	items := make([]domain.KitchenTicketItem, 0, len(itemsRows))
	for _, ir := range itemsRows {
		items = append(items, toDomainItem(ir))
	}

	ticket := toDomainTicketFromKitchenTicket(row)
	ticket.Items = items
	return &ticket, nil
}

func (r *PostgresKitchenRepository) ListByStore(
	ctx context.Context,
	storeID uuid.UUID,
	status *string,
) ([]domain.KitchenTicket, error) {
	var statusParam pgtype.Text
	if status != nil && *status != "" && *status != "ALL" {
		statusParam = pgtype.Text{String: *status, Valid: true}
	}

	rows, err := r.queries.ListKitchenTicketsByStore(
		ctx,
		kitchendb.ListKitchenTicketsByStoreParams{
			StoreID: uuidToPg(storeID),
			Status:  statusParam,
		},
	)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return []domain.KitchenTicket{}, nil
	}

	ticketIDs := make([]pgtype.UUID, 0, len(rows))
	for _, row := range rows {
		ticketIDs = append(ticketIDs, row.ID)
	}

	itemsRows, err := r.queries.ListKitchenTicketItemsByTicketIDs(ctx, ticketIDs)
	if err != nil {
		return nil, err
	}

	itemsByTicket := make(map[uuid.UUID][]domain.KitchenTicketItem)
	for _, ir := range itemsRows {
		tid := uuid.UUID(ir.TicketID.Bytes)
		itemsByTicket[tid] = append(itemsByTicket[tid], toDomainItem(ir))
	}

	tickets := make([]domain.KitchenTicket, 0, len(rows))
	for _, row := range rows {
		t := toDomainTicketFromListRow(row)
		t.Items = itemsByTicket[t.ID]
		tickets = append(tickets, t)
	}

	return tickets, nil
}

func (r *PostgresKitchenRepository) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status domain.TicketStatus,
) (*domain.KitchenTicket, error) {
	row, err := r.queries.UpdateKitchenTicketStatus(
		ctx,
		kitchendb.UpdateKitchenTicketStatusParams{
			ID:     uuidToPg(id),
			Status: string(status),
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTicketNotFound
		}
		return nil, err
	}

	// Update ticket items status as well
	_ = r.queries.UpdateKitchenTicketItemsStatusByTicketID(
		ctx,
		kitchendb.UpdateKitchenTicketItemsStatusByTicketIDParams{
			TicketID: uuidToPg(id),
			Status:   string(status),
		},
	)

	itemsRows, err := r.queries.ListKitchenTicketItemsByTicketID(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	items := make([]domain.KitchenTicketItem, 0, len(itemsRows))
	for _, ir := range itemsRows {
		items = append(items, toDomainItem(ir))
	}

	ticket := toDomainTicketFromKitchenTicket(row)
	ticket.Items = items
	return &ticket, nil
}

func toDomainTicketFromGetByID(row kitchendb.GetKitchenTicketByIDRow) domain.KitchenTicket {
	return domain.KitchenTicket{
		ID:          uuid.UUID(row.ID.Bytes),
		CompanyID:   uuid.UUID(row.CompanyID.Bytes),
		StoreID:     uuid.UUID(row.StoreID.Bytes),
		OrderID:     uuid.UUID(row.OrderID.Bytes),
		OrderNumber: row.OrderNumber,
		OrderType:   row.OrderType,
		TableID:     pgToOptUuid(row.TableID),
		TableNumber: pgToOptText(row.TableNumber),
		Status:      domain.TicketStatus(row.Status),
		Priority:    domain.TicketPriority(row.Priority),
		Notes:       pgToOptText(row.Notes),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
		StartedAt:   pgToOptTime(row.StartedAt),
		ReadyAt:     pgToOptTime(row.ReadyAt),
		ServedAt:    pgToOptTime(row.ServedAt),
		Items:       []domain.KitchenTicketItem{},
	}
}

func toDomainTicketFromListRow(row kitchendb.ListKitchenTicketsByStoreRow) domain.KitchenTicket {
	return domain.KitchenTicket{
		ID:          uuid.UUID(row.ID.Bytes),
		CompanyID:   uuid.UUID(row.CompanyID.Bytes),
		StoreID:     uuid.UUID(row.StoreID.Bytes),
		OrderID:     uuid.UUID(row.OrderID.Bytes),
		OrderNumber: row.OrderNumber,
		OrderType:   row.OrderType,
		TableID:     pgToOptUuid(row.TableID),
		TableNumber: pgToOptText(row.TableNumber),
		Status:      domain.TicketStatus(row.Status),
		Priority:    domain.TicketPriority(row.Priority),
		Notes:       pgToOptText(row.Notes),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
		StartedAt:   pgToOptTime(row.StartedAt),
		ReadyAt:     pgToOptTime(row.ReadyAt),
		ServedAt:    pgToOptTime(row.ServedAt),
		Items:       []domain.KitchenTicketItem{},
	}
}

func toDomainTicketFromKitchenTicket(row kitchendb.KitchenTicket) domain.KitchenTicket {
	return domain.KitchenTicket{
		ID:          uuid.UUID(row.ID.Bytes),
		CompanyID:   uuid.UUID(row.CompanyID.Bytes),
		StoreID:     uuid.UUID(row.StoreID.Bytes),
		OrderID:     uuid.UUID(row.OrderID.Bytes),
		OrderNumber: row.OrderNumber,
		OrderType:   row.OrderType,
		TableID:     pgToOptUuid(row.TableID),
		Status:      domain.TicketStatus(row.Status),
		Priority:    domain.TicketPriority(row.Priority),
		Notes:       pgToOptText(row.Notes),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
		StartedAt:   pgToOptTime(row.StartedAt),
		ReadyAt:     pgToOptTime(row.ReadyAt),
		ServedAt:    pgToOptTime(row.ServedAt),
		Items:       []domain.KitchenTicketItem{},
	}
}

func toDomainItem(row kitchendb.KitchenTicketItem) domain.KitchenTicketItem {
	return domain.KitchenTicketItem{
		ID:          uuid.UUID(row.ID.Bytes),
		TicketID:    uuid.UUID(row.TicketID.Bytes),
		OrderItemID: pgToOptUuid(row.OrderItemID),
		MenuItemID:  uuid.UUID(row.MenuItemID.Bytes),
		ItemName:    row.ItemName,
		SKU:         row.Sku,
		Quantity:    numericToInt64(row.Quantity),
		Notes:       pgToOptText(row.Notes),
		Status:      domain.TicketStatus(row.Status),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}

func uuidToPg(u uuid.UUID) pgtype.UUID            { return pgconv.UUID(u) }
func optUuidToPg(u *uuid.UUID) pgtype.UUID        { return pgconv.OptUUID(u) }
func pgToOptUuid(p pgtype.UUID) *uuid.UUID        { return pgconv.ToOptUUID(p) }
func textToPg(s *string) pgtype.Text              { return pgconv.OptText(s) }
func pgToOptText(p pgtype.Text) *string           { return pgconv.ToOptText(p) }
func timeToPg(t time.Time) pgtype.Timestamptz     { return pgconv.Timestamptz(t) }
func pgToOptTime(p pgtype.Timestamptz) *time.Time { return pgconv.ToOptTime(p) }
func int64ToNumeric(value int64) pgtype.Numeric   { return pgconv.Int64ToNumeric(value) }
func numericToInt64(value pgtype.Numeric) int64   { return pgconv.NumericToInt64Safe(value) }
