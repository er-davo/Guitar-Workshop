//go:build integration

package repository_test

import (
	"context"
	"testing"

	"api-gateway/internal/models"
	"api-gateway/internal/repository"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/stretchr/testify/require"
)

func TestTabRepository_CRUD(t *testing.T) {
	repo := repository.NewTabRepository(db, trmpgx.DefaultCtxGetter)
	ctx := context.Background()

	// CREATE
	tab := &models.Tab{
		Name: "Master of Puppets",
		Path: "/tabs/mop.gp",
	}

	err := repo.Create(ctx, tab)
	require.NoError(t, err)
	require.NotEmpty(t, tab.ID)

	// GET BY ID
	fromDB, err := repo.Get(ctx, tab.ID)
	require.NoError(t, err)
	require.Equal(t, tab.Name, fromDB.Name)
	require.Equal(t, tab.Path, fromDB.Path)

	// FIND BY NAME LIKE
	list, err := repo.FindByNameLike(ctx, "puppet")
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, tab.ID, list[0].ID)

	// DELETE
	err = repo.Delete(ctx, tab.ID)
	require.NoError(t, err)

	// GET AFTER DELETE
	_, err = repo.Get(ctx, tab.ID)
	require.Error(t, err)
}
