package main

import (
	"flag"
	"os"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	goscaleio "github.com/codedellemc/goscaleio"

	"github.com/dvonthenen/goprojects/scaleio-test/config"
)

// ----------------------- func init() ------------------------- //

func init() {
	log.SetOutput(os.Stdout)
	log.Infoln("Initializing the ScaleIO Scheduler...")
}

func main() {
	cfg := config.NewConfig()
	fs := flag.NewFlagSet("scheduler", flag.ExitOnError)
	cfg.AddFlags(fs)
	fs.Parse(os.Args[1:])

	if len(cfg.SdsList) == 0 {
		log.Fatalln("SDS List is empty")
	}

	log.Infoln("IP:", cfg.GatewayIP)
	log.Infoln("Version:", cfg.Version)
	log.Infoln("Username:", cfg.Username)
	log.Infoln("Password:", cfg.Password)
	log.Infoln("SDSList:", cfg.SdsList)

	endpoint := "https://" + cfg.GatewayIP + "/api"
	log.Infoln("Endpoint:", endpoint)

	client, err := goscaleio.NewClientWithArgs(endpoint, cfg.Version, true, false)
	if err != nil {
		log.Fatalln("NewClientWithArgs Error:", err)
	}

	_, err = client.Authenticate(&goscaleio.ConfigConnect{
		Endpoint: endpoint,
		Username: cfg.Username,
		Password: cfg.Password,
		//Version:  cfg.Version,
	})
	if err != nil {
		log.Fatalln("Authenticate Error:", err)
	}
	log.Infoln("Successfuly logged in to ScaleIO Gateway at", client.SIOEndpoint.String())

	system, err := client.FindSystem("", "scaleio", "")
	if err != nil {
		log.Fatalln("FindSystem Error:", err)
	}
	log.Infoln("Found system \"scaleio\"")

	tmpPd, err := system.FindProtectionDomain("", "default", "")
	if err != nil {
		//create it!
		pdID, err := system.CreateProtectionDomain("default")
		if err != nil {
			log.Fatalln("CreateProtectionDomain Error:", err)
		}
		tmpPd, err = system.FindProtectionDomain("", "default", "")
		if err != nil {
			log.Fatalln("FindProtectionDomain Error:", err)
		}
		if pdID != tmpPd.ID {
			log.Fatalln("Bad PD:", pdID, "!=", tmpPd.ID)
		}
		log.Infoln("PD Found:", pdID, "=", tmpPd.ID)
	} else {
		log.Infoln("PD Found:", tmpPd.ID)
	}

	pd := goscaleio.NewProtectionDomainEx(client, tmpPd)

	tmpSp, err := pd.FindStoragePool("", "default", "")
	if err != nil {
		//Create it!
		spID, err := pd.CreateStoragePool("default")
		if err != nil {
			log.Fatalln("CreateStoragePool Error:", err)
		}
		tmpSp, err = pd.FindStoragePool("", "default", "")
		if err != nil {
			log.Fatalln("FindStoragePool Error:", err)
		}
		if spID != tmpSp.ID {
			log.Fatalln("Bad SP:", spID, "!=", tmpSp.ID)
		}
		log.Infoln("SP Found:", spID, "=", tmpSp.ID)
	} else {
		log.Infoln("SP Found:", tmpSp.ID)
	}

	sp := goscaleio.NewStoragePoolEx(client, tmpSp)

	sdsIPs := strings.Split(cfg.SdsList, ",")
	for i := 0; i < len(sdsIPs); i++ {
		sdsIDstr := "sds" + strconv.Itoa(i+1)
		tmpSds, err := pd.FindSds("name", sdsIDstr)
		if err == nil && len(tmpSds.ID) > 0 {
			log.Fatalln("FindSds Found:", sdsIDstr)
			continue
		}

		//Create it!
		log.Infoln("Creating", sdsIDstr, "for", sdsIPs[i])

		sdsID, err := pd.CreateSds(sdsIDstr, []string{sdsIPs[i]})
		if err != nil {
			log.Fatalln("CreateSds Error:", err)
		}
		tmpSds, err = pd.FindSds("Name", sdsIDstr)
		if err != nil {
			log.Fatalln("FindSds Error:", err)
		}
		if sdsID != tmpSds.ID {
			log.Fatalln("Bad SP:", sdsID, "!=", tmpSds.ID)
		}
		log.Infoln("SDS Found:", sdsID, "=", tmpSds.ID)

		sds := goscaleio.NewSdsEx(client, tmpSds)

		//Add device
		devID, err := sp.AttachDevice("/dev/xvdf", sds.Sds.ID)
		if err != nil {
			log.Fatalln("AttachDevice Error:", err)
		}
		log.Fatalln("DEV Added:", devID)
	}

}
