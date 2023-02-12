package helpers

import (
	"fmt"

	"github.com/rsan92/teste-vibbra/internal/infra/configuration"
)

func GetConnString(conf configuration.AppConfigurations) string {
	return fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		conf.Database.User,
		conf.Database.Password,
		conf.Name)

}
