package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	AppVersion   = "0.0.1"
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

var (
	argConfigPath = flag.String("config", "config.json", "設定ファイルを指定.")
	argVersion    = flag.Bool("version", false, "バージョンを出力.")
	argEnable     = flag.Bool("enable", false, "通知を有効化.")
	argDisable    = flag.Bool("disable", false, "通知を無効化.")
	argCheck      = flag.Bool("check", false, "通知を状態を確認.")
	argHosts      = flag.String("hosts", "", "通知を操作するホスト名を指定.")
	argService    = flag.String("service", "", "通知を操作するサービス名を指定.")
)

type NagillaConfig struct {
	Nagios struct {
		URL string `json:"url"`
	} `json:"nagios"`
	Targets struct {
		Host string `json:"host"`
	} `json:"targets"`
}

func loadConfig() (*NagillaConfig, error) {
	f, err := os.Open(*argConfigPath)
	if err != nil {
		log.Fatalf("Config File Load Error:", err)
	}
	defer f.Close()
	var cfg NagillaConfig
	err = json.NewDecoder(f).Decode(&cfg)

	return &cfg, err
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func opeRequest(nagiosUrl string, typ string, targetHost string, c *NagillaConfig) string {
	var url string
	if typ == "" {
		url = nagiosUrl + "/nagios/cgi-bin//extinfo.cgi?type=1&host=" + targetHost
	} else {
		url = nagiosUrl + "/nagios/cgi-bin//cmd.cgi?cmd_mod=2&cmd_typ=" + typ + "&host=" + targetHost
	}

	// リクエストを生成する
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var nUser string
	var nPass string

	if os.Getenv("NAGIOS_USER") == "" {
		fmt.Println("Nagios のログインユーザー名を指定して下さい.")
		os.Exit(1)
	} else {
		nUser = os.Getenv("NAGIOS_USER")
	}

	if os.Getenv("NAGIOS_PASS") == "" {
		fmt.Println("Nagios のログインパスワードを指定して下さい.")
		os.Exit(1)
	} else {
		nPass = os.Getenv("NAGIOS_PASS")
	}

	req.Header.Add("Authorization", "Basic "+basicAuth(nUser, nPass))

	// Nagios で監視しているホストの通知設定を変更
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Unable to connect nagios server.\n")
	} else if res.StatusCode != 200 {
		log.Fatalf("Unable to get url : http status %d\n", res.StatusCode)
	}
	defer res.Body.Close()

	if typ == "" {
		ParseCheckHostStatus(targetHost, res.Body)
		os.Exit(0)
	}

	return ParseOpeResult(res.Body)
}

func opeNagios(nagiosUrl string, typ string, targetHost string, c *NagillaConfig) {
	switch typ {
	case "enable":
		fmt.Printf(WarningColor, targetHost+"\n")
		fmt.Print("上記のホストの通知を有効しますか?(y/n): ")
	case "disable":
		fmt.Printf(WarningColor, targetHost+"\n")
		fmt.Print("上記のホストの通知を無効しますか?(y/n): ")
	case "check":
		opeRequest(nagiosUrl, "", targetHost, c)
	default:
		opeRequest(nagiosUrl, "", targetHost, c)
	}

	// Nagios コマンドの実行結果を格納する
	var result string
	results := []string{}

	var stdin string
	fmt.Scan(&stdin)
	switch stdin {
	case "y", "Y":
		switch typ {
		case "enable":
			fmt.Println("通知を有効にします.")
			// 24, 28
			typs := [2]string{"24", "28"}
			for _, t := range typs {
				result = opeRequest(nagiosUrl, t, targetHost, c)
				results = append(results, result)
			}
		case "disable":
			fmt.Println("通知を無効にします.")
			// 25, 29
			typs := [2]string{"25", "29"}
			for _, t := range typs {
				result = opeRequest(nagiosUrl, t, targetHost, c)
				results = append(results, result)
			}
		}
	case "n", "N":
		fmt.Println("処理を停止します.")
		os.Exit(0)
	default:
		fmt.Println("処理を停止します.")
		os.Exit(0)
	}

	if strings.Join(results, ",") == "ok,ok" {
		switch typ {
		case "enable":
			fmt.Println("通知を有効にしました.")
		case "disable":
			fmt.Println("通知を無効にしました.")
		}
	} else {
		switch typ {
		case "enable":
			fmt.Println("通知有効処理に失敗しました.")
		case "disable":
			fmt.Println("通知無効処理に失敗しました.")
		}
	}
}

func main() {
	flag.Parse()

	if *argVersion {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	// 各種設定ファイルの読み込み (config.json から読み込む)
	c, _ := loadConfig()
	nagiosUrl := c.Nagios.URL

	// 操作するホスト名を変数にセットする (コマンドライン引数か, config.json から読み込む)
	var targetHosts string
	if *argHosts != "" {
		targetHosts = *argHosts
	} else {
		targetHosts = c.Targets.Host
	}

	// 操作するホストが設定されているか確認する
	if targetHosts == "" {
		fmt.Println("操作するホスト名を設定して下さい.")
		os.Exit(1)
	}

	if *argEnable {
		opeNagios(nagiosUrl, "enable", targetHosts, c)
	} else if *argDisable {
		opeNagios(nagiosUrl, "disable", targetHosts, c)
	} else if *argCheck {
		opeNagios(nagiosUrl, "check", targetHosts, c)
	} else {
		fmt.Println("操作タイプ (enable, disable, check) を設定して下さい.")
		os.Exit(1)
	}
}
