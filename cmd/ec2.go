package cmd

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/minamijoyo/myaws/myaws"
)

func init() {
	RootCmd.AddCommand(newEC2Cmd())
}

func newEC2Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ec2",
		Short: "Manage EC2 resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newEC2LsCmd(),
		newEC2VLsCmd(),
		newEC2ILsCmd(),
		newEC2ALsCmd(),
		newEC2SLsCmd(),
		newEC2StartCmd(),
		newEC2StopCmd(),
		newEC2SSHCmd(),
	)

	return cmd
}

func newEC2SLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sls",
		Short: "List Snapshot",
		RunE:  runEC2SLsCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("all", "a", false, "List all instances (by default, list running instances only)")
	flags.BoolP("quiet", "q", false, "Only display InstanceIDs")
	flags.StringP("filter-tag", "t", "",
		"Filter instances by tag, such as \"Name:app-production\". The value of tag is assumed to be a partial match",
	)
	flags.StringP("fields", "F", "InstanceId InstanceType PublicIpAddress PrivateIpAddress AvailabilityZone StateName LaunchTime Tag:Name Tag:Service 'Tag:In Charge'", "Output fields list separated by space")
	flags.StringP("domain", "D", "", "Please enter the domain you wish to search")
	viper.BindPFlag("ec2.sls.all", flags.Lookup("all"))
	viper.BindPFlag("ec2.sls.quiet", flags.Lookup("quiet"))
	viper.BindPFlag("ec2.sls.filter-tag", flags.Lookup("filter-tag"))
	viper.BindPFlag("ec2.sls.fields", flags.Lookup("fields"))
	viper.BindPFlag("ec2.sls.domain", flags.Lookup("domain"))

	return cmd
}

func runEC2SLsCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	options := myaws.EC2SLsOptions{
		All:       viper.GetBool("ec2.sls.all"),
		Quiet:     viper.GetBool("ec2.sls.quiet"),
		FilterTag: viper.GetString("ec2.sls.filter-tag"),
		Fields:    viper.GetStringSlice("ec2.sls.fields"),
		Domain:    viper.GetStringSlice("ec2.sls.domain"),
	}

	return client.EC2SLs(options)
}

func newEC2ALsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "als",
		Short: "List Images",
		RunE:  runEC2ALsCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("all", "a", false, "List all instances (by default, list running instances only)")
	flags.BoolP("quiet", "q", false, "Only display InstanceIDs")
	flags.StringP("filter-tag", "t", "",
		"Filter instances by tag, such as \"Name:app-production\". The value of tag is assumed to be a partial match",
	)
	flags.StringP("fields", "F", "InstanceId InstanceType PublicIpAddress PrivateIpAddress AvailabilityZone StateName LaunchTime Tag:Name Tag:Service 'Tag:In Charge'", "Output fields list separated by space")
	flags.StringP("domain", "D", "", "Please enter the domain you wish to search")
	viper.BindPFlag("ec2.als.all", flags.Lookup("all"))
	viper.BindPFlag("ec2.als.quiet", flags.Lookup("quiet"))
	viper.BindPFlag("ec2.als.filter-tag", flags.Lookup("filter-tag"))
	viper.BindPFlag("ec2.als.fields", flags.Lookup("fields"))
	viper.BindPFlag("ec2.als.domain", flags.Lookup("domain"))

	return cmd
}

func runEC2ALsCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	options := myaws.EC2ALsOptions{
		All:       viper.GetBool("ec2.als.all"),
		Quiet:     viper.GetBool("ec2.als.quiet"),
		FilterTag: viper.GetString("ec2.als.filter-tag"),
		Fields:    viper.GetStringSlice("ec2.als.fields"),
		Domain:    viper.GetStringSlice("ec2.als.domain"),
	}

	return client.EC2ALs(options)
}

func newEC2ILsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ils",
		Short: "List ElasticIP",
		RunE:  runEC2ILsCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("all", "a", false, "List all instances (by default, list running instances only)")
	flags.BoolP("quiet", "q", false, "Only display InstanceIDs")
	flags.StringP("filter-tag", "t", "",
		"Filter instances by tag, such as \"Name:app-production\". The value of tag is assumed to be a partial match",
	)
	flags.StringP("fields", "F", "InstanceId InstanceType PublicIpAddress PrivateIpAddress AvailabilityZone StateName LaunchTime Tag:Name Tag:Service 'Tag:In Charge'", "Output fields list separated by space")
	flags.StringP("domain", "D", "", "Please enter the domain you wish to search")
	viper.BindPFlag("ec2.ils.all", flags.Lookup("all"))
	viper.BindPFlag("ec2.ils.quiet", flags.Lookup("quiet"))
	viper.BindPFlag("ec2.ils.filter-tag", flags.Lookup("filter-tag"))
	viper.BindPFlag("ec2.ils.fields", flags.Lookup("fields"))
	viper.BindPFlag("ec2.ils.domain", flags.Lookup("domain"))

	return cmd
}

func runEC2ILsCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	options := myaws.EC2ILsOptions{
		All:       viper.GetBool("ec2.ils.all"),
		Quiet:     viper.GetBool("ec2.ils.quiet"),
		FilterTag: viper.GetString("ec2.ils.filter-tag"),
		Fields:    viper.GetStringSlice("ec2.ils.fields"),
		Domain:    viper.GetStringSlice("ec2.ils.domain"),
	}

	return client.EC2ILs(options)
}

