// SPDX-License-Identifier: Apache-2.0
// Copyright 2020-2021 Authors of Cilium

package cmd

import (
	"context"
	"os"

	"github.com/cilium/cilium-cli/defaults"
	"github.com/cilium/cilium-cli/hubble"

	"github.com/spf13/cobra"
)

func newCmdHubble() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hubble",
		Short: "Hubble observability",
		Long:  ``,
	}

	cmd.AddCommand(
		newCmdHubbleEnable(),
		newCmdHubbleDisable(),
		newCmdPortForwardCommand(),
		newCmdUI(),
	)

	return cmd
}

func newCmdHubbleEnable() *cobra.Command {
	var params = hubble.Parameters{
		Writer: os.Stdout,
	}

	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable Hubble observability",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			h := hubble.NewK8sHubble(k8sClient, params)
			if err := h.Enable(context.Background()); err != nil {
				fatalf("Unable to enable Hubble:  %s", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&params.Namespace, "namespace", "n", "kube-system", "Namespace Cilium is running in")
	cmd.Flags().BoolVar(&params.Relay, "relay", true, "Deploy Hubble Relay")
	// It can be deprecated since we have a helm option for it
	cmd.Flags().StringVar(&params.RelayImage, "relay-image", "", "Image path to use for Relay")
	// It can be deprecated since we have a helm option for it
	cmd.Flags().StringVar(&params.RelayVersion, "relay-version", "", "Version of Relay to deploy")
	// It can be deprecated since there is not a helm option for it and
	cmd.Flags().StringVar(&params.RelayServiceType, "relay-service-type", "ClusterIP", "Type of Kubernetes service to expose Hubble Relay")
	cmd.Flags().MarkDeprecated("relay-service-type", "value is no longer used for relay-service")
	cmd.Flags().BoolVar(&params.UI, "ui", false, "Enable Hubble UI")

	// It can be deprecated since we have a helm option for it
	cmd.Flags().StringVar(&params.UIImage, "ui-image", "", "Image path to use for UI")
	// It can be deprecated since we have a helm option for it
	cmd.Flags().StringVar(&params.UIBackendImage, "ui-backend-image", "", "Image path to use for UI backend")
	// It can be deprecated since we have a helm option for it
	cmd.Flags().StringVar(&params.UIVersion, "ui-version", "", "Version of UI to deploy")
	cmd.Flags().BoolVar(&params.CreateCA, "create-ca", true, "Automatically create CA if needed")
	cmd.Flags().StringVar(&contextName, "context", "", "Kubernetes configuration context")
	cmd.Flags().BoolVar(&params.Wait, "wait", true, "Wait for status to report success (no errors)")
	cmd.Flags().DurationVar(&params.WaitDuration, "wait-duration", defaults.StatusWaitDuration, "Maximum time to wait for status")

	cmd.Flags().StringVar(&params.K8sVersion, "k8s-version", "", "Kubernetes server version in case auto-detection fails")
	cmd.Flags().StringVar(&params.BaseVersion, "base-version", defaults.Version,
		"Specify the base Cilium version for configuration purpose in case the --version flag doesn't indicate the actual Cilium version")
	cmd.Flags().StringVar(&params.HelmChartDirectory, "chart-directory", "", "Helm chart directory")
	cmd.Flags().StringSliceVar(&params.HelmOpts.ValueFiles, "helm-values", []string{}, "Specify helm values in a YAML file or a URL (can specify multiple)")
	cmd.Flags().StringArrayVar(&params.HelmOpts.Values, "helm-set", []string{}, "Set helm values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	cmd.Flags().StringArrayVar(&params.HelmOpts.StringValues, "helm-set-string", []string{}, "Set helm STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	cmd.Flags().StringArrayVar(&params.HelmOpts.FileValues, "helm-set-file", []string{}, "Set helm values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)")
	cmd.Flags().StringVar(&params.HelmGenValuesFile, "helm-auto-gen-values", "", "Write an auto-generated helm values into this file")

	for flagName := range hubble.FlagsToHelmOpts {
		// TODO(aanm) Do not mark the flags has deprecated for now.
		// msg := fmt.Sprintf("use --helm-set=%s<=value> instead", helmOpt)
		// err := cmd.Flags().MarkDeprecated(flagName, msg)
		// if err != nil {
		// 	panic(err)
		// }
		hubble.FlagValues[flagName] = cmd.Flags().Lookup(flagName).Value
	}

	return cmd
}

func newCmdHubbleDisable() *cobra.Command {
	var params = hubble.Parameters{
		Writer: os.Stdout,
	}

	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable Hubble observability",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			h := hubble.NewK8sHubble(k8sClient, params)
			if err := h.Disable(context.Background()); err != nil {
				fatalf("Unable to disable Hubble:  %s", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&params.Namespace, "namespace", "n", "kube-system", "Namespace Cilium is running in")
	cmd.Flags().StringVar(&contextName, "context", "", "Kubernetes configuration context")

	return cmd
}

func newCmdPortForwardCommand() *cobra.Command {
	var params = hubble.Parameters{
		Writer: os.Stdout,
	}

	cmd := &cobra.Command{
		Use:   "port-forward",
		Short: "Forward the relay port to the local machine",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			h := hubble.NewK8sHubble(k8sClient, params)
			if err := h.PortForwardCommand(context.Background()); err != nil {
				fatalf("Unable to port forward: %s", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&params.Namespace, "namespace", "n", "kube-system", "Namespace Cilium is running in")
	cmd.Flags().StringVar(&params.Context, "context", "", "Kubernetes configuration context")
	cmd.Flags().IntVar(&params.PortForward, "port-forward", 4245, "Local port to forward to")

	return cmd
}

func newCmdUI() *cobra.Command {
	var params = hubble.Parameters{
		Writer: os.Stdout,
	}

	cmd := &cobra.Command{
		Use:   "ui",
		Short: "Open the Hubble UI",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := params.UIPortForwardCommand(context.Background()); err != nil {
				fatalf("Unable to port forward: %s", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&params.Namespace, "namespace", "n", "kube-system", "Namespace Cilium is running in")
	cmd.Flags().StringVar(&params.Context, "context", "", "Kubernetes configuration context")
	cmd.Flags().IntVar(&params.UIPortForward, "port-forward", 12000, "Local port to use for the port forward")

	return cmd
}
