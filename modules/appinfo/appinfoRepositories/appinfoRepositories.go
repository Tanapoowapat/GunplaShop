package appinforepositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/Tanapoowapat/GunplaShop/modules/appinfo"
	"github.com/jmoiron/sqlx"
)

type IAppinfoRepositories interface {
	FindCategory(req *appinfo.CategoryFiter) ([]*appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) error
	DeleteCategory(category_id int) error
}

type appinfoRepositories struct {
	db *sqlx.DB
}

func AppinfoRepositories(db *sqlx.DB) IAppinfoRepositories {
	return &appinfoRepositories{
		db: db,
	}
}

func (r *appinfoRepositories) FindCategory(req *appinfo.CategoryFiter) ([]*appinfo.Category, error) {
	query := `SELECT "id", "title" FROM "categories"`
	filterVals := make([]interface{}, 0)

	// Add WHERE clause if title filter is provided
	if req.Title != "" {
		query += ` WHERE LOWER("title") LIKE LOWER($1)`
		filterVals = append(filterVals, "%"+strings.ToLower(req.Title)+"%")
	}

	query += ";"

	category := make([]*appinfo.Category, 0)
	if err := r.db.Select(&category, query, filterVals...); err != nil {
		return nil, fmt.Errorf("select category failed: %v", err)
	}

	return category, nil
}

func (r *appinfoRepositories) InsertCategory(req []*appinfo.Category) error {
	ctx := context.Background()
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, category := range req {
		query := `INSERT INTO "categories" ("title") VALUES ($1) RETURNING "id"`
		var categoryId int
		if err := tx.QueryRowxContext(ctx, query, category.Title).Scan(&categoryId); err != nil {
			return fmt.Errorf("insert category failed: %v", err)
		}
		category.Id = categoryId
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r appinfoRepositories) DeleteCategory(category_id int) error {
	query := `DELETE FROM "categories" where "id" = $1`
	ctx := context.Background()

	if _, err := r.db.ExecContext(ctx, query, category_id); err != nil {
		return fmt.Errorf("delete category failed: %v", err)
	}
	return nil

}
