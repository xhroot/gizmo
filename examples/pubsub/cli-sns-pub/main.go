package main

import (
	"github.com/xhroot/gizmo/config"
	"github.com/xhroot/gizmo/examples/nyt"
	"github.com/xhroot/gizmo/pubsub"
	"github.com/Sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfigFromEnv()

	pub, err := pubsub.NewSNSPublisher(cfg.SNS)
	if err != nil {
		pubsub.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("unable to init publisher")
	}

	catArticle := &nyt.SemanticConceptArticle{
		Title:  "It's a Cat World",
		Byline: "By JP Robinson",
		Url:    "http://www.nytimes.com/2015/11/25/its-a-cat-world",
	}

	err = pub.Publish(catArticle.Url, catArticle)
	if err != nil {
		pubsub.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("unable to publish message")
	}

	pubsub.Log.WithFields(logrus.Fields{
		"articles": catArticle,
	}).Info("successfully published cat article")
}
