package bot

import (
	"fmt"
	"job-aggregator/internal/models"
	"strings"
)

func FormatVacancy(v models.Vacancy) string {

	cleanTitle := strings.TrimSpace(v.Title)

	var b strings.Builder

	b.WriteString(fmt.Sprintf("🚀 <b>%s</b>\n\n", cleanTitle))

	b.WriteString(fmt.Sprintf("🏢 <b>Компания:</b> <code>%s</code>\n", v.Company))

	if v.Salary != "" && v.Salary != "Зарплата не указана" {
		b.WriteString(fmt.Sprintf("💰 <b>Зарплата:</b> <b>%s</b>\n", v.Salary))
	} else {
		b.WriteString("💰 <b>Зарплата:</b> <i>по результатам собеседования</i>\n")
	}

	b.WriteString("⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯\n")

	if v.Description != "" {

		prettyStack := strings.ReplaceAll(v.Description, " • ", "  ")
		b.WriteString(fmt.Sprintf("🛠 <b>Стек и требования:</b>\n<i>%s</i>\n\n", prettyStack))
	}

	b.WriteString(fmt.Sprintf("🔗 <a href='%s'>Посмотреть на Habr Career</a>\n\n", v.URL))
	b.WriteString("#Habr #Golang #Job")

	return b.String()
}
