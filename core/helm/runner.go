package helm

import (
	"errors"
	"fmt"
	"net"
	"path"
	"strings"
	"time"

	"github.com/hashwing/log"

	"golang.org/x/crypto/ssh"
)

// Runner helm runner
type Runner struct {
	c        *ssh.Client
	Host     string
	Port     string
	User     string
	Password string
	Key      string
}

// NewRunner new a runner
func NewRunner(host, user, password, key, port string) (*Runner, error) {
	if host == "" || (password == "" && key == "") || port == "" {
		return nil, errors.New("parse error")
	}
	if user == "" {
		user = "root"
	}
	return &Runner{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Key:      key,
	}, nil
}

func (r *Runner) client() (*ssh.Client, error) {
	if r.c != nil {
		return r.c, nil
	}
	if r.Password != "" {
		authMethods := []ssh.AuthMethod{}

		keyboardInteractiveChallenge := func(
			user,
			instruction string,
			questions []string,
			echos []bool,
		) (answers []string, err error) {
			if len(questions) == 0 {
				return []string{}, nil
			}
			return []string{r.Password}, nil
		}

		authMethods = append(authMethods, ssh.KeyboardInteractive(keyboardInteractiveChallenge))
		authMethods = append(authMethods, ssh.Password(r.Password))

		c, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", r.Host, r.Port), &ssh.ClientConfig{
			User:    r.User,
			Auth:    authMethods,
			Timeout: 15 * time.Second,
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		})
		if err != nil {
			return nil, err
		}
		r.c = c
		return c, nil
	}

	signer, err := ssh.ParsePrivateKey([]byte(r.Key))
	if err != nil {
		return nil, err
	}

	c, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", r.Host, r.Port), &ssh.ClientConfig{
		User:    r.User,
		Auth:    []ssh.AuthMethod{ssh.PublicKeys(signer)},
		Timeout: 15 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})

	if err != nil {
		return nil, err
	}
	r.c = c
	return c, nil
}

type HelmRelease struct {
	Name       string `json:"name"`
	Revison    string `json:"revison"`
	Upadted    string `json:"upadted"`
	Status     string `json:"status"`
	Chart      string `json:"chart"`
	AppVersion string `json:"app_version"`
	Namespace  string `json:"namespace"`
}

func (r *Runner) Close() {
	if r.c != nil {
		r.c.Close()
	}

}

func (r *Runner) HelmList() ([]HelmRelease, error) {
	c, err := r.client()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	session, err := c.NewSession()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	buf, err := session.CombinedOutput("helm list -a")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	listsStr := strings.Split(string(buf), "\n")
	lists := make([]HelmRelease, 0)
	for _, l := range listsStr[1 : len(listsStr)-1] {
		infoArr := strings.Split(l, "\t")
		release := HelmRelease{
			Name:       strings.TrimSpace(infoArr[0]),
			Revison:    strings.TrimSpace(infoArr[1]),
			Upadted:    strings.TrimSpace(infoArr[2]),
			Status:     strings.TrimSpace(infoArr[3]),
			Chart:      strings.TrimSpace(infoArr[4]),
			AppVersion: strings.TrimSpace(infoArr[5]),
			Namespace:  strings.TrimSpace(infoArr[6]),
		}
		lists = append(lists, release)
	}

	return lists, err
}

func (r *Runner) copyFile(content, dest string) error {
	c, err := r.client()
	if err != nil {
		log.Error(err)
		return err
	}
	session, err := c.NewSession()
	if err != nil {
		log.Error(err)
		return err
	}
	defer session.Close()
	dir, name := path.Split(dest)
	go func() {
		w, err := session.StdinPipe()
		if err != nil {
			log.Error(err)
		}
		defer w.Close()

		fmt.Fprintln(w, "C0644", len(content), name)
		fmt.Fprint(w, content)
		fmt.Fprint(w, "\x00")
	}()
	if err := session.Run("mkdir -p " + dir + " &&/usr/bin/scp -qrt " + dir); err != nil {
		if err != nil {
			if err.Error() != "Process exited with: 1. Reason was:  ()" {
				log.Error(err)
				return err
			}
		}
	}
	return nil

}

func (r *Runner) addHelmRepo(repo string, name ...string) error {
	repoName := "ansible-manager"
	if len(name) > 0 {
		repoName = name[0]
	}
	c, err := r.client()
	if err != nil {
		log.Error(err)
		return err
	}
	session, err := c.NewSession()
	if err != nil {
		log.Error(err)
		return err
	}
	defer session.Close()
	output, err := session.CombinedOutput("helm repo add " + repoName + " " + repo)
	if err != nil {
		log.Error(string(output))
		return err
	}
	return nil
}

type ChartVesion struct {
	Name    string
	Version string
}

