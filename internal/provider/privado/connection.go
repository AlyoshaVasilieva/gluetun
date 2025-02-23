package privado

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Privado) GetConnection(selection settings.ServerSelection) (
	connection models.Connection, err error) {
	defaults := utils.NewConnectionDefaults(0, 1194, 0) //nolint:gomnd
	return utils.GetConnection(p.servers, selection, defaults, p.randSource)
}
