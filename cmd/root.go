package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/webitel/wlog"
	"golang.org/x/sync/errgroup"

	"github.com/kirychukyurii/fd-import/config"
	"github.com/kirychukyurii/fd-import/models"
	"github.com/kirychukyurii/fd-import/pkg/db"
	"github.com/kirychukyurii/fd-import/pkg/filestorage"
	"github.com/kirychukyurii/fd-import/pkg/s3"
)

var (

	// version is the app's semantic version.
	version = "0.0.0"

	// commit is the git commit used to build the App.
	commit     = "hash"
	commitDate = "date"

	requesterNameRegexp        = regexp.MustCompile(`^/([^/]+)`)
	attachmentRegexp           = regexp.MustCompile(`^/.*/(\d+)/attachments/(\d+)-.*\.(.*)$`)
	attachmentWithoutExtRegexp = regexp.MustCompile(`^/.*/(\d+)/attachments/(\d+)-.*$`)
)

func Execute() {
	if err := command().Execute(); err != nil {
		os.Exit(-1)
	}
}

func command() *cobra.Command {
	log := wlog.NewLogger(&wlog.LoggerConfiguration{
		EnableConsole: true,
		ConsoleLevel:  "debug",
	})

	cfg := config.New()
	c := &cobra.Command{
		Use:          "fd-import",
		Short:        "FreshDesk Import - easy import .json files, exported from API",
		SilenceUsage: true,
		Version:      fmt.Sprintf("%s, commit %s, date %s", version, commit, commitDate),
		Args: func(cmd *cobra.Command, args []string) error {

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// apply CLI args to config
			if err := cmd.ParseFlags(os.Args[1:]); err != nil {
				return fmt.Errorf("parsing flags: %w", err)
			}

			log := wlog.NewLogger(&wlog.LoggerConfiguration{
				EnableConsole: true,
				ConsoleLevel:  "info",
				EnableFile:    true,
				FileLevel:     cfg.LogLevel,
				FileLocation:  cfg.LogFile,
			})

			// os.Interrupt to gracefully shutdown on Ctrl+C which is SIGINT
			// syscall.SIGTERM is the usual signal for termination and the default one (it can be modified)
			// for docker containers, which is also used by kubernetes.
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer cancel()

			dbpool, err := db.New(ctx, log, cfg.DSN)
			if err != nil {
				return err
			}

			cfg.AttachmentDir = filepath.Join(cfg.AttachmentDir, cfg.Domain)
			a := &app{
				log:    log,
				cfg:    cfg,
				dbpool: dbpool,
				bucket: s3.New(log, cfg.S3),
				keys:   make(map[string][]string),
				stats: &stats{
					processed:   atomic.Uint64{},
					exists:      atomic.Uint64{},
					tickets:     atomic.Uint64{},
					attachments: atomic.Uint64{},
				},
			}

			// This blocks until the context is finished or until an error is produced
			if err = a.run(ctx); err != nil {
				a.log.Error("run app", wlog.Err(err))
			}

			a.log.Info("processed items", wlog.Any("processed", a.stats.processed.Load()),
				wlog.Any("exists", a.stats.exists.Load()), wlog.Any("tickets", a.stats.tickets.Load()),
				wlog.Any("attachments", a.stats.attachments.Load()))

			return err
		},
	}

	flagSet(c.PersistentFlags(), cfg)
	c.AddCommand(migrateCommand(cfg, log))

	return c
}

// flagSet sets up the command line flags using the given FlagSet and Config objects.
// The function binds the flags to the corresponding fields in the Config object.
// The flag names, shorthand flags, default values, and usage descriptions are specified
// for each flag.
func flagSet(fs *pflag.FlagSet, cfg *config.Config) {
	fs.StringVarP(&cfg.LogLevel, "log-level", "l", "debug", "log level")
	fs.StringVar(&cfg.LogFile, "log-file", "./fd-import.log", "log file")
	fs.StringVarP(&cfg.ExportedPath, "path", "p", "./export-data", "base path to exported files")
	fs.StringVarP(&cfg.DSN, "dsn", "d", "", "database connection string")
	fs.StringVar(&cfg.Domain, "domain", "", "domain name")
	fs.StringVarP(&cfg.AttachmentDir, "attachment", "a", "./attachments", "directory to store attachment files")
	fs.IntVarP(&cfg.Workers, "workers-count", "w", 100, "number of concurrent workers")

	fs.StringVar(&cfg.S3.AccessKeyID, "s3.access-key", "", "S3 access key ID")
	fs.StringVar(&cfg.S3.SecretAccessKey, "s3.secret-key", "", "S3 secret access key")
	fs.StringVar(&cfg.S3.Region, "s3.region", "", "S3 region")
	fs.StringVar(&cfg.S3.Bucket, "s3.bucket", "", "S3 bucket")
}

