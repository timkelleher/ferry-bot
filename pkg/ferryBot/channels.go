package ferryBot

import (
	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/state"
	"github.com/sirupsen/logrus"
	"github.com/timkelleher/ferry-bot/pkg/wsdot"
)

// todo: env this
var guildID = discord.GuildID(858197866262822932)

var excludedCategories = []string{"Text Channels"}

type ChannelManager struct {
	logger *logrus.Logger
	s      *state.State
}

func NewChannelManager(s *state.State, logger *logrus.Logger) ChannelManager {
	return ChannelManager{s: s, logger: logger}
}

// Compare wsdot.Terminals to category channels that exist in Discord.
// If a category for a specific Terminal does not exist, create it.
func (cm ChannelManager) SyncCategories(terminals []wsdot.Terminal) error {
	channels, err := cm.s.Channels(guildID)
	if err != nil {
		return err
	}

	unmappedTerminals := make([]wsdot.Terminal, 0)
	categories := getCategories(channels)
	cm.logger.WithField("count", len(channels)).Debug("Category Count")

	for _, terminal := range terminals {
		mapped := false
		for _, name := range categories {
			if name == terminal.Description {
				mapped = true
				break
			}
		}
		if !mapped {
			unmappedTerminals = append(unmappedTerminals, terminal)
		} else {
			cm.logger.WithFields(logrus.Fields{
				"terminal": terminal.Description,
			}).Info("Category Exists")
		}
	}

	for _, terminal := range unmappedTerminals {
		cm.logger.WithField("category", terminal.Description).Info("Creating Category")
		data := api.CreateChannelData{
			Name: terminal.Description,
			Type: discord.GuildCategory,
		}
		_, err := cm.s.CreateChannel(guildID, data)
		if err != nil {
			cm.logger.WithFields(logrus.Fields{
				"category": terminal.Description,
				"error":    err.Error(),
			}).Error("Error Creating Category")
		}
	}

	return nil
}

func (cm ChannelManager) SyncChannels() error {
	// todo
	return nil
}

func (cm ChannelManager) GetChannels() (map[discord.ChannelID][]discord.Channel, error) {
	sortedChannels := make(map[discord.ChannelID][]discord.Channel)

	channels, err := cm.s.Channels(guildID)
	if err != nil {
		return sortedChannels, err
	}

	categories := getCategories(channels)
	for _, channel := range channels {
		if channel.Type == discord.GuildCategory || !validParentChannel(channel.CategoryID, categories) {
			continue
		}
		key := channel.CategoryID
		sortedChannels[key] = append(sortedChannels[key], channel)
	}
	return sortedChannels, nil
}

func getCategories(channels []discord.Channel) map[discord.ChannelID]string {
	categories := make(map[discord.ChannelID]string)
	for _, channel := range channels {
		if channel.Type == discord.GuildCategory && !excludeChannel(channel.Name, excludedCategories) {
			categories[channel.ID] = channel.Name
		}
	}
	return categories
}

func excludeChannel(name string, excludedChannels []string) bool {
	for _, excluded := range excludedChannels {
		if excluded == name {
			return true
		}
	}
	return false
}

func validParentChannel(id discord.ChannelID, categories map[discord.ChannelID]string) bool {
	for categoryID, _ := range categories {
		if categoryID == id {
			return true
		}
	}
	return false
}

// Debug: this is useful when first building categories
func deleteCategories() {
	/*
		cm.logger.WithFields(logrus.Fields{
			"count": len(channels),
		}).Info("Channel Count")
		deleteCats := getCategories(channels)
		for id, _ := range deleteCats {
			cm.s.DeleteChannel(id)
		}
		channels, err = cm.s.Channels(guildID)
		if err != nil {
			return err
		}
		cm.logger.WithFields(logrus.Fields{
			"count": len(channels),
		}).Info("Channel Count")
	*/
}
