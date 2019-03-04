package k8s

import (
	"crypto/tls"
	"net/http"
	"os/exec"
)

//KubernetesCmd is a thing
type KubernetesCmd struct {
	netClient http.Client
}

// NewKubernetesCmd does stuff
func NewKubernetesCmd(insecureSsl bool) KubernetesCmd {
	config := &tls.Config{
		InsecureSkipVerify: insecureSsl,
	}
	tr := &http.Transport{TLSClientConfig: config}
	return KubernetesCmd{
		netClient: http.Client{
			Transport: tr,
		},
	}
}

// ExecuteCommand does stuff
func (kc KubernetesCmd) ExecuteCommand(command ...string) (string, error) {

	baseCommand := command

	cmd := exec.Command("kubectl", baseCommand...) //Execute without kubeconfig option
	//log.Printf("%v", cmd.Args)
	out, err := cmd.CombinedOutput()
	//fmt.Println(string(out))
	return string(out), err
}
