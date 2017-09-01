// Read POM elements and show on stdout
// Exit codes:
//      1: missing arguments
//      2: unsupported ModelVersion in POM

package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

const (
	supportedMavenVersion = "4.0.0"
)

var (
	showArtifact = flag.Bool("artifact", false, "Show artifact")
	showGroup    = flag.Bool("group", false, "Show group")
	showVersion  = flag.Bool("version", true, "Show version")
)

// Project is a subset of Maven POM 4.0.0
type Project struct {
	XMLName xml.Name `xml:"project"`

	ModelVersion string `xml:"modelVersion"`

	Group    string `xml:"groupId"`
	Artifact string `xml:"artifactId"`
	Version  string `xml:"version"`
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage: %s [options] <path to pom.xml> ...\n",
			path.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "Show POM elements on stdout\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	n := flag.NArg()
	if n == 0 {
		flag.Usage()
		os.Exit(1)
	}

	for _, filename := range flag.Args() {
		buf, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		var p Project
		if err := xml.Unmarshal(buf, &p); err != nil {
			log.Fatal(err)
		}
		if p.ModelVersion != supportedMavenVersion {
			s := "File %s contains model version %s, " +
				"the only supported model version is %s\n"
			fmt.Fprintf(os.Stderr, s, filename, p.ModelVersion,
				supportedMavenVersion)
			os.Exit(2)
		}
		var gav []string
		if *showGroup {
			gav = append(gav, p.Group)
		}
		if *showArtifact {
			gav = append(gav, p.Artifact)
		}
		if *showVersion {
			gav = append(gav, p.Version)
		}
		fmt.Println(strings.Join(gav, ":"))
	}
}
