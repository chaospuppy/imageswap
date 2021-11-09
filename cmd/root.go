/*
Copyright © 2021 Tim Seagren

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
	// "fmt"
	"context"
	"github.com/spf13/cobra"
	server "imageswap/server"
	"k8s.io/klog/v2"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

var ecrHostname string
var httpPort string
var ecrPattern = regexp.MustCompile(`^(\d{12})\.dkr\.ecr(\-fips)?\.([a-zA-Z0-9][a-zA-Z0-9-_]*)\.(amazonaws\.com(\.cn)?|sc2s\.sgov\.gov|c2s\.ic\.gov)$`)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "imageswap",
	Short: "",
	Long: `A binary used as part of a webhook to replace the existing hostname of a Pod
	image: field with the desired ecr hostname.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			klog.Fatalf("Missing requred ECR hostname positional argument")
		}
		splitURL := ecrPattern.FindStringSubmatch(args[0])
		if len(splitURL) < 4 {
			klog.Fatalf("%s is not a valid ECR repository URL", args[0])
		}

		server := server.NewHTTPServer(httpPort)
		go func() {
			if err := server.ListenAndServe(); err != nil {
				klog.Errorf("Failed to listen and serve: %v", err)
			}
		}()

		klog.Infof("Server listening on port: %v", httpPort)

		// listen shutdown signal
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan

		klog.Infof("Shutdown gracefully...")
		if err := server.Shutdown(context.Background()); err != nil {
			klog.Error(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(&httpPort, "httpPort", "8443", "The port the webhook HTTP Server with listen on.  Defaults to 8443")
}
