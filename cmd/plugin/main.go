package main

import (
	"github.com/BajomoDavid/kubectl-search-plugin/cmd/plugin/cli"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // required for GKE
)

func main() {
	cli.InitAndExecute()
}
