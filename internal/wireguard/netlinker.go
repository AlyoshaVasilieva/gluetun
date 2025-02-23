package wireguard

import "github.com/qdm12/gluetun/internal/netlink"

//go:generate mockgen -destination=netlinker_mock_test.go -package wireguard . NetLinker

type NetLinker interface {
	AddrAdd(link netlink.Link, addr *netlink.Addr) error
	RouteList(link netlink.Link, family int) (routes []netlink.Route, err error)
	RouteAdd(route *netlink.Route) error
	RuleAdd(rule *netlink.Rule) error
	RuleDel(rule *netlink.Rule) error
	LinkAdd(link netlink.Link) (err error)
	LinkList() (links []netlink.Link, err error)
	LinkByName(name string) (link netlink.Link, err error)
	LinkSetUp(link netlink.Link) error
	LinkSetDown(link netlink.Link) error
	LinkDel(link netlink.Link) error
	IsWireguardSupported() (ok bool, err error)
}
