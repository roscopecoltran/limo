package apk

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"github.com/agrison/go-tablib"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/jinzhu/configor"
	"github.com/roscopecoltran/sniperkit-limo/plugins/docker/alpine/version"
)

const (
	// BANNER is what is printed for help/info output
	APP_BANNER = `             _          __ _ _
  __ _ _ __ | | __     / _(_) | ___
 / _` + "`" + ` | '_ \| |/ /____| |_| | |/ _ \
| (_| | |_) |   <_____|  _| | |  __/
 \__,_| .__/|_|\_\    |_| |_|_|\___|
      |_|

 Search apk package contents via the command line.
 Version: %s

`
	alpineContentsSearchURI = "https://pkgs.alpinelinux.org/contents"
)

type fileInfo struct {
	path, pkg, branch, repo, arch string
}

var (

	// search args
	arch string
	repo string

	// output
	output string
	result string
	prefixPath string 
	filename string

	// user inputs
	input_file 	string 
	input_path 	string 

	// app features
	debug bool
	vrsn  bool
	save bool

	cfgFormat string = "yaml" 

	// valid inputs
	validOutput = []string{"markdown", "csv", "yaml", "json", "xlsx", "xml", "tsv", "mysql", "postgres", "html", "ascii"}
	validArches = []string{"x86", "x86_64", "armhf"}
	validRepos  = []string{"main", "community", "testing"}

	// config 
	cfg 	*Config

	// logs
	log 	= logrus.New()

)

type Config struct {

	App struct {
		Name string `default:"app name"`
		Debug bool `default:"false"`
		Version string `default:"dev"`
		Config struct {
			PrefixPaths []string{"~/.apk-file/apk-file", "./shared/conf.d/apk-file", "./shared/conf.d/apk-file", "../../shared/conf.d/apk-file"}
			Formats []string{"yaml", "json", "xml"}
			Write bool `default:"true"`
			Print bool `default:"false"`
		}
	}

	Queries []Query 

	Results struct {
		Output string
		Result string
	}

	Valid struct {
		Output []string{"md", "csv", "yaml", "json", "xlsx", "xml", "tsv", "mysql", "postgres", "html", "ascii"} `long:"valid-output" description:"valid output formats"`
		Arches []string{"x86", "x86_64", "armhf"} `long:"valid-archs" description:"valid architectures list"`
		Repos  []string{"main", "community", "testing"} `long:"valid-repos" description:"valid repos list"`
	}

}{}

type Query struct {
	File string 
	Path string
	Branch string 
	Repo string 
	Arch string 
}

func init() {

	// logs
	log.Out = os.Stdout

	// Parse flags
	flag.StringVar(&arch, "arch", "", "arch to search for ("+strings.Join(validArches, ", ")+")")
	flag.StringVar(&repo, "repo", "", "repository to search in ("+strings.Join(validRepos, ", ")+")")
	flag.StringVar(&output, "output", "", "output results with  ("+strings.Join(validOutput, ", ")+") format.")
	flag.StringVar(&prefixPath, "./output", "results", "output results to prefix_path (default: ./output).")
	flag.StringVar(&filename, "filename", "results", "output results to filename: (default: ./results.[FORMAT]).")

	flag.StringVar(&cfgFormat, "format-config", "yaml", "default data format for configuration file.")
	flag.Bool(&cfgPrint, "print-config", false, "print configuration file.")
	flag.Bool(&cfgSave, "save-config", false, "save output results to the output_file.[FORMAT].")

	flag.BoolVar(&save, "save", true, "save output results to the output_file.[FORMAT].")
	flag.BoolVar(&vrsn, "version", false, "print version and exit")
	flag.BoolVar(&vrsn, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&debug, "d", false, "run in debug mode")

	configor.Load(&cfg, "./shared/conf.d/apk-file.yml", "../shared/conf.d/apk-file.yml", "../../shared/conf.d/apk-file.yml", "~/.sniperkit/plugins/apk-file.yml")
	fmt.Printf("config: %#v", Config)

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(APP_BANNER, version.SNK_PLUGIN_APK_FILE_VERSION))
		flag.PrintDefaults()
	}

	flag.Parse()

	if vrsn {
		fmt.Printf("apk-file version %s, build %s", version.SNK_PLUGIN_APK_FILE_VERSION, version.SNK_PLUGIN_APK_FILE_VCS_GIT_COMMIT)
		os.Exit(0)
	}

	// Set log level
	if debug {
		// log.SetLevel(log.DebugLevel)
		log.Level = logrus.DebugLevel
	}

	if arch != "" && !stringInSlice(arch, validArches) {
		log.Fatalf("%s is not a valid arch", arch)
	}

	if repo != "" && !stringInSlice(repo, validRepos) {
		log.Fatalf("%s is not a valid repo", repo)
	}

}

	
func check(err error) {
    if err != nil {
        panic(err)
    }
}

