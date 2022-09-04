package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	user2 "os/user"
	"path"
	"strings"
)
type config struct {
	Sender string `json:"sender"`
	Receiver string `json:"receiver"`
	StorePath string `json:"storepath"`
	Password string `json:"password"`
	Server string `json:"server"`
	Port int `json:"port"`
}
var instance *config
func isGmail(mail string) bool {
	//use regex here
	return strings.Contains(mail, "@gmail.com")
}
func DefaultConfigPath() (string, error){
	user, err:=user2.Current()
	if err!=nil{
		log.Fatal("Couldn't get current user ", err)
	}
	return path.Join(user.HomeDir, "KindleConfig.json"), nil
}
func exists(filename string) bool{
	if _,err:=os.Stat(filename);err!=nil{
		log.Println(err)
		return false
	}
	return true
}
func NewConfig()*config{
	config:=config{}
	config.Server="smtp.gmail.com"
	config.Port=465
	return &config
}
func CreateConfig() *config{
	configuration:=NewConfig()
	reader:=bufio.NewReader(os.Stdin)
	fmt.Printf("Email of your kindle and press enter (eg. purple_terminal@kindle.com) : ")
	fmt.Scan(&configuration.Receiver)
	fmt.Printf("Email that'll be used to send documents to kindle (eg. yourname@gmail.com) : ")
	fmt.Scan(&configuration.Sender)

	if isGmail(configuration.Sender)==false{
		fmt.Println("Sender email is different then Gmail, " +
			"can you help with SMTP server address and SMTP port for your email provider\n" +
			"Just search SMTP settings for <your email domain>.com on internet \n" +
			"-----------------------------------------")

		fmt.Printf("Enter SMTP Server Address (eg. smtp.gmail.com) : ")
		fmt.Scan(&configuration.Server)
		fmt.Printf("Enter SMTP port (usually 587 or 485) : ")
		fmt.Scan(&configuration.Port)
	}

	fmt.Printf("Enter password for Sender %s (password remains encrypted in your machine) : ",configuration.Receiver)
	fmt.Scan(&configuration.Password)


	fmt.Printf("File path to store all the documents on your computer (empty is ok) :")
	configuration.StorePath, _ =reader.ReadString('\n')

	configuration.StorePath=strings.Trim(configuration.StorePath, "\n")

	encryptedPass, err := Encrypt(configuration.Sender, configuration.Password)
	if err!=nil{
		log.Println("Error encrypting password: ", err)
		os.Exit(1)
	}
	configuration.Password=encryptedPass

	if err!=nil{
		log.Println(err)
		os.Exit(1)
	}
	return configuration
}
func handleCreation(filename string) error {
	fmt.Println("Configuration file doesn't exist\n Answer next few questions to create config file\n")
	configuration:=CreateConfig()
	err := Save(*configuration, filename)
	if err!=nil{
		log.Println("Error while writing config to ",filename, err)
		return err
	}
	fmt.Printf("Config created successfully and stored at %s, you can directly edit it later on ", filename)
	return nil
}
func Load(filename string) (config, error) {
	if !exists(filename){
		err:=handleCreation(filename)
		if err!=nil{
			return config{}, err
		}
	}
	data, err:=ioutil.ReadFile(filename)
	if err!=nil{
		log.Println("Error reading config ", err)
		return config{}, err
	}
	var c config
	err = json.Unmarshal(data, &c)
	if err!=nil{
		log.Println("Error converting config to json ", err)
		return config{}, err
	}
	decryptedPass, err:=Decrypt(c.Sender, c.Password)
	if err!=nil{
		log.Println("Error decrypting password : ", err)
		os.Exit(1)
	}
	c.Password=decryptedPass
	log.Println("loaded configuration")
	InitializeConfig(&c)
	return c, nil
}

func Save(c config, filename string) error {
	data, err:=json.MarshalIndent(c, "", "	")
	if err!=nil{
		log.Println("Error parsing configuration for writing")
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func InitializeConfig(c *config){
	if instance==nil{
		instance=c
		log.Println("Initialized configuration instance")
	}
}

func GetInstance() *config{
	return instance
}
