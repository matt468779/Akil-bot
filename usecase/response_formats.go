package usecase

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Opportunity struct {
	Data    []OpportunityData `json:"data"`
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Errors  []error           `json:"errors"`
}

type OpportunityData struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Responsibilities string    `json:"responsibilities"`
	IdealCandidate   string    `json:"idealCandidate"`
	Categories       []string  `json:"categories"`
	OpType           string    `json:"opType"`
	StartDate        time.Time `json:"startDate"`
	EndDate          time.Time `json:"endDate"`
	Deadline         time.Time `json:"deadline"`
	Location         []string  `json:"location"`
	RequiredSkills   []string  `json:"requiredSkills"`
	WhenAndWhere     string    `json:"whenAndWhere"`
	DatePosted       time.Time `json:"datePosted"`
	Status           string    `json:"status"`
	ApplicantsCount  int       `json:"applicantsCount"`
}

func FormatOpportunitiesMarkdown(opportunityData []OpportunityData, baseURL string) string {
	if len(opportunityData) == 0 {
		return "I can't find any opportunities with the given filter \nPlease try agian with other filters"
	}

	var builder strings.Builder

	for _, opportunityData := range opportunityData {

		// Title
		builder.WriteString(fmt.Sprintf("🌟 *[%s](%sorganization/opportunities/%s)* 🌟\n", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, opportunityData.Title), tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, baseURL), tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, opportunityData.ID)))

		// Description
		description := opportunityData.Description
		if len(description) > 200 {
			description = description[:200] + "..."
		}
		builder.WriteString(tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, description) + "\n\n")

		// Responsibilities
		builder.WriteString("✅ *Responsibilities:*\n")
		builder.WriteString(tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, opportunityData.Responsibilities) + "\n\n")

		// Categories
		builder.WriteString("📚 *Categories:*\n")
		builder.WriteString(tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, strings.Join(opportunityData.Categories, ", ")) + "\n\n")

		// Start and end dates
		builder.WriteString("📅 *Start Date:* ")
		builder.WriteString(tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, opportunityData.StartDate.Format("January 2, 2006")) + "\n")
		builder.WriteString("📅 *End Date:* ")
		builder.WriteString(tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, opportunityData.EndDate.Format("January 2, 2006")) + "\n\n")

		// Deadline
		builder.WriteString("⏰ *Deadline:* ")
		builder.WriteString(tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, opportunityData.Deadline.Format("January 2, 2006")) + "\n\n")

		// Location
		builder.WriteString("📍 *Location:* ")
		builder.WriteString(tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, strings.Join(opportunityData.Location, ", ")) + "\n\n")

		// Status and applicants count
		builder.WriteString("📢 *Status:* ")
		builder.WriteString(tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, opportunityData.Status) + "\n")
		builder.WriteString("👥 *Applicants Count:* ")
		builder.WriteString(fmt.Sprintf("%d\n", opportunityData.ApplicantsCount))

		builder.WriteString(tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, strings.Repeat("-", 50))) // Separator between job entries
		builder.WriteString("\n\n")
	}

	return builder.String()
}
