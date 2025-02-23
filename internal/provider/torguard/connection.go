package torguard

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Torguard) GetConnection(selection settings.ServerSelection) (
	connection models.Connection, err error) {
	defaults := utils.NewConnectionDefaults(1912, 1912, 0) //nolint:gomnd
	return utils.GetConnection(p.servers, selection, defaults, p.randSource)
}
