package article

import (
	"time"

	"github.com/sirupsen/logrus"
)

type UsecaseProperty struct {
	ServiceName string
	UTCLoc      *time.Location
	Logger      *logrus.Logger
	Repository  Repository
}
