package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ChizhovVadim/counterdev/pkg/common"
	"github.com/ChizhovVadim/counterdev/pkg/engine"
	"github.com/ChizhovVadim/counterdev/pkg/evalnn"
	"github.com/ChizhovVadim/counterdev/pkg/uci"
)

/*
Counter Copyright (C) 2017-2026 Vadim Chizhov
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
You should have received a copy of the GNU General Public License along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

/* TODO
go:embed n-30-5268.nn
var content embed.FS
*/

func main() {
	var logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)

	//TODO content
	var weights, err = loadNetworkWeights(mapPath("~/chess/n-30-5268.nn"), true)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Loaded nnue weights")

	var options = engine.NewMainOptions(func() common.IEvaluator {
		return evalnn.NewEvaluationService(weights, 1.0)
	})
	var eng = engine.NewEngine(options)
	var protocol = uci.New("Counter", "Vadim Chizhov", "5.5", eng,
		[]uci.Option{
			&uci.IntOption{Name: "Hash", Min: 4, Max: 1 << 16, Value: &eng.Options.Hash},
			&uci.IntOption{Name: "Threads", Min: 1, Max: runtime.NumCPU(), Value: &eng.Options.Threads},
			//&uci.BoolOption{Name: "ExperimentSettings", Value: &eng.Options.ExperimentSettings},
		},
	)
	protocol.Run(logger)
}

func loadNetworkWeights(path string, oldFormat bool) (*evalnn.Weights, error) {
	var f, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return evalnn.LoadWeights(f, oldFormat)
}

func mapPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		curUser, err := user.Current()
		if err != nil {
			return path
		}
		return filepath.Join(curUser.HomeDir, strings.TrimPrefix(path, "~/"))
	}
	if strings.HasPrefix(path, "./") {
		var exePath, err = os.Executable()
		if err != nil {
			return path
		}
		return filepath.Join(filepath.Dir(exePath), strings.TrimPrefix(path, "./"))
	}
	return path
}
