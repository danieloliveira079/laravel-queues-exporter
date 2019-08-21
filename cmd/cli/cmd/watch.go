// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/ui/widget"
	"github.com/spf13/cobra"
	"log"
	"runtime"
	"time"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch queues metrics from remote server",
	Long:  `Queues metrics are retrieved from the remote server in a time bases interval.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())

		watcher, err := NewWatcher(serverAddress, time.Duration(2)*time.Second)
		if err != nil {
			log.Fatal(err)
		}
		watcher.Run()

		dashboard := widget.NewDashboardWidget(time.Duration(1)*time.Second, watcher.GetMetrics)
		dashboard.Render(ctx, cancel)

		<-ctx.Done()
		runtime.GC()
		log.Println("stopped")
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
