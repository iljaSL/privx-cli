package cmd

import (
	"errors"
	"os"
	"strings"

	apiConfig "github.com/SSHcom/privx-sdk-go/api/config"
	"github.com/SSHcom/privx-sdk-go/api/hoststore"
	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/spf13/cobra"
)

var (
	hostID         string
	filter         string
	deployStatus   bool
	disabledStatus bool
)

func init() {
	rootCmd.AddCommand(hostListCmd)
	hostListCmd.Flags().IntVar(&offset, "offset", 0, "where to start fetching the items")
	hostListCmd.Flags().IntVar(&limit, "limit", 50, "number of items to return")
	hostListCmd.Flags().StringVar(&sortdir, "sortdir", "", "sort direction, ASC or DESC (default ASC)")
	hostListCmd.Flags().StringVar(&sortkey, "sortkey", "", "sort object by name, updated, or created.")
	hostListCmd.Flags().StringVar(&filter, "filter", "", "filter hosts, possible values: accessible or configured")

	hostListCmd.AddCommand(hostSearchCmd)
	hostSearchCmd.Flags().IntVar(&offset, "offset", 0, "where to start fetching the items")
	hostSearchCmd.Flags().IntVar(&limit, "limit", 50, "number of items to return")
	hostSearchCmd.Flags().StringVar(&sortdir, "sortdir", "", "sort direction, ASC or DESC (default ASC)")
	hostSearchCmd.Flags().StringVar(&sortkey, "sortkey", "", "sort object by name, updated, or created.")
	hostSearchCmd.Flags().StringVar(&filter, "filter", "", "filter hosts, possible values: accessible or configured")

	hostListCmd.AddCommand(hostCreateCmd)

	hostListCmd.AddCommand(hostShowCmd)
	hostShowCmd.Flags().StringVar(&hostID, "id", "", "host ID")
	hostShowCmd.MarkFlagRequired("id")

	hostListCmd.AddCommand(hostUpdateCmd)
	hostUpdateCmd.Flags().StringVar(&hostID, "id", "", "host ID")
	hostUpdateCmd.MarkFlagRequired("id")

	hostListCmd.AddCommand(hostDeleteCmd)
	hostDeleteCmd.Flags().StringVar(&hostID, "id", "", "host ID")
	hostDeleteCmd.MarkFlagRequired("id")

	hostListCmd.AddCommand(hostResolveCmd)

	hostListCmd.AddCommand(hostDeployableCmd)
	hostDeployableCmd.Flags().StringVar(&hostID, "id", "", "host ID")
	hostDeployableCmd.Flags().BoolVar(&deployStatus, "status", false, "host deploy status")
	hostDeployableCmd.MarkFlagRequired("id")
	hostDeployableCmd.MarkFlagRequired("status")

	hostListCmd.AddCommand(hostDisableCmd)
	hostDisableCmd.Flags().StringVar(&hostID, "id", "", "host ID")
	hostDisableCmd.Flags().BoolVar(&disabledStatus, "status", false, "host disabled status")
	hostDisableCmd.MarkFlagRequired("id")
	hostDisableCmd.MarkFlagRequired("status")

	hostListCmd.AddCommand(hostSettingListCmd)

	hostListCmd.AddCommand(hostsDeployCmd)
}

//
//
var hostListCmd = &cobra.Command{
	Use:   "hosts",
	Short: "List and manage PrivX hosts",
	Long:  `List and manage PrivX hosts`,
	Example: `
privx-cli hosts [access flags] --offset <OFFSET> --sortkey <SORTKEY>
	`,
	SilenceUsage: true,
	RunE:         hostList,
}

func hostList(cmd *cobra.Command, args []string) error {
	api := hoststore.New(curl())

	hosts, err := api.Hosts(offset, limit, sortkey, strings.ToUpper(sortdir), filter)
	if err != nil {
		return err
	}

	return stdout(hosts)
}

//
//
var hostSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search hosts",
	Long:  `Search hosts`,
	Example: `
privx-cli hosts search [access flags] --offset <OFFSET> --sortkey <SORTKEY>
privx-cli hosts search [access flags] --limit <LIMIT> JSON-FILE
	`,
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE:         hostSearch,
}

func hostSearch(cmd *cobra.Command, args []string) error {
	var searchObject hoststore.HostSearchObject
	api := hoststore.New(curl())

	if len(args) == 1 {
		err := decodeJSON(args[0], &searchObject)
		if err != nil {
			return err
		}
	}

	hosts, err := api.SearchHost(sortkey, strings.ToUpper(sortdir), filter, offset, limit, &searchObject)
	if err != nil {
		return err
	}

	return stdout(hosts)
}

//
//
var hostCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new host",
	Long:  `Create new host`,
	Example: `
privx-cli hosts create [access flags] JSON-FILE
	`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         hostCreate,
}

func hostCreate(cmd *cobra.Command, args []string) error {
	var newHost hoststore.Host
	api := hoststore.New(curl())

	err := decodeJSON(args[0], &newHost)
	if err != nil {
		return err
	}

	id, err := api.CreateHost(newHost)
	if err != nil {
		return err
	}

	return stdout(id)
}

