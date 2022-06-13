/*Package cmd root.go
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
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/chaospuppy/imageswap/pkg/imageswap"
	"github.com/chaospuppy/imageswap/pkg/server"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var is imageswap.Imageswap

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "imageswap",
	Short: "",
	Long: `A binary used as part of a webhook to replace the existing hostname of a Pod
	image: field with the desired registry hostname.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			klog.Fatalf("Missing requred hostname positional argument")
		}

		// TODO run additional validation on provided hostname to ensure it matches a registry regex
		svr := is.NewHTTPServer()

		idleConnsClosed := make(chan struct{})
		go func() {
			// Listen for shutdown signal
			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
			<-signalChan

			klog.Infof("Shutdown gracefully...")
			if err := svr.Shutdown(context.Background()); err != nil {
				klog.Error(err)
			}
			// Wait for graceful termination
			close(idleConnsClosed)
		}()

		server.RunHTTPServer(svr)
		klog.Infof("Server listening on port: %v", is.HttpPort)
		<-idleConnsClosed
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(&is.HttpPort, "httpPort", "8443", "The port the webhook HTTP Server with listen on.  Defaults to 8443")
	rootCmd.PersistentFlags().StringVar(&is.TlsKey, "tlsKeyFile", "/etc/webhook/certs/key.pem", "Path to TLS key")
	rootCmd.PersistentFlags().StringVar(&is.TlsCert, "tlsCertFile", "/etc/webhook/certs/cert.pem", "Path to TLS certificate")
	rootCmd.PersistentFlags().StringVar(&is.Annotation, "imageswapAnnotation", "imageswap.chaospuppy.github.com", "Annotation key for imageswap to monitor")
	rootCmd.PersistentFlags().StringVar(&is.RegistryHostname, "registryHostname", "", "registry fqdn to replace original registry with")
}
