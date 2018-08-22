package function

import (
	"strings"
	"github.com/satori/go.uuid"
)

func NewUuidString() string {
	id := uuid.Must(uuid.NewV4())
	return strings.Replace(id.String(), "-", "", -1)
}

func NewUuidV4String() string {
	id := uuid.Must(uuid.NewV4())
	return id.String() 
}
 