type app struct {
	log *wlog.Logger
	cfg *config.Config

	dbpool *db.Connection
	bucket *s3.Bucket

	domain int64
	keys   map[string][]string
	stats  *stats
}

type stats struct {
	processed   atomic.Uint64
	exists      atomic.Uint64
	tickets     atomic.Uint64
	attachments atomic.Uint64
}

// run executes the main logic of the application. It performs the following steps:
//   - Fetches the domain ID from the database pool based on the given domain name.
//   - If the domain doesn't exist, creates a new domain in the database pool.
//   - Sets the retrieved or created domain ID as the app's domain.
//   - Uses an errgroup to concurrently process objects in the object pool.
//   - Processes each object by calling the process method of the app.
//   - Waits for all processing to complete.
//   - Logs the keys found during processing.
func (a *app) run(ctx context.Context) error {
	domain, err := a.dbpool.Domain(ctx, a.cfg.Domain)
	if err != nil {
		if errors.Is(err, db.ErrDBNoExists) {
			var cerr error
			domain, cerr = a.dbpool.CreateDomain(ctx, a.cfg.Domain)
			if cerr != nil {
				return fmt.Errorf("create domain: %v", cerr)
			}
		} else {
			return fmt.Errorf("domain: %v", err)
		}
	}

	a.domain = domain
	eg, gctx := errgroup.WithContext(ctx)
	workers := a.cfg.Workers
	eg.SetLimit(workers + 1)
	eg.Go(func() error {
		if err := a.bucket.ListObjects(gctx, a.cfg.ExportedPath, ""); err != nil {
			return err
		}

		a.log.Debug("complete list objects")

		return nil
	})

	objpool := a.bucket.ObjectPool()
	for o := range objpool {
		eg.Go(func() error {
			a.bucket.DequeueObjectPool()
			a.log.Debug("process", wlog.Any("object", o))
			if err := a.process(gctx, o); err != nil {
				return fmt.Errorf("process key (%s): %v", o, err)
			}

			return nil
		})
	}

	a.log.Debug("wait for all")
	if err := eg.Wait(); err != nil {
		return err
	}

	/*
		for o, k := range a.keys {
			a.log.Debug(fmt.Sprintf("found %s keys", o), wlog.String("keys", fmt.Sprintf("%v", unique(k))))
		}
	*/

	return nil
}

// process executes the processing logic for the given key. It performs the following steps:
//   - Checks if the ticket already exists in the database for the given domain and key. If so, returns without further processing.
//   - Retrieves the metadata of the S3 object using the `HeadObject` method of the bucket.
//   - Logs the metadata of the S3 object.
//   - Checks the content type of the S3 object and performs the appropriate processing based on the content type.
//   - If the content type is "application/json", calls the `processJSON` method to process the JSON object.
//   - Otherwise, calls the `processAttachment` method to process the attachment.
//   - Returns any processing errors that occur.
func (a *app) process(ctx context.Context, key string) error {
	a.stats.processed.Add(1)
	ok, err := a.dbpool.Ticket(ctx, a.domain, key)
	if err != nil {
		if !errors.Is(err, db.ErrDBNoExists) {
			return err
		}
	}

	if ok {
		a.stats.exists.Add(1)
		a.log.Debug("exists", wlog.Any("key", key))

		return nil
	}

	if strings.Contains(key, "/attachments/") {
		a.stats.attachments.Add(1)
		if err := a.processAttachment(ctx, key); err != nil {
			return fmt.Errorf("attachment: %v", err)
		}
	} else {
		a.stats.tickets.Add(1)
		if err := a.processJSON(ctx, key); err != nil {
			return fmt.Errorf("json: %v", err)
		}
	}

	/*
		head, err := a.bucket.HeadObject(ctx, key)
		if err != nil {
			return fmt.Errorf("head object: %v", err)
		}

		a.log.With(wlog.Any("content-type", head.ContentType), wlog.Any("key", key)).
			Debug("head S3 object")

		ct := *head.ContentType
		switch ct {
		case "application/json":
			if err := a.processJSON(ctx, key); err != nil {
				return fmt.Errorf("json: %v", err)
			}
		default:
			if err := a.processAttachment(ctx, key); err != nil {
				return fmt.Errorf("attachment: %v", err)
			}
		}
	*/

	return nil
}