func SearchApkCli() (string, string) {
	if flag.NArg() < 1 {
		log.Fatal("must pass a file to search for.")
	}

	if cfgPrint || cfgSave {
		localConfigFilePath := findLocalConfig()
		if localConfigFilePath != "" {
			log.WithFields(logrus.Fields{
				"action": "localConfigFilePath",
				"localConfigFilePath": localConfigFilePath,
			}).Infof("Local config file path found at: %#s\n", localConfigFilePath)
		} else {
			log.WithFields(logrus.Fields{
				"action": "localConfigFilePath",
			}).Warnf("No Local config file path found")
		}
	}

	if cfgPrint {
		cfg.App.Config.Print = cfgPrint
	}

	if cfgSave {
		cfg.App.Config.Write = cfgPrint
	}

	// handle print-config and write-config before auto detection to prevent
	// auto detected values from being written to the config file
	if err := processConfigOptions(); err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"action": "processConfigOptions",
		}).Fatal("could not process configuration file.")
	}

	input_file, input_path 	= getFileAndPath(flag.Arg(0))	
	return input_file, input_path
}

func SearchAPK(cliMode bool, queryString string) {

	if cliMode != false {
	 	input_file, input_path = SearchApkCli()
	} else {
		if queryString != "" {
			input_file, input_path = getFileAndPath(queryString)	
		}
	}

	if input_file == "" && input_path == "" {
		log.WithFields(logrus.Fields{
			"input_file": input_file,
			"input_path": input_path,
		}).Fatal("must pass a file to search for.")
	}

	query := url.Values{
		"file":   {input_file},
		"path":   {input_path},
		"branch": {""},
		"repo":   {repo},
		"arch":   {arch},
	}

	log.WithFields(logrus.Fields{
		"input_file": input_file,
		"path": path,
		"repo": repo,
		"branch": "",
		"arch": arch,
		"query":  query,
	}).Info("SearchAPK by filename")

	uri := fmt.Sprintf("%s?%s", alpineContentsSearchURI, query.Encode())
	doc, err := goquery.NewDocument(uri)
	if err != nil {
		log.Fatalf("requesting %s failed: %v", uri, err)
	}

	files := getFilesInfo(doc)

	ds 		:= tablib.NewDataset([]string{"file", "package", "branch", "repository", "architecture"})
	
	for _, f := range files {
		// https://github.com/agrison/go-tablib
		ds.AppendValues(f.path, f.pkg, f.branch, f.repo, f.arch)
	}

	switch output {

		case "csv":
			result, _ := ds.CSV()
			if save == true {
				if result.WriteFile(prefixPath+"/"+filename+"."+output, 0644) != nil {
				    fmt.Println(err)
				}
			}
			fmt.Println(result)
		case "tsv":
			result, _ := ds.TSV()
			if save == true {
				if result.WriteFile(prefixPath+"/"+filename+"."+output, 0644) != nil {
				    fmt.Println(err)
				}
			}
			fmt.Println(result)
		case "yaml":
			result, _ := ds.YAML()
			if save == true {
				if result.WriteFile(prefixPath+"/"+filename+"."+output, 0644) != nil {
				    fmt.Println(err)
				}
			}
			fmt.Println(result)
		case "json":
			result, _ := ds.JSON()
			if save == true {
				if result.WriteFile(prefixPath+"/"+filename+"."+output, 0644) != nil {
				    fmt.Println(err)
				}
			}
			fmt.Println(result)
		case "xlsx":
			result, _ := ds.XLSX()
			if save == true {
				if result.WriteFile(prefixPath+"/"+filename+"."+output, 0644) != nil {
				    fmt.Println(err)
				}
			}
			fmt.Println(result)
		case "xml":
			result, _ := ds.XML()
			if save == true {
				if result.WriteFile(prefixPath+"/"+filename+"."+output, 0644) != nil {
				    fmt.Println(err)
				}
			}
			fmt.Println(result)
		case "html":
			result, _ := ds.XLSX()
			fmt.Println(result)
			if result.WriteFile(prefixPath+"/"+filename+"."+output, 0644) != nil {
			    fmt.Println(err)
			}
			fmt.Println(result)
		case "ascii":
		default:
			ascii := ds.Tabular("grid" /* tablib.TabularGrid */)	
			fmt.Println(ascii)
	}


}

