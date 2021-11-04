package cmd

import (
	"github.com/minamijoyo/myaws/myaws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(newELBV2Cmd())
}

func newELBV2Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "elbv2",
		Short: "Manage ELBV2 resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newELBV2LsCmd(),
		newELBV2PsCmd(),
	)

	return cmd
}

func newELBV2LsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List ELBV2 instances",
		RunE:  runELBV2LsCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("all", "a", false, "List all instances (by default, list running instances only)")
	flags.BoolP("quiet", "q", false, "Only display InstanceIDs")
	flags.StringP("filter-tag", "t", "",
		"Filter instances by tag, such as \"Name:app-production\". The value of tag is assumed to be a partial match",
	)
	flags.StringP("fields", "F", "InstanceId InstanceType PublicIpAddress PrivateIpAddress AvailabilityZone StateName LaunchTime Tag:Name Tag:Service 'Tag:In Charge'", "Output fields list separated by space")
	flags.StringP("domain", "D", "", "Please enter the domain you wish to search")
	viper.BindPFlag("elbv2.ls.all", flags.Lookup("all"))
	viper.BindPFlag("elbv2.ls.quiet", flags.Lookup("quiet"))
	viper.BindPFlag("elbv2.ls.filter-tag", flags.Lookup("filter-tag"))
	viper.BindPFlag("elbv2.ls.fields", flags.Lookup("fields"))
	viper.BindPFlag("elbv2.ls.domain", flags.Lookup("domain"))

	return cmd
}

func runELBV2LsCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	options := myaws.ELBv2Options{
		All:       viper.GetBool("elbv2.ls.all"),
		Quiet:     viper.GetBool("elbv2.ls.quiet"),
		FilterTag: viper.GetString("elbv2.ls.filter-tag"),
		Fields:    viper.GetStringSlice("elbv2.ls.fields"),
		Domain:    viper.GetStringSlice("elbv2.ls.domain"),
	}

	return client.ELBV2Ls(options)
}

func newELBV2PsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ps TARGET_GROUP_NAME",
		Short: "Show ELBV2 target group health",
		RunE:  runELBV2PsCmd,
	}

	return cmd
}

func runELBV2PsCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("TARGET_GROUP_NAME is required")
	}

	options := myaws.ELBV2PsOptions{
		TargetGroupName: args[0],
	}

	return client.ELBV2Ps(options)
}
