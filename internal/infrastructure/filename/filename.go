package filename

import (
	"fmt"
	"time"
)

func GetCertificateName() string {
	return fmt.Sprintf("certificate_%d.json", time.Now().UnixNano())
}
