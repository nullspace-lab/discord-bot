package handlers

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nullspace-lab/discord-bot/config"
)

func StartScheduler(s *discordgo.Session, cfg *config.Config) {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			waitFor := time.Until(next)

			log.Printf("next member promotion scheduled for: %s\n", next.Format("2006-01-02 15:04:05"))
			time.Sleep(waitFor)

			promoteNewMembers(s, cfg)
		}
	}()
}

func promoteNewMembers(s *discordgo.Session, cfg *config.Config) {
	log.Println("starting new member promotion...")

	var promoted int
	after := ""

	for {
		members, err := s.GuildMembers(cfg.GuildID, after, 1000)
		if err != nil {
			log.Println("failed to fetch members:", err)
			return
		}

		for _, member := range members {
			hasNewRole := false
			for _, roleID := range member.Roles {
				if roleID == cfg.NewMemberRoleID {
					hasNewRole = true
					break
				}
			}

			if !hasNewRole {
				continue
			}

			err = s.GuildMemberRoleAdd(cfg.GuildID, member.User.ID, cfg.OficialMemberRoleID)
			if err != nil {
				log.Printf("failed to promote %s: %v\n", member.User.Username, err)
				continue
			}

			err = s.GuildMemberRoleRemove(cfg.GuildID, member.User.ID, cfg.NewMemberRoleID)
			if err != nil {
				log.Printf("failed to remove new member role from %s: %v\n", member.User.Username, err)
				continue
			}

			promoted++
			log.Printf("member %s promoted\n", member.User.Username)
		}

		if len(members) < 1000 {
			break
		}

		after = members[len(members)-1].User.ID
	}

	log.Printf("promotion complete: %d members promoted\n", promoted)
}