//
//
var hostShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Get host by ID",
	Long:  `Get host by ID. Host ID's are separated by commas when using multiple values, see example`,
	Example: `
privx-cli hosts show [access flags] --id <HOST-ID>,<HOST-ID>
	`,
	SilenceUsage: true,
	RunE:         hostShow,
}

func hostShow(cmd *cobra.Command, args []string) error {
	api := hoststore.New(curl())
	hosts := []hoststore.Host{}

	for _, id := range strings.Split(hostID, ",") {
		host, err := api.Host(id)
		if err != nil {
			return err
		}
		hosts = append(hosts, *host)
	}

	return stdout(hosts)
}

//
//
var hostUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update host",
	Long:  `Update host`,
	Example: `
privx-cli hosts update [access flags] JSON-FILE --id <HOST-ID>
	`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         hostUpdate,
}

func hostUpdate(cmd *cobra.Command, args []string) error {
	var updateHost hoststore.Host
	api := hoststore.New(curl())

	err := decodeJSON(args[0], &updateHost)
	if err != nil {
		return err
	}

	err = api.UpdateHost(hostID, &updateHost)
	if err != nil {
		return err
	}

	return err
}

//
//
var hostDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete host",
	Long:  `Delete host. Host ID's are separated by commas when using multiple values, see example`,
	Example: `
privx-cli hosts delete [access flags] --id <HOST-ID>,<HOST-ID>
	`,
	SilenceUsage: true,
	RunE:         hostDelete,
}

func hostDelete(cmd *cobra.Command, args []string) error {
	api := hoststore.New(curl())

	for _, id := range strings.Split(hostID, ",") {
		err := api.DeleteHost(id)
		if err != nil {
			return err
		}
	}

	return nil
}

//
//
var hostResolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Resolve host",
	Long:  `Resolve service and address to a single host`,
	Example: `
privx-cli hosts resolve [access flags] JSON-FILE
	`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         hostResolve,
}

func hostResolve(cmd *cobra.Command, args []string) error {
	var service hoststore.Service
	api := hoststore.New(curl())

	err := decodeJSON(args[0], &service)
	if err != nil {
		return err
	}

	host, err := api.ResolveHost(service)
	if err != nil {
		return err
	}

	return stdout(host)
}

//
//
var hostDeployableCmd = &cobra.Command{
	Use:   "deployable",
	Short: "Set a host to be depoyable or undeployable",
	Long:  `Set a host to be depoyable or undeployable. Host ID's are separated by commas when using multiple values, see example`,
	Example: `
privx-cli hosts deployable [access flags] --id <HOST-ID>,<HOST-ID> --status=true
	`,
	SilenceUsage: true,
	RunE:         hostDeployable,
}

func hostDeployable(cmd *cobra.Command, args []string) error {
	api := hoststore.New(curl())

	for _, id := range strings.Split(hostID, ",") {
		err := api.UpdateDeployStatus(id, deployStatus)
		if err != nil {
			return err
		}
	}

	return nil
}

//
//
var hostDisableCmd = &cobra.Command{
	Use:   "disabled",
	Short: "Enable/disable host",
	Long:  `Enable(false)/disable(true) host. Host ID's are separated by commas when using multiple values, see example`,
	Example: `
privx-cli hosts disabled [access flags] --id <HOST-ID>,<HOST-ID> --status=true
	`,
	SilenceUsage: true,
	RunE:         hostDisable,
}

func hostDisable(cmd *cobra.Command, args []string) error {
	api := hoststore.New(curl())

	for _, id := range strings.Split(hostID, ",") {
		err := api.UpdateDisabledHostStatus(id, disabledStatus)
		if err != nil {
			return err
		}
	}

	return nil
}

//
//
var hostSettingListCmd = &cobra.Command{
	Use:   "settings",
	Short: "Get the default service options",
	Long:  `Get the default service options`,
	Example: `
privx-cli hosts settings [access flags]
	`,
	SilenceUsage: true,
	RunE:         hostSettingList,
}

func hostSettingList(cmd *cobra.Command, args []string) error {
	api := hoststore.New(curl())

	settings, err := api.ServiceOptions()
	if err != nil {
		return err
	}

	return stdout(settings)
}

//
//
var hostsDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Creates target hosts deployment config",
	Long:  `Creates target hosts deployment config`,
	Example: `
privx-cli hosts deploy [access flags] <NAME>
	`,
	SilenceUsage: true,
	RunE:         hostDeploy,
}

func hostDeploy(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires name of deployment configuration as an argument")
	}
	name := args[0]

	curl := curl()
	store := userstore.New(curl)

	seq, err := store.TrustedClients()
	if err != nil {
		return err
	}

	cli := findClientID(seq, name)
	if cli == "" {
		cli, err = store.CreateTrustedClient(
			userstore.HostProvisioning(name),
		)
		if err != nil {
			return err
		}
	}

	conf := apiConfig.New(curl)
	file, err := conf.ConfigDeploy(cli)
	if err != nil {
		return err
	}

	os.Stdout.Write(file)
	return nil
}

func findClientID(seq []userstore.TrustedClient, name string) string {
	for _, cli := range seq {
		if cli.Name == name {
			return cli.ID
		}
	}
	return ""
}
