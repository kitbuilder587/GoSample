package db

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/kitbuilder587/cryptotrack/internal/migrations"
)

func waitForDB(db *sqlx.DB, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if err := db.Ping(); err == nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("DB did not respond in %s", timeout)
}

func setupTestDB(t *testing.T) *sqlx.DB {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "cryptotrack_test",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("database system is ready to accept connections"),
		),
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		postgresC.Terminate(ctx)
	})

	host, err := postgresC.Host(ctx)
	require.NoError(t, err)
	port, err := postgresC.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := fmt.Sprintf("postgres://test:test@%s:%s/cryptotrack_test?sslmode=disable", host, port.Port())
	db, err := sqlx.Open("pgx", dsn)
	require.NoError(t, err)

	require.NoError(t, waitForDB(db, 10*time.Second))

	// Миграции
	require.NoError(t, migrations.Run(db))

	return db
}

func TestRepo_InsertAndGetLatest(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepo(db)
	ctx := context.Background()
	now := time.Now().UTC()

	price := decimal.NewFromFloat(70000.123456)
	pl := PriceLog{Coin: "bitcoin", PriceUSD: price, Timestamp: now}
	require.NoError(t, repo.InsertPrice(ctx, pl))

	got, err := repo.GetLatest(ctx, "bitcoin")
	require.NoError(t, err)
	require.Equal(t, pl.Coin, got.Coin)
	require.True(t, got.PriceUSD.Equal(price))
}

func TestRepo_GetHistory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepo(db)
	ctx := context.Background()
	now := time.Now().UTC()

	for i := 0; i < 5; i++ {
		price := decimal.NewFromFloat(70000.123456 + float64(i))
		pl := PriceLog{Coin: "bitcoin", PriceUSD: price, Timestamp: now.Add(time.Duration(i) * time.Minute)}
		require.NoError(t, repo.InsertPrice(ctx, pl))
	}

	history, err := repo.GetHistory(ctx, "bitcoin", 3)
	require.NoError(t, err)
	require.Len(t, history, 3)
	// Убедимся, что цены убывают по времени (по последнему времени — первый элемент самый новый)
	require.True(t, history[0].Timestamp.After(history[1].Timestamp) || history[0].Timestamp.Equal(history[1].Timestamp))
	require.True(t, history[1].Timestamp.After(history[2].Timestamp) || history[1].Timestamp.Equal(history[2].Timestamp))
}
