package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws-contrib/aurora/cmd"
	"github.com/aws-contrib/aurora/internal/database/ent"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "aurora",
		Usage: "Manage your database schema as code",
		Commands: []*cli.Command{
			{
				Name:  "migrate",
				Usage: "Manage versioned migration files",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Usage:    "select config (project) file using URL format",
						Value:    "file://aurora.hcl",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "env",
						Usage:    "set which env from the config file to use",
						Required: true,
					},
				},
				Before: func(ctx context.Context, command *cli.Command) (context.Context, error) {
					path, err := cmd.GetPath(command.String("config"))
					if err != nil {
						return nil, err
					}

					data, err := os.ReadFile(path)
					if err != nil {
						return nil, err
					}

					config := &cmd.Config{}
					if err := config.UnmarshalText(data); err != nil {
						return nil, err
					}

					if name := command.String("env"); name != "" {
						command.Root().Metadata = make(map[string]any)

						if config := config.GetEnvironment(name); config != nil {
							conn, err := config.GetURL()
							if err != nil {
								return nil, err
							}

							directory, err := config.Migration.GetDir()
							if err != nil {
								return nil, err
							}

							gateway, err := ent.Open(ctx, conn)
							if err != nil {
								return nil, err
							}

							command.Root().Metadata["directory"], err = cmd.GetPath(directory)
							if err != nil {
								return nil, err
							}

							if err := gateway.CreateTableRevisions(ctx); err != nil {
								return nil, err
							}

							command.Root().Metadata["gateway"] = gateway

						} else {
							return nil, fmt.Errorf("environment %s not found in config", name)
						}
					}

					return ctx, nil
				},
				Commands: []*cli.Command{
					{
						Name:  "apply",
						Usage: "Applies pending migration files on the connected database.",
						Action: func(ctx context.Context, command *cli.Command) error {
							repository := &ent.RevisionRepository{
								Gateway:    command.Root().Metadata["gateway"].(ent.Gateway),
								FileSystem: os.DirFS(command.Root().Metadata["directory"].(string)),
							}

							lock := &ent.LockRevisionParams{
								Revision: ent.Mutex,
								Timeout:  1 * time.Minute,
							}
							fmt.Println("Applying migrations started...")
							// lock the execution
							if err := repository.LockRevision(ctx, lock); err != nil {
								return err
							}

							time.Sleep(30 * time.Second)

							unlock := &ent.UnlockRevisionParams{
								Revision: ent.Mutex,
							}
							fmt.Println("Applying migrations completed...")
							// unlock the execution
							if err := repository.UnlockRevision(ctx, unlock); err != nil {
								return err
							}

							return nil
						},
					},
					{
						Name:  "status",
						Usage: "Get information about the current migration status.",
						Action: func(ctx context.Context, command *cli.Command) error {
							repository := &ent.RevisionRepository{
								Gateway:    command.Root().Metadata["gateway"].(ent.Gateway),
								FileSystem: os.DirFS(command.Root().Metadata["directory"].(string)),
							}

							revisions, err := repository.ListRevisions(ctx, &ent.ListRevisionsParams{})
							if err != nil {
								return err
							}

							for _, revision := range revisions {
								fmt.Println("Current migration status:", revision.GetName())
							}

							return nil
						},
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
