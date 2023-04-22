package cmd

import (
	"github.com/huoxue1/qinglong-go/controller"
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/service"
	env_check "github.com/huoxue1/qinglong-go/utils/env-check"
	"github.com/spf13/cobra"
)

var (
	address string

	rootCMd = &cobra.Command{
		Use:     "qinglong",
		Aliases: []string{"ql"},
		PreRun: func(cmd *cobra.Command, args []string) {
			models.InitModels()
			initCron()
			env_check.CheckStatic()
			service.AppInit()
		},

		Run: func(cmd *cobra.Command, args []string) {

			engine := controller.Router()
			_ = engine.Run(address)
		},
		PostRun: func(cmd *cobra.Command, args []string) {

		},
	}
)

func Run() {
	rootCMd.PersistentFlags().StringVarP(&address, "address", "a", "0.0.0.0:5700", "the ql listen address!")
	InitLog()
	if err := rootCMd.Execute(); err != nil {
		return
	}
}
