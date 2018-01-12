package factory

import (
	"github.com/jinzhu/gorm"
	"github.com/theplant/plantpkg/template"
	"github.com/theplant/plantpkg/template/internal"
)

var _ template.TemplateService = (*internal.TemplateImpl)(nil)

func New(db *gorm.DB, cfg *template.Config) (service *internal.TemplateImpl) {
	var err error
	service, err = internal.New(db, cfg)
	if err != nil {
		panic(err)
	}
	return
}
