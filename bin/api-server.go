package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/kirsle/configdir"
	"github.com/omecodes/ditt"
	"github.com/omecodes/ditt/info"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const (
	cookiesFile       = "cookies-key"
	cookieStoreDir    = "cookies-store"
	adminPasswordFile = "admin-auth"
)

var (
	port        int
	dataDirname string
	databaseURI string
	cmd         *cobra.Command
)

func init() {
	programName := filepath.Base(os.Args[0])

	cmd = &cobra.Command{
		Use: programName,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	versionCommand := &cobra.Command{
		Use:   "version",
		Short: "Shows the version info",
		Run: func(cmd *cobra.Command, args []string) {
			showVersion()
		},
	}

	startCommand := &cobra.Command{
		Use:   "start",
		Short: "Starts the server",
		Run: func(cmd *cobra.Command, args []string) {
			startServer()
		},
	}

	flags := startCommand.PersistentFlags()
	flags.IntVar(&port, "port", 80, "The HTTP server port")
	flags.StringVar(&databaseURI, "db-uri", "localhost", "The database URI")
	flags.StringVar(&dataDirname, "data-dir", "", "Directory path in where file data are saved")

	cmd.AddCommand(versionCommand)
	cmd.AddCommand(startCommand)
}

func showVersion() {
	fmt.Println()
	fmt.Println("       Version:", info.ApplicationName)
	fmt.Println("          Name:", info.Version)
	fmt.Println("Build Revision:", info.BuildRevision)
	fmt.Println("   Build Stamp:", info.BuildStamp)
	fmt.Println()
}

func startServer() {
	configDir := configdir.LocalConfig(info.ApplicationName)
	err := os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("CONFIG DIR: ", configDir)

	setupDataDir(configDir)
	setupCookies(configDir)
	setAdminAuthentication(configDir)
	setupMongoDB()

	err = ditt.Serve(&ditt.Config{Port: port})
	if err != nil {
		log.Fatalln(err)
	}
}

func setupDataDir(configDir string) {
	if dataDirname == "" {
		dataDirname = filepath.Join(configDir, "data")
		err := os.MkdirAll(dataDirname, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	}
	ditt.Env.Files = ditt.NewDirFiles(dataDirname)
}

func setupCookies(configDir string) {
	cookieKeyFilename := filepath.Join(configDir, cookiesFile)
	keyData, err := ioutil.ReadFile(cookieKeyFilename)
	if err != nil {
		if os.IsNotExist(err) {
			keyData = make([]byte, 64)
			_, err = rand.Read(keyData)
			if err != nil {
				log.Fatalln(err)
			}

			err = ioutil.WriteFile(cookieKeyFilename, keyData, 0600)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Fatalln(err)
		}
	}
	cookiesStoreDirname := filepath.Join(configDir, cookieStoreDir)
	err = os.MkdirAll(cookiesStoreDirname, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	ditt.Env.CookiesStore = sessions.NewFilesystemStore(cookiesStoreDirname, keyData[:31], keyData[32:])
}

func setupMongoDB() {
	store, err := ditt.NewMongoUserDataStore(databaseURI)
	if err != nil {
		log.Fatalln("Mongo", err)
	}
	ditt.Env.DataStore = store
}

func setAdminAuthentication(configDir string) {
	adminPasswordFilename := filepath.Join(configDir, adminPasswordFile)

	adminPasswordBytes, err := ioutil.ReadFile(adminPasswordFilename)
	if err != nil {
		if os.IsNotExist(err) {
			adminPasswordBytes = generateRandomPassword(16)
			err = ioutil.WriteFile(adminPasswordFilename, adminPasswordBytes, 0600)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Fatalln(err)
		}
	}
	ditt.Env.AdminPassword = string(adminPasswordBytes)
}

func generateRandomPassword(length int) []byte {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "=+*!@#$"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" + digits + specials
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return buf
}

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
