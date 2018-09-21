/*
Copyright (c) 2016-2017 Bitnami

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/epimorphics/s3-trigger/pkg/controller"
	s3utils "github.com/epimorphics/s3-trigger/pkg/utils"
	"github.com/epimorphics/s3-trigger/pkg/version"
	kubelessutils "github.com/kubeless/kubeless/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	globalUsage = `S3 trigger controller adds support for S3 changes as event source to Kubeless functions`
)

var rootCmd = &cobra.Command{
	Use:   "s3-controller",
	Short: "s3-controller",
	Long:  globalUsage,
	Run: func(cmd *cobra.Command, args []string) {

		kubelessClient, err := kubelessutils.GetFunctionClientInCluster()
		if err != nil {
			logrus.Fatalf("Cannot get kubeless CR API client: %v", err)
		}

		s3Client, err := s3utils.GetTriggerClientInCluster()
		if err != nil {
			logrus.Fatalf("Cannot get s3 trigger CR API client: %v", err)
		}

		s3TriggerCfg := controller.S3TriggerConfig{
			KubeCli:        kubelessutils.GetClient(),
			TriggerClient:  s3Client,
			KubelessClient: kubelessClient,
		}

		s3TriggerController := controller.NewS3TriggerController(s3TriggerCfg)

		stopCh := make(chan struct{})
		defer close(stopCh)

		go s3TriggerController.Run(stopCh)

		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGTERM)
		signal.Notify(sigterm, syscall.SIGINT)
		<-sigterm
	},
}

func main() {
	logrus.Infof("Running S3 controller version: %v", version.Version)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
