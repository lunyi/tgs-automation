package main

import (
	"cdnetwork/internal/log"
	"cdnetwork/pkg/cli"
	"fmt"

	"github.com/spf13/cobra"
)

func main() {

	// domains := namecheap.GetExpiredDomains()
	// domainsForExcel := postgresql.GetAgents(domains)
	// fmt.Println("domain's length:", len(domainsForExcel))

	//googlesheet.CreateExpiredDomainExcel(nil)

	//cloudflare.GetDnsInfo("autotest1.rpgp.cc")

	// var rootCmd = &cobra.Command{
	// 	Use:   "cdnctl",
	// 	Short: "cdnctl is a command-line tool for managing CDNetworks",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		fmt.Println("Usage: cdnctl <command> [options]")
	// 	},
	// }

	// addCommands(rootCmd)

	// if err := rootCmd.Execute(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
}

func addCommands(rootCmd *cobra.Command) {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "create is used to create resources",
	}
	getCreateCommand(createCmd)
	createCertificateCommand(createCmd)
	rootCmd.AddCommand(createCmd)

	rootCmd.AddCommand(getDeleteCommand())
	rootCmd.AddCommand(getDisableCommand())
	rootCmd.AddCommand(getCertificateCommand())
	rootCmd.AddCommand(getEnableDomainCommand())
	rootCmd.AddCommand(getDomainIdCommand())
}

func getCreateCommand(createCmd *cobra.Command) {
	var cmd = &cobra.Command{
		Use:   "domain [domain] [orginset]",
		Short: "Create a new domain and its originset",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Creating domain: %s, originset: %s \n", args[0], args[1])
			cdnCommand := &cli.CdnCommand{
				DomainName: args[0],
				OriginSet:  args[1],
			}
			cdnCommand.CreateDomain()
		},
	}
	createCmd.AddCommand(cmd)
}

func getDeleteCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete <domain>",
		Short: "Delete an existing domain",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Deleting domain: %s\n", args[0])
			cdnCommand := &cli.CdnCommand{
				DomainName: args[0],
			}
			cdnCommand.DeleteDomain()
		},
	}
	return cmd
}

func getDisableCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "disable <domain>",
		Short: "Disable an existing domain",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.LogInfo(fmt.Sprintf("Disabling domain: %s\n", args[0]))
			cdnCommand := &cli.CdnCommand{
				DomainName: args[0],
			}
			cdnCommand.DisableDomain()
		},
	}
	return cmd
}

func getEnableDomainCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "enable <domain>",
		Short: "Disable an existing domain",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.LogInfo(fmt.Sprintf("Enabling domain: %s\n", args[0]))
			cdnCommand := &cli.CdnCommand{
				DomainName: args[0],
			}
			cdnCommand.EnableDomain()
		},
	}
	return cmd
}

func getCertificateCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "getcert <domains>",
		Short: "Get certificate of domains",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cdnCommand := &cli.CdnCommand{
				DomainName: args[0],
			}
			cert := cdnCommand.GetCertificate()
			fmt.Printf("Cert is %s\n", cert)
		},
	}
	return cmd
}

func createCertificateCommand(createCmd *cobra.Command) {
	var cmd = &cobra.Command{
		Use:   "cert [domain]",
		Short: "Create certificate of domain",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cdnCommand := &cli.CdnCommand{
				DomainName: args[0],
			}
			cdnCommand.CreateCertificate()
		},
	}
	createCmd.AddCommand(cmd)
}

func getDomainIdCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "getdomainid <domains>",
		Short: "Get domain's id",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cdnCommand := &cli.CdnCommand{
				DomainName: args[0],
			}
			res := cdnCommand.GetDomainId()
			fmt.Printf("Domain id: %d\n", res)
		},
	}
	return cmd
}
