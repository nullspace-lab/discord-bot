package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/nullspace-lab/discord-bot/config"
)

func rulesEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "📜 Regras da Comunidade",
		Color: 0x5865F2,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "🧭 Valores da Comunidade",
				Value: "> Respeito mútuo por todos\n> Inclusão e diversidade\n> Comunicação direta, sem enrolação\n> Conhecimento como pilar\n> Evolução técnica, mental e pessoal\n> Alta performance com consistência — sem papo de coach",
			},
			{
				Name:  "📌 Regras Gerais",
				Value: "> Use o grupo certo pra cada tema (Frontend, Backend, Vagas etc)\n> Proibido divulgar serviços, produtos ou canais sem permissão dos admins\n> Não mande mensagem privada pra membros sem consentimento\n> Eventos e parcerias: fale antes com a moderação",
			},
			{
				Name:  "✅ O que esperamos",
				Value: "> Comunicação respeitosa e não-violenta\n> Poste dúvidas completas (não só \"alguém pode ajudar?\")\n> Compartilhe o que pode ajudar outros a crescer\n> Apoie quem tá começando ou precisa de uma força",
			},
			{
				Name:  "🚫 O que não aceitamos",
				Value: "> Assédio ou preconceito\n> Ironias ofensivas, ameaças ou linguagem agressiva\n> Conteúdo ilegal, pirataria ou pornografia\n> Venda irregular de licenças\n> Compra coletiva de softwares pagos\n> Discussões políticas/religiosas que gerem conflito\n> Flamewar de tecnologia — aqui todo stack tem vez\n> Pornografia ou conteúdo altamente sugestivo",
			},
			{
				Name:  "⚠️ Quebrou as regras?",
				Value: "> Advertência → Remoção do grupo",
			},
			{
				Name:  "📩 Denúncias ou dúvidas?",
				Value: "> Fala no privado com um admin. Tudo será tratado com sigilo.",
			},
			{
				Name:  "​",
				Value: "Reaja com ✅ abaixo para confirmar que leu e aceita as regras e receber acesso ao servidor!",
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Nullspace Lab — Comunidade de Devs",
		},
	}
}

func SetupRules(s *discordgo.Session, cfg *config.Config) {
	msgs, err := s.ChannelMessages(cfg.RulesChannelID, 10, "", "", "")
	if err != nil {
		log.Println("failed to fetch channel messages:", err)
		return
	}

	for _, msg := range msgs {
		if msg.Author.ID == s.State.User.ID {
			_, err = s.ChannelMessageEditEmbed(cfg.RulesChannelID, msg.ID, rulesEmbed())
			if err != nil {
				log.Println("failed to update rules embed:", err)
			}
			log.Println("rules embed updated")
			return
		}
	}

	msg, err := s.ChannelMessageSendEmbed(cfg.RulesChannelID, rulesEmbed())
	if err != nil {
		log.Println("failed to send rules embed:", err)
		return
	}

	s.MessageReactionAdd(cfg.RulesChannelID, msg.ID, "✅")
	log.Println("rules embed created")
}

func ReactionAddHandler(cfg *config.Config) func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	return func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		if r.UserID == s.State.User.ID {
			return
		}

		if r.ChannelID != cfg.RulesChannelID || r.Emoji.Name != "✅" {
			return
		}

		member, err := s.GuildMember(cfg.GuildID, r.UserID)
		if err != nil {
			log.Println("failed to fetch member:", err)
			return
		}

		for _, roleID := range member.Roles {
			if roleID == cfg.NewMemberRoleID || roleID == cfg.OficialMemberRoleID {
				log.Printf("user %s already has the role, skipping\n", r.UserID)
				s.MessageReactionRemove(cfg.RulesChannelID, r.MessageID, "✅", r.UserID)
				return
			}
		}

		err = s.GuildMemberRoleAdd(cfg.GuildID, r.UserID, cfg.NewMemberRoleID)
		if err != nil {
			log.Println("failed to assign role:", err)
			return
		}

		s.MessageReactionRemove(cfg.RulesChannelID, r.MessageID, "✅", r.UserID)
		log.Printf("role assigned to user %s\n", r.UserID)
	}
}