// processJSON processes the JSON object with the given key. It performs the following steps:
//   - Retrieves the object from the bucket using the ReadObject method.
//   - Unmarshals the JSON object into a models.Ticket struct.
//   - Sets additional fields of the Ticket struct.
//   - Creates the ticket in the database using the CreateTicket method of the dbpool.
//   - Returns any error that occurs during the processing.
func (a *app) processJSON(ctx context.Context, key string) error {
	object, err := a.bucket.ReadObject(ctx, key)
	if err != nil {
		return fmt.Errorf("read object: %v", err)
	}

	var target models.Ticket
	if err = json.Unmarshal(object, &target); err != nil {
		return fmt.Errorf("unmarshal: %v", err)
	}

	match := requesterNameRegexp.FindStringSubmatch(strings.TrimPrefix(key, a.cfg.ExportedPath))
	target.AWSKey = key
	target.Raw = object
	target.DomainID = a.domain
	target.RequesterName = match[1]

	if err := a.dbpool.CreateTicket(ctx, &target); err != nil {
		return fmt.Errorf("create ticket: %v", err)
	}

	/*
		mar, err := json.Marshal(target)
		if err != nil {
			return err
		}

		patch, err := jsondiff.CompareJSON(object, mar)
		if err != nil {
			return err
		}

		for _, p := range patch {
			switch p.Type {
			case jsondiff.OperationRemove:
				a.keys["remove"] = append(a.keys["remove"], strings.TrimLeft(p.Path, "/"))
			case jsondiff.OperationAdd:
				a.keys["add"] = append(a.keys["add"], strings.TrimLeft(p.Path, "/"))
			}
		}
	*/

	return nil
}

func (a *app) processAttachment(ctx context.Context, key string) error {
	var (
		ticketID     string
		attachmentID string
		extension    string
		file         string
		fileName     string
		err          error
	)

	f := attachmentRegexp.FindStringSubmatch(strings.TrimPrefix(key, a.cfg.ExportedPath))
	if len(f) < 4 {
		f = attachmentWithoutExtRegexp.FindStringSubmatch(strings.TrimPrefix(key, a.cfg.ExportedPath))
		ticketID = f[1]
		attachmentID = f[2]
		fileName = attachmentID
	} else {
		ticketID = f[1]
		attachmentID = f[2]
		extension = f[3]
		fileName = fmt.Sprintf("%s.%s", attachmentID, extension)
	}

	attachmentPath := filepath.Join(a.cfg.AttachmentDir, ticketID)
	file = filepath.Join(attachmentPath, fileName)
	if filestorage.IsExist(file) {
		a.log.Debug("exists", wlog.String("key", key), wlog.String("file", file))

		return nil
	}

	defer func() {
		if err != nil {
			if filestorage.IsExist(file) {
				if err := filestorage.Remove(file); err != nil {
					a.log.Debug("remove file", wlog.Err(err), wlog.String("file", file))

					return
				}
			}
		}
	}()

	if err = filestorage.InsureDir(attachmentPath); err != nil {
		return err
	}

	if err = a.bucket.DownloadObject(ctx, key, file); err != nil {
		return fmt.Errorf("download: %v", err)
	}

	return nil
}

// unique returns a new slice containing unique elements from the input slice.
func unique(intSlice []string) []string {
	var list []string

	keys := make(map[string]bool)
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
