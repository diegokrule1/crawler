/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/diegokrule1/crawler/walker"
	"github.com/spf13/cobra"
	"log"
	"net/url"
)

// crawlerCmd represents the crawler command
var crawlerCmd = &cobra.Command{
	Use:   "crawler",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalln("You must provide the seed url to start crawling")
			return
		}

		consumer, producer, err := walker.Init()
		if err != nil {
			log.Fatalf("Could not create app %v", err)
			return
		}
		go consumer.Consume()
		//producer, b := cmd.Context().Value("producer").(walker.Producer)
		//
		//if !b {
		//	log.Fatalln("Could not get context")
		//	return
		//}

		urlToCrawl := args[0]

		parsedUrl, err := url.ParseRequestURI(urlToCrawl)
		if err != nil {
			log.Fatalf("%s is not a valid url", urlToCrawl)
		} else {
			log.Printf("Received valid url %s", parsedUrl)
		}

		producer.Produce(parsedUrl.Scheme, parsedUrl.Host, parsedUrl.Path, nil)

		<-producer.KillChan

		producer.Logger.Info("Exiting")

	},
}

func init() {
	rootCmd.AddCommand(crawlerCmd)
}
