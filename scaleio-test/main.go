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

	client, err := goscaleio.NewClient()
	if err != nil {
		log.Fatalln("NewClient Error:", err)
	}

	endpoint := "http://" + cfg.GatewayIP + "/api"

	_, err = client.Authenticate(&goscaleio.ConfigConnect{
		Endpoint: endpoint,
		Username: cfg.Username,
		Password: cfg.Password,
		Version:  cfg.Version,
	})
	if err != nil {
		log.Fatalln("Authenticate Error:", err)
	}
	log.Infoln("Successfuly logged in to ScaleIO Gateway at", client.SIOEndpoint.String())

	system, err := client.FindSystem("", "scaleio", "")
	if err != nil {
		log.Fatalln("FindSystem Error:", err)
	}

	pdID, err := system.CreateProtectionDomain("pd")
	if err != nil {
		log.Fatalln("CreateProtectionDomain Error:", err)
	}
	tmpPd, err := system.FindProtectionDomain("", "pd", "")
	if err != nil {
		log.Fatalln("FindProtectionDomain Error:", err)
	}
	if pdID != tmpPd.ID {
		log.Fatalln("Bad PD:", pdID, "!=", tmpPd.ID)
	}

	pd := goscaleio.NewProtectionDomainEx(client, tmpPd)

	spID, err := pd.CreateStoragePool("sp")
	if err != nil {
		log.Fatalln("CreateStoragePool Error:", err)
	}
	tmpSp, err := pd.FindStoragePool("", "sp", "")
	if err != nil {
		log.Fatalln("FindStoragePool Error:", err)
	}
	if spID != tmpSp.ID {
		log.Fatalln("Bad SP:", spID, "!=", tmpSp.ID)
	}

	sp := goscaleio.NewStoragePoolEx(client, tmpSp)

	sdsIPs := strings.Split(cfg.SdsList, ",")
	for i := 0; i < len(sdsIPs); i++ {
		sdsIDstr := "sds" + strconv.Itoa(i+1)
		sdsID, err := pd.CreateSds(sdsIDstr, []string{sdsIPs[i]}, pd.ProtectionDomain.ID,
			[]string{"/dev/xvdf"}, []string{sp.StoragePool.ID})
		if err != nil {
			log.Fatalln("CreateSds Error:", err)
		}
		tmpSds, err := pd.FindSds("name", sdsIDstr)
		if err != nil {
			log.Fatalln("FindSds Error:", err)
		}
		if sdsID != tmpSds.ID {
			log.Fatalln("Bad SP:", sdsID, "!=", tmpSds.ID)
		}
	}
}
