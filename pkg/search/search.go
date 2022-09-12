package search

import (
	"fmt"
	"github.com/BajomoDavid/kubectl-search-plugin/pkg/logger"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	v12 "k8s.io/client-go/kubernetes/typed/apps/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"strings"
	//"k8s.io/client-go/search/pkg/client/auth/gcp"
)

type Result struct {
	ID   int
	Name string
}

func RunPlugin(configFlags *genericclioptions.ConfigFlags, text string) error {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %w", err)
	}

	namespaces, err := clientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list namespaces: %w", err)
	}

	for _, namespace := range namespaces.Items {
		fmt.Println("Namespace: " + namespace.Name)

		ListDeployments(
			namespace.Name,
			clientSet.AppsV1(),
			text,
		)
		println()
		ListServices(
			namespace.Name,
			clientSet.CoreV1(),
			text,
		)
		println()
		ListSecrets(
			namespace.Name,
			clientSet.CoreV1(),
			text,
		)
		println()
		println()
	}
	ListNodes(
		clientSet.CoreV1(),
		text,
	)
	return nil
}

func ListSecrets(
	namespace string,
	clientSet v1.CoreV1Interface,
	needle string,
) {
	secrets, err := clientSet.Secrets(namespace).List(metav1.ListOptions{})
	if err != nil {
		log := logger.NewLogger()
		log.Error(err)
	}
	fmt.Println("Secrets: ")
	tbl := TableHeader("ID", "Name")
	count := 1
	for _, secret := range secrets.Items {
		result := NeedleInHaystack(
			count,
			needle,
			secret.Name,
			secret.Annotations,
			secret.Labels,
		)
		if result.Name != "" {
			tbl.AddRow(result.ID, result.Name)
			count++
		}
	}
	tbl.Print()
}

func ListDeployments(
	namespace string,
	clientSet v12.AppsV1Interface,
	needle string,
) {
	deployments, err := clientSet.Deployments(namespace).List(metav1.ListOptions{})
	if err != nil {
		log := logger.NewLogger()
		log.Error(err)
	}
	fmt.Println("Deployments: ")
	tbl := TableHeader("ID", "Name")
	count := 1
	for _, deployment := range deployments.Items {
		result := NeedleInHaystack(
			count,
			needle,
			deployment.Name,
			deployment.Annotations,
			deployment.Labels,
		)
		if result.Name != "" {
			tbl.AddRow(result.ID, result.Name)
			count++
		}
	}
	tbl.Print()
}

func ListNodes(
	clientSet v1.CoreV1Interface,
	needle string,
) {
	nodes, err := clientSet.Nodes().List(metav1.ListOptions{})
	if err != nil {
		log := logger.NewLogger()
		log.Error(err)
	}
	fmt.Println("Nodes from all namespaces: ")
	tbl := TableHeader("ID", "Name")
	count := 1
	for _, node := range nodes.Items {
		result := NeedleInHaystack(
			count,
			needle,
			node.Name,
			node.Annotations,
			node.Labels,
		)
		if result.Name != "" {
			tbl.AddRow(result.ID, result.Name)
			count++
		}
	}
	tbl.Print()
}

func ListServices(
	namespace string,
	clientSet v1.CoreV1Interface,
	needle string,
) {
	services, err := clientSet.Services(namespace).List(metav1.ListOptions{})
	if err != nil {
		log := logger.NewLogger()
		log.Error(err)
	}
	fmt.Println("Services: ")
	tbl := TableHeader("ID", "Name")
	count := 1
	for _, service := range services.Items {
		result := NeedleInHaystack(
			count,
			needle,
			service.Name,
			service.Annotations,
			service.Labels,
		)
		if result.Name != "" {
			tbl.AddRow(result.ID, result.Name)
			count++
		}
	}
	tbl.Print()
}

func NeedleInHaystack(
	id int,
	needle string,
	name string,
	annotations map[string]string,
	labels map[string]string,
) Result {
	//	Run through properties and find needle
	checkLabels := false
	for _, label := range labels {
		if strings.Contains(label, needle) {
			checkLabels = true
			break
		}
	}

	checkAnnotation := false
	for _, annotation := range annotations {
		if strings.Contains(annotation, needle) {
			checkAnnotation = true
			break
		}
	}

	checkName := strings.Contains(name, needle)
	if checkName || checkLabels || checkAnnotation {
		return Result{
			ID:   id,
			Name: name,
		}
	}
	return Result{}
}

func TableHeader(id string, name string) table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New(id, name)
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	return tbl
}
