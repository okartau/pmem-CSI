/*
Copyright 2020 Intel Coporation.

SPDX-License-Identifier: Apache-2.0
*/

package k8sutil

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/intel/pmem-csi/pkg/pmem-csi-operator/version"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

// NewInClusterClient connects code that runs inside a Kubernetes pod to the
// API server.
func NewInClusterClient() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("build in-cluster Kubernetes client configuration: %v", err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("create Kubernetes client: %v", err)
	}
	return client, nil
}

// GetKubernetesVersion returns kubernetes server version
func GetKubernetesVersion(cfg *rest.Config) (*version.Version, error) {
	client, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, err
	}
	ver, err := client.ServerVersion()
	if err != nil {
		return nil, err
	}

	klog.Infof("Kubernetes version read from server: %v.%v", ver.Major, ver.Minor)

	// Suppress all non digits, version might contain special charcters like, <number>+
	reg, _ := regexp.Compile("[^0-9]+")
	major, err := strconv.Atoi(reg.ReplaceAllString(ver.Major, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to parse Kubernetes major version %q: %v", ver.Major, err)
	}
	minor, err := strconv.Atoi(reg.ReplaceAllString(ver.Minor, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to parse Kubernetes minor version %q: %v", ver.Minor, err)
	}

	return version.NewVersion(uint(major), uint(minor)), nil
}
