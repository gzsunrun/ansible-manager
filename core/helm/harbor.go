package helm

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"time"

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

func GetDstCharts(addr string) ([]ChartRepo, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	reqURL := u.Scheme + "://" + u.Host + "/api/chartrepo/" + chartName + "/charts"
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
		reqURL := u.Scheme + "://" + u.Host + "/api/chartrepo/" + chartName + "/charts/" + c.Name
		resp, err := req.Get(reqURL)
		if err != nil {
			return nil, err
		}
		var cs []Chart
		err = resp.ToJSON(&cs)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		res = append(res, ChartRepo{Name: c.Name, Icon: c.Icon, Desc: cs[0].Desc, Cs: cs})
	}
	return res, nil

}

func SyncCharts(srcAddr, dstAddr string, repos []ChartVesion) error {
	srcName := fmt.Sprintf("src-%d", time.Now().UnixNano())
	output10, err := exec.Command("helm", "repo", "add", srcName, srcAddr).CombinedOutput()
	if err != nil {
		log.Error(string(output10))
		return err
	}
	defer func() {
		output3, err := exec.Command("helm", "repo", "remove", srcName).CombinedOutput()
		if err != nil {
			log.Error(string(output3))
		}
	}()
	log.Debug(srcAddr, dstAddr)
	dstName := fmt.Sprintf("dst-%d", time.Now().UnixNano())
	output11, err := exec.Command("helm", "repo", "add", dstName, dstAddr).CombinedOutput()
	if err != nil {
		log.Error(string(output11))
		return err
	}
	defer func() {
		output3, err := exec.Command("helm", "repo", "remove", dstName).CombinedOutput()
		if err != nil {
			log.Error(string(output3))
		}
	}()

	output, err := exec.Command("helm", "repo", "up").CombinedOutput()
	if err != nil {
		log.Error(string(output))
		return err
	}

	for _, repo := range repos {
		log.Debug(repo.Name, repo.Version)
		envVars()
		cmd1 := exec.Command("helm", "fetch", srcName+"/"+repo.Name, "--version", repo.Version, "--destination", "/tmp")
		cmd1.Env = envVars()
		output1, err := cmd1.CombinedOutput()
		if err != nil {
			log.Error(string(output1))
			return err
		}
		cmd2 := exec.Command("helm", "push", "/tmp/"+repo.Name+"-"+repo.Version+".tgz", dstName, "--username=admin", "--password=Harbor12345")
		cmd2.Env = envVars()
		output2, err := cmd2.CombinedOutput()
		if err != nil {
			log.Error(string(output2))
			return err
		}
		os.Remove("/tmp/" + repo.Name + "-" + repo.Version + ".tgz")
	}
	log.Debug("finish sync")
	return nil
}

func envVars() []string {
	env := os.Environ()
	//env = append(env, fmt.Sprintf("HOME=/tmp/ansible-mamager"))
	env = append(env, fmt.Sprintf("PWD=/tmp/ansible-mamager"))
	//env = append(env, fmt.Sprintln("PYTHONUNBUFFERED=1"))
	return env
}
