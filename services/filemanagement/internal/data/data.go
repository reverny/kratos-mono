package data

import (
	"github.com/reverny/kratos-mono/services/filemanagement/internal/biz"
	"github.com/reverny/kratos-mono/services/filemanagement/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewData, NewFileStorage)

type Data struct {
	log *log.Helper
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		log: log.NewHelper(logger),
	}, cleanup, nil
}

// NewFileStorage creates a file storage based on configuration
func NewFileStorage(c *conf.Storage, logger log.Logger) biz.FileStorage {
	log := log.NewHelper(logger)
	
	switch c.Type {
	case "local":
		log.Infof("Using local file storage: %s", c.UploadPath)
		return NewLocalFileStorage(c.BaseUrl, c.UploadPath)
	case "minio":
		// TODO: Implement MinIO storage
		log.Warn("MinIO storage not implemented yet, falling back to local storage")
		return NewLocalFileStorage(c.BaseUrl, c.UploadPath)
	case "s3":
		// TODO: Implement S3 storage
		log.Warn("S3 storage not implemented yet, falling back to local storage")
		return NewLocalFileStorage(c.BaseUrl, c.UploadPath)
	default:
		log.Warnf("Unknown storage type: %s, using local storage", c.Type)
		return NewLocalFileStorage(c.BaseUrl, c.UploadPath)
	}
}
