package productsrepositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Tanapoowapat/GunplaShop/config"
	"github.com/Tanapoowapat/GunplaShop/modules/entities"
	filesusecase "github.com/Tanapoowapat/GunplaShop/modules/file/filesUsecase"
	"github.com/Tanapoowapat/GunplaShop/modules/products"
	productspatterns "github.com/Tanapoowapat/GunplaShop/modules/products/productsPatterns"
	"github.com/jmoiron/sqlx"
)

type IProductRepositorise interface {
	FindOneProducts(productId string) (*products.Products, error)
	FindProduct(req *products.ProductFilter) ([]*products.Products, int)
	InsertProduct(req *products.Products) (*products.Products, error)
	DeleteProduct(productId string) error
	UpdateProduct(req *products.Products) (*products.Products, error)
}

type productsRepositories struct {
	db           *sqlx.DB
	cfg          config.IConfig
	filesUsecase filesusecase.IFileUsecase
}

func NewProductRepositories(db *sqlx.DB, cfg config.IConfig, filesUsecase filesusecase.IFileUsecase) IProductRepositorise {
	return &productsRepositories{
		db:           db,
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}

func (repo *productsRepositories) FindOneProducts(productId string) (*products.Products, error) {
	query := `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
			"p"."id",
			"p"."title",
			"p"."description",
			"p"."price",
			(
				SELECT
					to_jsonb("ct")
				FROM (
					SELECT
						"c"."id",
						"c"."title"
					FROM "categories" "c"
						LEFT JOIN "products_categories" "pc" ON "pc"."category_id" = "c"."id"
					WHERE "pc"."product_id" = "p"."id"
				) AS "ct"
			) AS "category",
			"p"."created_at",
			"p"."updated_at",
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "images" "i"
					WHERE "i"."product_id" = "p"."id"
				) AS "it"
			) AS "media"
		FROM "products" "p"
		WHERE "p"."id" = $1
		LIMIT 1
	) AS "t";`

	productBytes := make([]byte, 0)
	product := &products.Products{
		Images: make([]*entities.Images, 0),
	}

	if err := repo.db.Get(&productBytes, query, productId); err != nil {
		return nil, fmt.Errorf("get products fails: %v", err)
	}
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, fmt.Errorf("fail to unmarshal json: %v", err)
	}

	return product, nil
}

func (repo *productsRepositories) FindProduct(req *products.ProductFilter) ([]*products.Products, int) {
	builder := productspatterns.NewFindProductBuilder(repo.db, req)
	engineer := productspatterns.NewFindProductEngineer(builder)

	result := engineer.FindProduct().Result()
	count := engineer.CountProduct().Count()
	return result, count
}

func (repo *productsRepositories) DeleteProduct(productId string) error {
	query := `DELETE FROM "producuts" WHERE "id" = $1;`
	if _, err := repo.db.ExecContext(context.Background(), query, productId); err != nil {
		return fmt.Errorf("delete products fail: %v", err)
	}
	return nil
}

func (repo *productsRepositories) InsertProduct(req *products.Products) (*products.Products, error) {
	builder := productspatterns.NewInsertProductBuilder(repo.db, req)
	productId, err := productspatterns.NewInsertProductEngineer(builder).InsertProduct()

	if err != nil {
		return nil, err
	}

	product, err := repo.FindOneProducts(productId)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (repo *productsRepositories) UpdateProduct(req *products.Products) (*products.Products, error) {
	builder := productspatterns.NewUpdateProductBuilder(repo.db, req, repo.filesUsecase)
	engineer := productspatterns.NewUpdateProductEngineer(builder)

	if err := engineer.UpdateProduct(); err != nil {
		return nil, err
	}

	product, err := repo.FindOneProducts(req.Id)
	if err != nil {
		return nil, err
	}

	return product, nil
}
