package cli

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/coline-carle/battleaxe/battle"
)

// Version of the app
const Version = "0.0.10"

// OauthTokenURL blizard oauth2 url to fetch token
const OauthTokenURL = "https://us.battle.net/oauth/token"

var logger *log.Logger
var errClientIDEmpty = errors.New("client id can't be empty")
var errClientSecretEmpty = errors.New("client secret can't be empty")

// AppName : determine app name bbased on Game value
func AppName(game battle.Game) string {
	switch game {
	case battle.WoW:
		return "wowaxe"
	case battle.D3:
		return "daxe"
	case battle.SC2:
		return "scaxe"

	default:
		return "battleaxe"
	}
}

type app struct {
	flags        *appFlags
	clientID     string
	clientSecret string
	inURL        string
	game         battle.Game
}

func init() {
	logger = log.New(os.Stderr, "", 0)
}

func buildQueryMap(f *appFlags) map[string]string {
	queryMap := make(map[string]string)

	if f.locale != "" {
		queryMap["locale"] = f.locale
	}

	if f.fields != "" {
		queryMap["fields"] = f.fields
	}

	return queryMap
}

func buildURL(url string, game battle.Game, flags *appFlags) (string, error) {
	queryMap := buildQueryMap(flags)

	return battle.ParseURL(url, queryMap, game)
}

func getCredentials(flags *appFlags) (clientID string, clientSecret string, err error) {
	if flags.clientID != "" {
		clientID = flags.clientID
	} else {
		clientID = os.Getenv("BLIZZARD_CLIENT_ID")
	}

	if clientID == "" {
		return "", "", errClientIDEmpty
	}

	if flags.clientSecret != "" {
		clientSecret = flags.clientSecret
	} else {
		clientSecret = os.Getenv("BLIZZARD_CLIENT_SECRET")
	}

	if clientSecret == "" {
		return "", "", errClientSecretEmpty
	}

	return clientID, clientSecret, nil
}

func buildClient(clientID string, clientSecret string) *http.Client {
	blizzOauth := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     OauthTokenURL,
	}

	return blizzOauth.Client(oauth2.NoContext)
}

func doRun(game battle.Game, args []string) (printHelp bool, err error) {
	flags, url, err := parseCommand(args[1:])
	if err != nil {
		return true, err
	}

	if err != nil {
		return true, err
	}

	url, err = buildURL(url, game, flags)
	if err != nil {
		return true, err
	}

	if flags.dry {
		fmt.Println(url)
		return true, err
	}

	clientID, clientSecret, err := getCredentials(flags)
	if err != nil {
		return true, err
	}

	client := buildClient(clientID, clientSecret)
	resp, err := client.Get(url)

	if flags.head {
		PrintHeader(resp)
		return false, err
	}

	if flags.help {
		_ = PrintHelp()
		return false, nil
	}

	err = PrintBody(resp, flags.human)

	return false, err
}

// Run the app
func Run(game battle.Game, args []string) {
	printHelp, err := doRun(game, args)
	if err != nil {
		logger.Println(err)
	}
	if printHelp {
		PrintHelp()
	}
	os.Exit(0)
}
