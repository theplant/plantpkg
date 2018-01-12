package factory

import (
	"github.com/jinzhu/gorm"
	"github.com/theplant/plantpkg/template"
	"github.com/theplant/plantpkg/template/internal"
)

var _ template.HelloService = (*internal.TemplateImpl)(nil)

func New(db *gorm.DB, requiredParam2 string) (service *internal.TemplateImpl) {
	var err error
	service, err = internal.New(db, requiredParam2)
	if err != nil {
		panic(err)
	}
	return
}
