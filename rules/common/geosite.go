package common

import (
	"fmt"
	_ "unsafe"

	"github.com/Dreamacro/clash/component/geodata"
	_ "github.com/Dreamacro/clash/component/geodata/memconservative"
	"github.com/Dreamacro/clash/component/geodata/router"
	_ "github.com/Dreamacro/clash/component/geodata/standard"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
)

//go:linkname initGeoSite github.com/Dreamacro/clash/config.initGeoSite
func initGeoSite() error

type GEOSITE struct {
	*Base
	country    string
	adapter    string
	matcher    *router.DomainMatcher
	recodeSize int
}

func (gs *GEOSITE) RuleType() C.RuleType {
	return C.GEOSITE
}

func (gs *GEOSITE) Match(metadata *C.Metadata) bool {
	if metadata.AddrType != C.AtypDomainName {
		return false
	}

	domain := metadata.Host
	return gs.matcher.ApplyDomain(domain)
}

func (gs *GEOSITE) Adapter() string {
	return gs.adapter
}

func (gs *GEOSITE) Payload() string {
	return gs.country
}

func (gs *GEOSITE) GetDomainMatcher() *router.DomainMatcher {
	return gs.matcher
}

func (gs *GEOSITE) GetRecodeSize() int {
	return gs.recodeSize
}

func NewGEOSITE(country string, adapter string) (*GEOSITE, error) {
	if !initFlag {
		if err := initGeoSite(); err != nil {
			log.Errorln("can't initial GeoSite: %s", err)
			return nil, err
		}
		initFlag = true
	}

	matcher, size, err := geodata.LoadGeoSiteMatcher(country)
	if err != nil {
		return nil, fmt.Errorf("load GeoSite data error, %s", err.Error())
	}

	log.Infoln("Start initial GeoSite rule %s => %s, records: %d", country, adapter, size)

	geoSite := &GEOSITE{
		Base:       &Base{},
		country:    country,
		adapter:    adapter,
		matcher:    matcher,
		recodeSize: size,
	}

	return geoSite, nil
}

var _ C.Rule = (*GEOSITE)(nil)
