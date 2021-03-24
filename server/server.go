package server

import (
	"context"
	"os"
	"regexp"
	"sync"
	"time"

	server "github.com/balazshorvath/go-srv"
	"github.com/cloudflare/cloudflare-go"
	"github.com/rs/zerolog"

	"CloudflareDDNS/config"
	"CloudflareDDNS/ipify"
)

type cfServer struct {
	server.BasicServer
	logger *zerolog.Logger
	config *config.Config
}

func New(ctx context.Context, group *sync.WaitGroup) server.Server {
	logger := zerolog.New(os.Stderr)
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config.yaml"
	}
	conf, err := config.LoadFile(path)
	if err != nil {
		logger.Fatal().Err(err).Msgf("Failed to load config file")
	}
	return &cfServer{
		BasicServer: server.BasicServer{
			Ctx:   ctx,
			Group: group,
		},
		config: conf,
		logger: &logger,
	}
}

func (c *cfServer) Init() {
}

func (c *cfServer) Start() {
	go func() {
		c.Group.Add(1)
		defer c.Group.Done()
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		previous := ""
		matcher := regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
		for {
			select {
			case <-c.Ctx.Done():
				return
			case <-ticker.C:
				ip, err := ipify.GetIp()
				if err != nil {
					c.logger.Error().Err(err).Msg("Failed to query ip address")
					break
				}
				if !matcher.MatchString(ip) {
					c.logger.Error().Msgf("Not a valid ip address")
					break
				}
				if previous == ip {
					break
				}
				err = c.UpdateRecords(c.Ctx, ip)
				if err != nil {
					break
				}
				previous = ip
			}
		}
	}()
}

func (c *cfServer) Shutdown(ctx context.Context) error {
	return nil
}

func (c *cfServer) UpdateRecords(ctx context.Context, ip string) error {
	api, err := cloudflare.New(c.config.CfKey, c.config.CfEmail)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to create Cloudflare API client")
		return err
	}
	zoneId, err := api.ZoneIDByName(c.config.Zone)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to query zoneId from Cloudflare API")
		return err
	}
	records, err := api.DNSRecords(ctx, zoneId, cloudflare.DNSRecord{})
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to query zone dns records from Cloudflare API")
		return err
	}
	var lastError error = nil
	for _, record := range records {
		if record.Type != "A" {
			continue
		}
		for _, name := range c.config.Names {
			if name != record.Name {
				continue
			}
			record.Content = ip
			err = api.UpdateDNSRecord(ctx, zoneId, record.ID, record)
			if err != nil {
				c.logger.Error().Err(err).Msgf("Failed to update new ip (%s) in dns record %v", ip, record)
				lastError = err
			}
			c.logger.Info().Msgf("[%s] record updated with ip [%s] in zone [%s]", record.Name, ip, zoneId)
		}
	}
	return lastError
}
