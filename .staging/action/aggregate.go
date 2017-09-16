package actions

/*
import (
	"fmt"
	"flag"
	"log"
	"os"
	"time"

	"github.com/aviddiviner/gin-limit"
	"github.com/gin-gonic/contrib/cache"
	"github.com/gin-gonic/contrib/secure"
	"github.com/gin-gonic/gin"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/config/viper"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/logging/gologging"
	"github.com/devopsfaith/krakend/proxy"
	krakendgin "github.com/devopsfaith/krakend/router/gin"	
	// "strconv"
	// "strings"
	// "github.com/blevesearch/bleve"
	"github.com/roscopecoltran/sniperkit-limo/config"
	// "github.com/roscopecoltran/sniperkit-limo/model"
	"github.com/spf13/cobra"
)

// AggregateCmd does a full-text gateway
var AggregateCmd = &cobra.Command{
	Use:     "aggregate <vcs uri>",
	Aliases: []string{"collect"},
	Short:   "Aggregate info on stars",
	Long:    "Perform a full-text aggregate search on your stars",
	Example: fmt.Sprintf("  %s aggregate robust", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {

		// Get configuration
		cfg, err := getConfiguration()
		fatalOnError(err)

		serviceConfig := config

		// cfg.Debug = cfg.Debug
		if cfg.Aggregate.Port != 0 {
			serviceConfig.Port = *port
		}

		if cfg.Aggregate.Config.FilePath != "" {
			serviceConfig.Port = *port
		}

		logger, err := gologging.NewLogger(*logLevel, os.Stdout, "[KRAKEND]")
		if err != nil {
			log.Fatal("ERROR:", err.Error())
		}

		store := cache.NewInMemoryStore(time.Minute)

		mws := []gin.HandlerFunc{
			secure.Secure(secure.Options{
				AllowedHosts:          []string{"127.0.0.1:8080", "example.com", "ssl.example.com"},
				SSLRedirect:           false,
				SSLHost:               "ssl.example.com",
				SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
				STSSeconds:            315360000,
				STSIncludeSubdomains:  true,
				FrameDeny:             true,
				ContentTypeNosniff:    true,
				BrowserXssFilter:      true,
				ContentSecurityPolicy: "default-src 'self'",
			}),
			limit.MaxAllowed(20),
		}

		// routerFactory := krakendgin.DefaultFactory(proxy.DefaultFactory(logger), logger)

		routerFactory := krakendgin.NewFactory(krakendgin.Config{
			Engine:       gin.Default(),
			ProxyFactory: customProxyFactory{logger, proxy.DefaultFactory(logger)},
			Middlewares:  mws,
			Logger:       logger,
			HandlerFactory: func(configuration *config.EndpointConfig, proxy proxy.Proxy) gin.HandlerFunc {
				return cache.CachePage(store, configuration.CacheTTL, krakendgin.EndpointHandler(configuration, proxy))
			},
		})

		routerFactory.New().Run(serviceConfig)

	},
}

// customProxyFactory adds a logging middleware wrapping the internal factory
type customProxyFactory struct {
	logger  logging.Logger
	factory proxy.Factory
}

// New implements the Factory interface
func (cf customProxyFactory) New(cfg *config.EndpointConfig) (p proxy.Proxy, err error) {
	p, err = cf.factory.New(cfg)
	if err == nil {
		p = proxy.NewLoggingMiddleware(cf.logger, cfg.Endpoint)(p)
	}
	return
}

func init() {
	RootCmd.AddCommand(AggregateCmd)
}

*/