func newEC2VLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vls",
		Short: "List EC2 Volumes",
		RunE:  runEC2VLsCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("all", "a", false, "List all instances (by default, list running instances only)")
	flags.BoolP("quiet", "q", false, "Only display InstanceIDs")
	flags.StringP("filter-tag", "t", "",
		"Filter instances by tag, such as \"Name:app-production\". The value of tag is assumed to be a partial match",
	)
	flags.StringP("fields", "F", "InstanceId InstanceType PublicIpAddress PrivateIpAddress AvailabilityZone StateName LaunchTime Tag:Name Tag:Service 'Tag:In Charge'", "Output fields list separated by space")
	flags.StringP("domain", "D", "", "Please enter the domain you wish to search")
	viper.BindPFlag("ec2.vls.all", flags.Lookup("all"))
	viper.BindPFlag("ec2.vls.quiet", flags.Lookup("quiet"))
	viper.BindPFlag("ec2.vls.filter-tag", flags.Lookup("filter-tag"))
	viper.BindPFlag("ec2.vls.fields", flags.Lookup("fields"))
	viper.BindPFlag("ec2.vls.domain", flags.Lookup("domain"))

	return cmd
}

func runEC2VLsCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	options := myaws.EC2VLsOptions{
		All:       viper.GetBool("ec2.vls.all"),
		Quiet:     viper.GetBool("ec2.vls.quiet"),
		FilterTag: viper.GetString("ec2.vls.filter-tag"),
		Fields:    viper.GetStringSlice("ec2.vls.fields"),
		Domain:    viper.GetStringSlice("ec2.vls.domain"),
	}

	return client.EC2VLs(options)
}

func newEC2LsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List EC2 instances",
		RunE:  runEC2LsCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("all", "a", false, "List all instances (by default, list running instances only)")
	flags.BoolP("quiet", "q", false, "Only display InstanceIDs")
	flags.StringP("filter-tag", "t", "",
		"Filter instances by tag, such as \"Name:app-production\". The value of tag is assumed to be a partial match",
	)
	flags.StringP("fields", "F", "InstanceId InstanceType PublicIpAddress PrivateIpAddress AvailabilityZone StateName LaunchTime Tag:Name Tag:Service 'Tag:In Charge'", "Output fields list separated by space")
	flags.StringP("domain", "D", "", "Please enter the domain you wish to search")
	viper.BindPFlag("ec2.ls.all", flags.Lookup("all"))
	viper.BindPFlag("ec2.ls.quiet", flags.Lookup("quiet"))
	viper.BindPFlag("ec2.ls.filter-tag", flags.Lookup("filter-tag"))
	viper.BindPFlag("ec2.ls.fields", flags.Lookup("fields"))
	viper.BindPFlag("ec2.ls.domain", flags.Lookup("domain"))

	return cmd
}

func runEC2LsCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	options := myaws.EC2LsOptions{
		All:       viper.GetBool("ec2.ls.all"),
		Quiet:     viper.GetBool("ec2.ls.quiet"),
		FilterTag: viper.GetString("ec2.ls.filter-tag"),
		Fields:    viper.GetStringSlice("ec2.ls.fields"),
		Domain:    viper.GetStringSlice("ec2.ls.domain"),
	}

	return client.EC2Ls(options)
}

func newEC2StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start INSTANCE_ID [...]",
		Short: "Start EC2 instances",
		RunE:  runEC2StartCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("wait", "w", false, "Wait until instance running")

	viper.BindPFlag("ec2.start.wait", flags.Lookup("wait"))

	return cmd
}

func runEC2StartCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("INSTANCE_ID is required")
	}
	instanceIds := aws.StringSlice(args)

	options := myaws.EC2StartOptions{
		InstanceIds: instanceIds,
		Wait:        viper.GetBool("ec2.start.wait"),
	}

	return client.EC2Start(options)
}

func newEC2StopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop INSTANCE_ID [...]",
		Short: "Stop EC2 instances",
		RunE:  runEC2StopCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("wait", "w", false, "Wait until instance stopped")

	viper.BindPFlag("ec2.stop.wait", flags.Lookup("wait"))

	return cmd
}

func runEC2StopCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("INSTANCE_ID is required")
	}
	instanceIds := aws.StringSlice(args)

	options := myaws.EC2StopOptions{
		InstanceIds: instanceIds,
		Wait:        viper.GetBool("ec2.stop.wait"),
	}

	return client.EC2Stop(options)
}

func newEC2SSHCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssh [USER@]INSTANCE_NAME [COMMAND...]",
		Short: "SSH to EC2 instances",
		RunE:  runEC2SSHCmd,
	}

	flags := cmd.Flags()
	flags.StringP("login-name", "l", "", "Login username")
	flags.StringP("identity-file", "i", "~/.ssh/id_rsa", "SSH private key file")
	flags.BoolP("private", "", false, "Use private IP to connect")

	viper.BindPFlag("ec2.ssh.login-name", flags.Lookup("login-name"))
	viper.BindPFlag("ec2.ssh.identity-file", flags.Lookup("identity-file"))
	viper.BindPFlag("ec2.ssh.private", flags.Lookup("private"))

	return cmd
}

func runEC2SSHCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("Instance name is required")
	}

	var loginName, instanceName string
	if strings.Contains(args[0], "@") {
		// parse loginName@instanceName format
		splitted := strings.SplitN(args[0], "@", 2)
		loginName, instanceName = splitted[0], splitted[1]
	} else {
		loginName = viper.GetString("ec2.ssh.login-name")
		instanceName = args[0]
	}

	filterTag := "Name:" + instanceName

	var command string
	if len(args) >= 2 {
		command = strings.Join(args[1:], " ")
	}
	options := myaws.EC2SSHOptions{
		FilterTag:    filterTag,
		LoginName:    loginName,
		IdentityFile: viper.GetString("ec2.ssh.identity-file"),
		Private:      viper.GetBool("ec2.ssh.private"),
		Command:      command,
	}

	return client.EC2SSH(options)
}
