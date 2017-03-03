package gitstats

import (
	"os"
	"strings"
	"testing"
	"time"

	"context"

	"github.com/google/go-github/github"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGitstatsPlugin(t *testing.T) {
	Convey("Meta should return metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, Name)
		So(meta.Version, ShouldResemble, Version)
		So(meta.Type, ShouldResemble, plugin.CollectorPluginType)
	})

	Convey("Create Gitstats Collector", t, func() {
		collector := &Gitstats{}
		Convey("So Gitstats collector should not be nil", func() {
			So(collector, ShouldNotBeNil)
		})
		Convey("So Gitstats collector should be of Gitstats type", func() {
			So(collector, ShouldHaveSameTypeAs, &Gitstats{})
		})
		Convey("collector.GetConfigPolicy() should return a config policy", func() {
			configPolicy, _ := collector.GetConfigPolicy()
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})
			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})
			Convey("So config policy namespace should be /raintank/Gitstats", func() {
				conf := configPolicy.Get([]string{"raintank", "apps", "gitstats"})
				So(conf, ShouldNotBeNil)
				So(conf.HasRules(), ShouldBeTrue)
				tables := conf.RulesAsTable()
				So(len(tables), ShouldEqual, 3)
				for _, rule := range tables {
					So(rule.Name, ShouldBeIn, "access_token", "user", "repo")
					switch rule.Name {
					case "access_token":
						So(rule.Required, ShouldBeTrue)
						So(rule.Type, ShouldEqual, "string")
					case "user":
						So(rule.Required, ShouldBeFalse)
						So(rule.Type, ShouldEqual, "string")
					case "repo":
						So(rule.Required, ShouldBeFalse)
						So(rule.Type, ShouldEqual, "string")
					}
				}
			})
		})
	})
}

var open = "open"
var closed = "closed"
var label1 = "label1"
var label2 = "label2"
var label3 = "label3"

var fakeIssues = []*github.Issue{
	{State: &open, Labels: []github.Label{{Name: &label1}}},
	{State: &open, Labels: []github.Label{{Name: &label1}}},
	{State: &open, Labels: []github.Label{{Name: &label2}}},
	{State: &closed, Labels: []github.Label{{Name: &label1}}},
	{State: &closed, Labels: []github.Label{{Name: &label2}}},
	{State: &closed, Labels: []github.Label{{Name: &label2}}},
	{State: &open, Labels: []github.Label{}},
	{State: &closed, Labels: []github.Label{}},
}
var cfg = setupCfg("grafana", "grafana")

type expected struct {
	Key   string
	Value int
}

var issuesByLabelCases = []struct {
	issues   []*github.Issue
	metrics  []plugin.MetricType
	expected []expected
}{
	{
		issues:  fakeIssues,
		metrics: getIssuesByLabelTypes(cfg),
		expected: []expected{
			{"raintank.apps.gitstats.repo.grafana.grafana.issuesbylabel.label1.open.count", 2},
			{"raintank.apps.gitstats.repo.grafana.grafana.issuesbylabel.label1.closed.count", 1},
			{"raintank.apps.gitstats.repo.grafana.grafana.issuesbylabel.label2.open.count", 1},
			{"raintank.apps.gitstats.repo.grafana.grafana.issuesbylabel.label2.closed.count", 2},
			{"raintank.apps.gitstats.repo.grafana.grafana.issuesbylabel.label3.open.count", 0},
			{"raintank.apps.gitstats.repo.grafana.grafana.issuesbylabel.label3.closed.count", 0},
			{"raintank.apps.gitstats.repo.grafana.grafana.issuesbylabel.NoLabel.open.count", 1},
			{"raintank.apps.gitstats.repo.grafana.grafana.issuesbylabel.NoLabel.closed.count", 1},
		},
	},
}

func TestIssueMetrics(t *testing.T) {
	Convey("Create Issue metrics", t, func() {
		for _, c := range issuesByLabelCases {
			metrics, err := collectIssueMetrics(
				c.metrics[0],
				time.Now(),
				"grafana",
				"grafana",
				[]*github.Label{{Name: &label1}, {Name: &label2}, {Name: &label3}}, c.issues)
			So(err, ShouldBeNil)

			So(metrics, ShouldNotBeNil)
			So(len(metrics), ShouldEqual, len(c.expected))
			for _, expected := range c.expected {
				for i := 0; i < len(metrics); i++ {
					if format(&metrics[i]) == expected.Key {
						So(metrics[i].Data(), ShouldEqual, expected.Value)
						t.Log(metrics[i].Namespace().String(), metrics[i].Data())
					}
				}
			}
		}
	})
}

func format(m *plugin.MetricType) string {
	return strings.Join(m.Namespace().Strings(), ".")
}

func SkipTestGetAllIssues(t *testing.T) {
	if os.Getenv("GITSTATS_ACCESS_TOKEN") == "" {
		return
	}

	Convey("Get All Issues", t, func() {
		client := NewClient(os.Getenv("GITSTATS_ACCESS_TOKEN"))
		issues, err := client.GetAllIssues(context.Background(), "grafana", "grafana")

		So(err, ShouldBeNil)
		So(len(issues), ShouldBeGreaterThan, 7000)
		So(issues[0], ShouldEqual, "test")

		labels, err := client.GetAllLabels(context.Background(), "grafana", "grafana")
		So(err, ShouldBeNil)
		So(len(labels), ShouldBeGreaterThan, 10)
	})

}

func TestGitstatsCollectMetrics(t *testing.T) {
	if os.Getenv("GITSTATS_ACCESS_TOKEN") == "" {
		return
	}
	cfg := setupCfg("grafana", "grafana")

	Convey("Ping collector", t, func() {
		p := &Gitstats{}
		mt, err := p.GetMetricTypes(cfg)
		if err != nil {
			t.Fatal("failed to get metricTypes", err)
		}
		So(len(mt), ShouldBeGreaterThan, 0)
		for _, m := range mt {
			t.Log(m.Namespace().String())
		}
		Convey("collect metrics", func() {
			mts := []plugin.MetricType{
				{
					Namespace_: core.NewNamespace(
						"raintank", "apps", "gitstats", "repo", "*", "*", "issuesbylabel", "*", "*", "count"),
					Config_: cfg.ConfigDataNode,
				},
			}
			metrics, err := p.CollectMetrics(mts)
			So(err, ShouldBeNil)
			So(metrics, ShouldNotBeNil)
			// So(len(metrics), ShouldEqual, 1)
			So(metrics[0].Namespace()[0].Value, ShouldEqual, "raintank")
			So(metrics[0].Namespace()[1].Value, ShouldEqual, "apps")
			So(metrics[0].Namespace()[2].Value, ShouldEqual, "gitstats")
			for _, m := range metrics {
				So(m.Namespace()[3].Value, ShouldEqual, "repo")
				So(m.Namespace()[4].Value, ShouldEqual, "grafana")
				t.Log(m.Namespace().String(), m.Data())
			}
		})
	})
}

func setupCfg(user, repo string) plugin.ConfigType {
	node := cdata.NewNode()
	node.AddItem("access_token", ctypes.ConfigValueStr{Value: os.Getenv("GITSTATS_ACCESS_TOKEN")})
	node.AddItem("user", ctypes.ConfigValueStr{Value: user})
	node.AddItem("repo", ctypes.ConfigValueStr{Value: repo})
	return plugin.ConfigType{ConfigDataNode: node}
}
