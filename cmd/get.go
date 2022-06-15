/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/hacker65536/findlb/pkg/myaws"
	"github.com/hacker65536/findlb/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			log.SetLevel(log.DebugLevel)
		}
		domain, host := util.GetDomain(args[0])

		dns := ""
		//	flg, _ := cmd.Flags().GetBool("private")

		zones := myaws.GetHostedZone(domain)

		if len(zones) > 1 {
			for _, v := range zones {
				log.WithFields(
					log.Fields{
						"name":    v.Name,
						"zoneid":  v.Id,
						"private": v.Private,
						"records": v.Records,
						//					"records" v.Records,
					}).Warn()
			}
		} else {
			dns = myaws.GetDNSFromRecoard(zones[0].Id, host)

		}

		arn := myaws.GetALB(dns)

		fmt.Println(arn)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getCmd.Flags().BoolP("private", "p", false, "private zone")
}
