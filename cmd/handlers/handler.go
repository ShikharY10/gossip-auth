package handlers

import (
	"github.com/ShikharY10/gbAUTH/cmd/admin"
	config "github.com/ShikharY10/gbAUTH/cmd/configs"
)

type Handler struct {
	Logger     *admin.Logger
	Cloudinary *config.Cloudinary
	Cache      *Cache
	DataBase   *DataBase
}
