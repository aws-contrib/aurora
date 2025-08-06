package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws-contrib/aurora/cmd"
	"github.com/aws-contrib/aurora/internal/database/ent"
	"github.com/aws-contrib/aurora/internal/database/ent/template"
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
						Name:  "init",
						Usage: "Creates the necessary tables in the database.",
						Action: func(ctx context.Context, command *cli.Command) error {
							gateway := command.Root().Metadata["gateway"].(ent.Gateway)
							return gateway.CreateTableRevisions(ctx)
						},
					},
					{
						Name:  "apply",
						Usage: "Applies pending migration files on the connected database.",
						Action: func(ctx context.Context, command *cli.Command) error {
							repository := &ent.MigrationRepository{
								Gateway:    command.Root().Metadata["gateway"].(ent.Gateway),
								FileSystem: os.DirFS(command.Root().Metadata["directory"].(string)),
							}

							lock := &ent.LockRevisionParams{
								Revision: ent.Mutex,
								Timeout:  1 * time.Minute,
							}
							// lock the execution
							if xerr := repository.LockRevision(ctx, lock); xerr != nil {
								return xerr
							}
							unlock := &ent.UnlockRevisionParams{
								Revision: ent.Mutex,
							}
							// unlock the execution
							defer repository.UnlockRevision(ctx, unlock)

							migrations, err := repository.ListMigrations(ctx, &ent.ListMigrationsParams{})
							if err != nil {
								return err
							}

							state := &ent.MigrationState{}
							// prepare the status
							for _, migration := range migrations {
								if state.Next == nil {
									state.Next = migration.Revision
								}

								if err == nil {
									if state.Current == nil || state.Current.Error == nil {
										params := &ent.ApplyMigrationParams{}
										params.Migration = migration
										// apply the migration
										err = repository.ApplyMigration(ctx, params)
										// update the migration state
										migration = params.Migration
									}
								}

								if migration.Revision.ExecutedAt.IsZero() {
									state.Pending = append(state.Pending, migration.Revision)
								} else {
									state.Executed = append(state.Executed, migration.Revision)
									state.Current = migration.Revision
									state.Next = nil
								}
							}

							// print the status
							template.Execute(os.Stdout, "status", state)

							if state.Current != nil && state.Current.Error != nil {
								// return the error
								return cli.Exit("", 1)
							}
							// done!
							return err
						},
					},
					{
						Name:  "status",
						Usage: "Get information about the current migration status.",
						Action: func(ctx context.Context, command *cli.Command) error {
							repository := &ent.MigrationRepository{
								Gateway:    command.Root().Metadata["gateway"].(ent.Gateway),
								FileSystem: os.DirFS(command.Root().Metadata["directory"].(string)),
							}

							migrations, err := repository.ListMigrations(ctx, &ent.ListMigrationsParams{})
							if err != nil {
								return err
							}

							state := &ent.MigrationState{}
							// prepare the status
							for _, migration := range migrations {
								if state.Next == nil {
									state.Next = migration.Revision
								}

								if migration.Revision.ExecutedAt.IsZero() {
									state.Pending = append(state.Pending, migration.Revision)
								} else {
									state.Executed = append(state.Executed, migration.Revision)
									state.Current = migration.Revision
									state.Next = nil
								}
							}
							// print the status
							template.Execute(os.Stdout, "status", state)
							// done!
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
