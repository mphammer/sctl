/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"

	alertsv1 "github.com/yashbhutwala/kb-synopsys-operator/api/v1"
	client "sigs.k8s.io/controller-runtime/pkg/client"

	//git "gopkg.in/src-d/go-git.v4"

	"github.com/spf13/cobra"
)

var label string

func readFromSource(source string) ([]byte, error) {
	return ReadFileData(source)
}

func sendToDestination(data []byte, destination string) error {
	return AddFileToGit(data, destination)
}

func updateCustomResource(data []byte, destination string) error {
	hasher := sha1.New()
	hasher.Write(data)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	cfg, err := GetKubeConfig("", false)
	if err != nil {
		return err
	}
	ctx := context.Background()

	myAlertScheme := runtime.NewScheme()
	alertsv1.AddToScheme(myAlertScheme)
	k8sClient, err := client.New(cfg, client.Options{Scheme: myAlertScheme})
	objKey := client.ObjectKey{
		Name:      "alert-sample",
		Namespace: "default",
	}
	var alert alertsv1.Alert
	fmt.Printf("setting desination: '%s'\n", destination)
	fmt.Printf("setting sha: '%s'\n", sha)

	if err = k8sClient.Get(ctx, objKey, &alert); err != nil {
		fmt.Printf("the Alert CR doesn't exist: %s", err)
		fmt.Printf("creating the Alert CR")
		// create
		alert = alertsv1.Alert{}
		alert.APIVersion = "alerts.synopsys.com/v1"
		alert.Kind = "Alert"
		alert.ObjectMeta.Name = "alert-sample"
		alert.ObjectMeta.Namespace = "default"
		alert.Spec.FinalYamlUrl = destination
		alert.Spec.ShaOfFinalYaml = sha
		if err = k8sClient.Create(ctx, &alert); err != nil {
			return err
		}
	} else {
		// update
		fmt.Printf("updating the Alert CR")
		alert.Spec.FinalYamlUrl = destination
		alert.Spec.ShaOfFinalYaml = sha
		if err = k8sClient.Update(ctx, &alert); err != nil {
			return err
		}
	}
	return nil
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:          "create-cr SOURCE DESINATION",
	Short:        "updates a CR in your cluster",
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			cmd.Help()
			return fmt.Errorf("this command takes two arguments")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		source := args[0]
		destination := args[1]

		fmt.Printf("getting yaml from source '%s'\n", source)
		data, err := ReadFileData(source)
		if err != nil {
			return err
		}

		fmt.Printf("putting yaml in destination '%s'\n", destination)
		//err = sendToDestination(data, destination)
		// if err != nil {
		// 	return err
		// }

		fmt.Printf("updating custom resource definition\n")
		err = updateCustomResource(data, destination)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&label, "customer-label", label, "New label for customers")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
