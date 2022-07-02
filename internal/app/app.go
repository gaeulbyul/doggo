package app

import (
	"math/rand"
	"time"

	"github.com/miekg/dns"
	"github.com/mr-karan/doggo/pkg/models"
	"github.com/mr-karan/doggo/pkg/resolvers"
	"github.com/mr-karan/logf"
)

// App represents the structure for all app wide configuration.
type App struct {
	Logger       *logf.Logger
	Version      string
	QueryFlags   models.QueryFlags
	Questions    []dns.Question
	Resolvers    []resolvers.Resolver
	ResolverOpts resolvers.Options
	Nameservers  []models.Nameserver
}

// NewApp initializes an instance of App which holds app wide configuration.
func New(logger *logf.Logger, buildVersion string) App {
	app := App{
		Logger:  logger,
		Version: buildVersion,
		QueryFlags: models.QueryFlags{
			QNames:      []string{},
			QTypes:      []string{},
			QClasses:    []string{},
			Nameservers: []string{},
		},
		Nameservers: []models.Nameserver{},
	}
	return app
}

// Attempts a DNS Lookup with retries for a given Question.
func (app *App) LookupWithRetry(attempts int, resolver resolvers.Resolver, ques dns.Question) (resolvers.Response, error) {
	resp, err := resolver.Lookup(ques)
	if err != nil {
		// Retry lookup.
		attempts--
		if attempts > 0 {
			// Add some random delay.
			time.Sleep(time.Millisecond*300 + (time.Duration(rand.Int63n(int64(time.Millisecond*100))))/2)
			app.Logger.Debug("retrying lookup")
			return app.LookupWithRetry(attempts, resolver, ques)
		}
		return resolvers.Response{}, err
	}
	return resp, nil
}