func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
}

func getFilesInfo(d *goquery.Document) []fileInfo {
	files := []fileInfo{}
	d.Find(".table tr:not(:first-child)").Each(func(j int, l *goquery.Selection) {
		f := fileInfo{}
		rows := l.Find("td")
		rows.Each(func(i int, s *goquery.Selection) {
			switch i {
			case 0:
				f.path = s.Text()
			case 1:
				f.pkg = s.Text()
			case 2:
				f.branch = s.Text()
			case 3:
				f.repo = s.Text()
			case 4:
				f.arch = s.Text()
			default:
				log.Warn("Unmapped value for column %d with value %s", i, s.Text())
			}
		})
		files = append(files, f)
	})
	return files
}

func getFileAndPath(arg string) (file string, dir string) {
	file = "*" + path.Base(arg) + "*"
	dir = path.Dir(arg)
	if dir != "" && dir != "." {
		dir = "*" + dir
		file = strings.TrimPrefix(file, "*")
	} else {
		dir = ""
	}
	return file, dir
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// findLocalConfig returns the path to the local config file.
// It searches the current directory and all parent directories for a config file.
// If no config file is found, findLocalConfig returns an empty string.
func findLocalConfig() string {
	curdir, err := os.Getwd()
	if err != nil {
		curdir = "."
	}
	path, err := filepath.Abs(curdir)
	if err != nil || path == "" {
		return ""
	}
	lp := ""
	for cfgPrefixPath := range Options.App.Config.PrefixPaths {
		for cfgFormat := range Options.App.Config.Formats {
			confpath := filepath.Join(path, cfgPrefixPath+"."+cfgFormat)
			if _, err := os.Stat(confpath); err == nil {
				return confpath
			}
			lp = path
			path = filepath.Dir(path)
		}
	}

	return ""
}

func (o *Options) processConfigOptions() error {

	var configFilePath string
	localConfigFilePath := findLocalConfig()
	if localConfigFilePath != "" {
		fmt.Fprintf(os.Stderr, "Local config file path: %s\n", localConfigFilePath)
	} else {
		fmt.Fprintf(os.Stderr, "No local config file found.\n")
	}


	if o.Options.Config.Print {

		// conf, err := json.MarshalIndent(o, "", "    ")
		conf, err := cfg.YAML()
		if err != nil {
			return fmt.Errorf("cannot convert config to JSON: %s", err)
		}
		fmt.Println(conf)
		fmt.Println(string(conf))
		os.Exit(0)
	}

	if o.Options.Config.Write {
		conf, err := json.MarshalIndent(o, "", "    ")
		if err != nil {
			return fmt.Errorf("cannot convert config to JSON: %s", err)
		}
		if err := ioutil.WriteFile(configFilePath, conf, os.ModePerm); err != nil {
			return fmt.Errorf("cannot write config file: %s", err)
		}
		fmt.Printf("Saved config to '%s'.\n", configFilePath)
		os.Exit(0)
	}

	return nil
}