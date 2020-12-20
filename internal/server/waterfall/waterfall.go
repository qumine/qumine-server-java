package waterfall

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/qumine/qumine-server-java/internal/server/common"
	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	serverVersion      = "latest"
	serverForceUpdate  = false
	serverWaterfallAPI = "https://papermc.io/api/v2/projects/waterfall/"
)

// Server is the struct for waterfall servers.
type Server struct {
	serverVersion     string
	serverForceUpdate bool
	waterfallAPI      string
}

// NewWaterfallServer creates a new waterfall server.
func NewWaterfallServer() *Server {
	return &Server{
		serverVersion:     utils.GetEnvString("SERVER_VERSION", serverVersion),
		serverForceUpdate: utils.GetEnvBool("SERVER_FORCE_UPDATE", serverForceUpdate),
		waterfallAPI:      utils.GetEnvString("SERVER_WATERFALL_API", serverWaterfallAPI),
	}
}

// Configure configures the server.
func (s *Server) Configure() error {
	logrus.WithFields(logrus.Fields{
		"type":        "WATERFALL",
		"version":     s.serverVersion,
		"forceUpdate": s.serverForceUpdate,
		"yatopiaApi":  s.waterfallAPI,
	}).Info("server configuring")

	logrus.Debug("server configured")
	return nil
}

// Update updates the resource, if supported uses cache.
func (s *Server) Update() error {
	logrus.WithFields(logrus.Fields{
		"type":         "WATERFALL",
		"version":      s.serverVersion,
		"forceUpdate":  s.serverForceUpdate,
		"waterfallAPI": s.waterfallAPI,
	}).Info("checking for server updates")

	version := ""
	if match, _ := regexp.MatchString("\\d*\\.\\d*\\.\\d", s.serverVersion); match {
		version = s.serverVersion
	} else if match, _ := regexp.MatchString("\\d*\\.\\d*", s.serverVersion); match {
		versionGroupDetailsDownloadURL := s.waterfallAPI + "version_group/" + s.serverVersion
		versionGroupDetails, err := getVersionGroupDetails(versionGroupDetailsDownloadURL)
		if err != nil {
			return err
		}
		version = versionGroupDetails.Versions[len(versionGroupDetails.Versions)-1]
	} else if s.serverVersion == "latest" {
		// TODO: Implement latest version resolver
	} else {
		return errors.New("Unsupported version")
	}

	versionDetailsDownloadURL := s.waterfallAPI + "versions/" + version
	versionDetails, err := getVersionDetails(versionDetailsDownloadURL)
	if err != nil {
		return err
	}

	buildDetailsURL := versionDetailsDownloadURL + "/builds/" + strconv.Itoa(versionDetails.Builds[len(versionDetails.Builds)-1])
	buildDetails, err := getBuildDetails(versionDetailsDownloadURL + "/builds/" + strconv.Itoa(versionDetails.Builds[len(versionDetails.Builds)-1]))
	if err != nil {
		return err
	}

	if common.CompareHash(s.serverForceUpdate, buildDetails.Downloads.Application.Sha256) {
		logrus.Info("updated server")
		return nil
	}

	if err := common.DownloadServerJar(buildDetailsURL + "/downloads/" + buildDetails.Downloads.Application.Name); err != nil {
		return err
	}

	if err := common.SaveHash(buildDetails.Downloads.Application.Sha256); err != nil {
		return err
	}

	logrus.Info("updated server")
	return nil
}

// StartupCommand retuns the command and arguments used to startup the server.
func (s *Server) StartupCommand() (string, []string) {
	return "java", []string{"-jar", "server.jar", "nogui"}
}