func (r *Runner) SyncHelmRepo(addr string, repos []ChartVesion) error {
	c, err := r.client()
	if err != nil {
		log.Error(err)
		return err
	}
	session, err := c.NewSession()
	if err != nil {
		log.Error(err)
		return err
	}
	defer session.Close()
	repoName := time.Now().UnixNano()
	output, err := session.CombinedOutput(fmt.Sprintf("helm repo add sync-%d %s", repoName, addr))
	if err != nil {
		log.Error(string(output))
		return err
	}

	for _, repo := range repos {
		output1, err := session.CombinedOutput(fmt.Sprintf("helm fetch sync-%d/%s --version %s", repoName, repo.Name, repo.Version))
		if err != nil {
			log.Error(string(output1))
			return err
		}
		output2, err := session.CombinedOutput(fmt.Sprintf("helm push sync-%d %s-%s.tgz", repoName, repo.Name, repo.Version))
		if err != nil {
			log.Error(string(output2))
			return err
		}
	}
	output3, err := session.CombinedOutput(fmt.Sprintf("helm repo remove sync-%d %s", repoName, addr))
	if err != nil {
		log.Error(string(output3))
		return err
	}
	return nil
}

func (r *Runner) updateRepo() error {
	session, err := r.newSession()
	if err != nil {
		log.Error(err)
		return err
	}
	defer session.Close()
	output, err := session.CombinedOutput("helm repo update")
	if err != nil {
		log.Error(string(output))
		return err
	}
	return nil
}

func (r *Runner) newSession() (*ssh.Session, error) {
	c, err := r.client()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	session, err := c.NewSession()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return session, nil

}

func (r *Runner) GetValues(name string) (string, error) {
	session, err := r.newSession()
	if err != nil {
		log.Error(err)
		return "", err
	}
	defer session.Close()
	output, err := session.CombinedOutput("helm get values " + name)
	if err != nil {
		log.Error(string(output))
		return "", err
	}
	log.Debug(string(output))
	return string(output), nil
}

// Install install / update a helm chart
func (r *Runner) Install(name, repo, chart, version, namespace, values string, update bool) error {
	valuesPath := "/opt/ansible-manager/helm-values/" + name + "-values.yaml"
	err := r.copyFile(values, valuesPath)
	if err != nil {
		log.Error(err)
		return err
	}

	err = r.addHelmRepo(repo)
	if err != nil {
		log.Error(err)
		return err
	}

	r.updateRepo()

	session, err := r.newSession()
	if err != nil {
		log.Error(err)
		return err
	}
	defer session.Close()
	cmd := ""
	if update {
		cmd = fmt.Sprintf("helm upgrade %s ansible-manager/%s  -f  %s", name, chart, valuesPath)
	} else {
		cmd = fmt.Sprintf("helm install --name %s --namespace %s ansible-manager/%s  -f  %s", name, namespace, chart, valuesPath)
	}
	if version != "" {
		cmd += " --version " + version
	}
	log.Debug(cmd)
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Error(string(output))
		return errors.New(string(output))
	}
	return err
}

type ReleaseStatus struct {
	PodsStatus []PodStatus `json:"pods_status"`
}
type PodStatus struct {
	Name    string `json:"name"`
	Ready   string `json:"ready"`
	Status  string `json:"status"`
	Restart string `json:"restart"`
	Age     string `json:"age"`
}

func (r *Runner) HelmStatus(name string) (*ReleaseStatus, error) {
	session, err := r.newSession()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer session.Close()
	output, err := session.CombinedOutput("helm status " + name)
	if err != nil {
		log.Error(string(output))
		return nil, err
	}
	log.Debug(string(output))
	lines := strings.Split(string(output), "\n")
	podFlag := false
	podsStatus := make([]PodStatus, 0)
	for i := 0; i < len(lines); i++ {
		log.Debug(lines[i])
		if podFlag && lines[i] == "" {
			break
		}
		if podFlag {
			podLine := strings.Split(lines[i], " ")
			j := 0
			for _, podItem := range podLine {
				if podItem != "" {
					podLine[j] = podItem
					j++
				}

			}
			log.Debug(len(podLine))
			if len(podLine) >= 5 {
				podStatus := PodStatus{
					Name:    strings.TrimSpace(podLine[0]),
					Ready:   strings.TrimSpace(podLine[1]),
					Status:  strings.TrimSpace(podLine[2]),
					Restart: strings.TrimSpace(podLine[3]),
					Age:     strings.TrimSpace(podLine[4]),
				}
				podsStatus = append(podsStatus, podStatus)
			}
		}

		if strings.HasPrefix(lines[i], "==> v1/Pod") {
			podFlag = true
			i = i + 2
		}
	}
	log.Debug(podsStatus)
	return &ReleaseStatus{PodsStatus: podsStatus}, err
}

// HelmDelete delete helm release
// @name release name @purge purge or not
func (r *Runner) HelmDelete(name string, purge bool) error {
	session, err := r.newSession()
	if err != nil {
		log.Error(err)
		return err
	}
	defer session.Close()
	cmd := "helm delete " + name
	if purge {
		cmd += " --purge"
	}
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Error(string(output))
		return errors.New(string(output))
	}
	return nil
}
