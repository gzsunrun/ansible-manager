package helm

import (
	"errors"

	"github.com/hashwing/log"
	"github.com/imroc/req"
)

var (
	harborURL = ""
	chartName = ""
)

func InitHarbor(hurl string, chart string) {
	harborURL = hurl
	chartName = chart
}

type ChartRepo struct {
	Name string  `json:"name"`
	Desc string  `json:"description"`
	Icon string  `json:"icon"`
	Cs   []Chart `json:"charts"`
}

type ChartList struct {
	Name string
	Icon string
}

type Chart struct {
	Name       string `json:"-"`
	Version    string `json:"version"`
	Desc       string `json:"description"`
	AppVersion string `json:"appVersion"`
	Icon       string `json:"icon"`
}

func GetCharts() ([]ChartRepo, error) {
	reqURL := harborURL + "/api/chartrepo/" + chartName + "/charts"
	log.Debug(reqURL)
	resp, err := req.Get(reqURL)
	if err != nil {
		return nil, err
	}
	if resp.Response().StatusCode != 200 {
		return nil, errors.New(resp.String())
	}
	var cs []ChartList
	err = resp.ToJSON(&cs)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	res := make([]ChartRepo, 0)
	for _, c := range cs {
		reqURL := harborURL + "/api/chartrepo/" + chartName + "/charts/" + c.Name
		resp, err := req.Get(reqURL)
		if err != nil {
			return nil, err
		}
		var cs []Chart
		err = resp.ToJSON(&cs)
		if err != nil {
			return nil, err
		}
		res = append(res, ChartRepo{Name: c.Name, Icon: c.Icon, Desc: cs[0].Desc, Cs: cs})
	}
	return res, nil

}

func GetValues(name, version string) ([]byte, error) {
	reqURL := harborURL + "/api/chartrepo/" + chartName + "/charts/" + name + "/" + version
	resp, err := req.Get(reqURL)
	if err != nil {
		return nil, err
	}
	return resp.ToBytes()
}
