package cfgwatch

import (
	"encoding/json"

	"github.com/master-g/playground/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var paramList = []config.Flag{
	{Name: "singular", Type: config.Bool, Shorthand: "", Value: false, Usage: "Singular mode, if true, this node will not connect to other service"},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start service short",
	Long:  `start service long`,
	Run: func(cmd *cobra.Command, args []string) {
		j, _ := json.Marshal(viper.AllSettings())
		log.Info(string(j))

		configService()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	err := config.BindFlags(serveCmd, paramList)
	if err != nil {
		log.Fatal(err)
	}
}

func configService() {
	log.Info(viper.GetBool("singular"))
	go signal.Start()

	<-signal.InterruptChan
}

func updateConfig() {
	log.Infof("config update: %v", viper.GetBool("singular"))
}
