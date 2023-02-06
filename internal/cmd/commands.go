package cmd

import (
	"context"
	"os"

	"github.com/kairos-io/immucore/internal/utils"
	"github.com/kairos-io/immucore/pkg/mount"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spectrocloud-labs/herd"
	"github.com/twpayne/go-vfs"
	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{

	{
		Name:      "start",
		Usage:     "start",
		UsageText: "starts",
		Description: `
Sends a generic event payload with the configuration found in the scanned directories.
`,
		Aliases: []string{},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:     "dry-run",
				EnvVars:  []string{"IMMUCORE_DRY_RUN"},
				Required: false,
			},
		},
		Action: func(c *cli.Context) (err error) {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
			g := herd.DAG()
			s := &mount.State{Logger: log.Logger, Rootdir: "/"}

			fs := vfs.OSFS

			err = s.Register(g)
			if err != nil {
				s.Logger.Err(err)
				return err
			}

			log.Print(s.WriteDAG(g))

			if c.Bool("dry-run") {
				return err
			}

			cdBoot, err := utils.BootedFromCD(fs)
			if err != nil {
				s.Logger.Err(err)
				return err
			}

			if cdBoot {
				log.Info().Msg("Seems we booted from CD, doing nothing. Bye!")
				return nil
			}

			log.Print("Calling dag")
			return g.Run(context.Background())
		},
	},
}

func writeDag(d [][]herd.GraphEntry) {
	for i, layer := range d {
		log.Printf("%d.", (i + 1))
		for _, op := range layer {
			if op.Error != nil {
				log.Printf(" <%s> (error: %s) (background: %t) (weak: %t)", op.Name, op.Error.Error(), op.Background, op.WeakDeps)
			} else {
				log.Printf(" <%s> (background: %t) (weak: %t)", op.Name, op.Background, op.WeakDeps)
			}
		}
		log.Print("")
	}
}
