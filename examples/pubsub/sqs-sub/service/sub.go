package service

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/xhroot/gizmo/config"
	"github.com/xhroot/gizmo/pubsub"
	"github.com/NYTimes/logrotate"
	"github.com/Sirupsen/logrus"
	"github.com/cyberdelia/go-metrics-graphite"
	"github.com/golang/protobuf/proto"
	"github.com/rcrowley/go-metrics"

	"github.com/xhroot/gizmo/examples/nyt"
)

var (
	Log = logrus.New()

	sub pubsub.Subscriber

	client nyt.Client

	articles []nyt.SemanticConceptArticle
)

type Config struct {
	*config.Config
	MostPopularToken string
	SemanticToken    string
}

func Init() {
	flag.Parse()

	var cfg *Config
	config.LoadJSONFile("./config.json", &cfg)
	config.SetLogOverride(cfg.Log)

	if *cfg.Log != "" {
		lf, err := logrotate.NewFile(*cfg.Log)
		if err != nil {
			Log.Fatalf("unable to access log file: %s", err)
		}
		Log.Out = lf
		Log.Formatter = &logrus.JSONFormatter{}
	} else {
		Log.Out = os.Stderr
	}

	pubsub.Log = Log

	if cfg.GraphiteHost != nil {
		initMetrics(*cfg.GraphiteHost)
	}

	client = nyt.NewClient(cfg.MostPopularToken, cfg.SemanticToken)

	var err error
	sub, err = pubsub.NewSQSSubscriber(cfg.SQS)
	if err != nil {
		Log.Fatal("unable to init pb subs SQS: ", err)
	}
}

func Run() (err error) {
	stream := sub.Start()

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
		Log.Infof("received kill signal %s", <-ch)
		err = sub.Stop()
	}()

	var article nyt.SemanticConceptArticle
	for msg := range stream {
		totalMsgsConsumed.Inc(1)

		if err = proto.Unmarshal(msg.Message(), &article); err != nil {
			Log.Error("unable to unmarshal article from SQS: ", err)
			errorCount.Inc(1)
			if err = msg.Done(); err != nil {
				Log.Error("unable to delete message from SQS: ", err)
			}
			continue
		}

		// do something!
		fmt.Println("Most Recent Article on 'Cats':")
		out, _ := json.MarshalIndent(article, "", "    ")
		fmt.Fprint(os.Stdout, string(out))
		articles = append(articles, article)

		if err = msg.Done(); err != nil {
			Log.WithFields(logrus.Fields{
				"article": article,
			}).Error("unable to delete message from SQS: ", err)
		}
	}

	return err
}

var (
	errorCount        = metrics.NewRegisteredCounter("error-count", metrics.DefaultRegistry)
	totalMsgsConsumed = metrics.NewRegisteredCounter("total-consumed", metrics.DefaultRegistry)
)

func initMetrics(graphiteHost string) {
	Log.Infof("connecting to graphite host: %s", graphiteHost)
	addr, err := net.ResolveTCPAddr("tcp", graphiteHost)
	if err != nil {
		Log.Errorf("unable to resolve graphite host: %s", err)
		return
	}
	go graphite.Graphite(metrics.DefaultRegistry, 30*time.Second, metricsRegistryName(), addr)
}

func metricsRegistryName() string {
	// get only server base name
	name, _ := os.Hostname()
	name = strings.SplitN(name, ".", 2)[0]
	// set it up to be paperboy.servername
	name = strings.Replace(name, "-", ".", 1)
	// add the 'apps' prefix  to keep things neat
	return "apps." + name + ".cats-subscriber"
}
