package appinfousecase

import (
	"github.com/Tanapoowapat/GunplaShop/modules/appinfo"
	appinforepositories "github.com/Tanapoowapat/GunplaShop/modules/appinfo/appinfoRepositories"
)

type IAppinfoUsecase interface {
	FindCategory(req *appinfo.CategoryFiter) ([]*appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) error
	DeleteCategory(category_id int) error
}

type appinfoUsecase struct {
	appinfo_repo appinforepositories.IAppinfoRepositories
}

func AppinfoRepositories(appinforepositories appinforepositories.IAppinfoRepositories) IAppinfoUsecase {
	return &appinfoUsecase{
		appinfo_repo: appinforepositories,
	}
}

func (u *appinfoUsecase) FindCategory(req *appinfo.CategoryFiter) ([]*appinfo.Category, error) {
	category, err := u.appinfo_repo.FindCategory(req)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (u *appinfoUsecase) InsertCategory(req []*appinfo.Category) error {
	if err := u.appinfo_repo.InsertCategory(req); err != nil {
		return nil
	}
	return nil
}

func (u *appinfoUsecase) DeleteCategory(category_id int) error {
	if err := u.appinfo_repo.DeleteCategory(category_id); err != nil {
		return err
	}
	return nil
}
