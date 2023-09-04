package helpers

import (
	"admin_api/db"
	"strings"
)

func ExtractMentionedStudentEmails(notification string) []string {
	words := strings.Fields(notification)
	var mentionedStudents []string
	for _, word := range words {
		if strings.HasPrefix(word, "@") && strings.Contains(word, "@") {
			studentEmail := strings.TrimPrefix(word, "@")
			_, _, err := db.CheckStudentExist(db.DB, studentEmail)
			if err != nil {
				continue
			}

			mentionedStudents = append(mentionedStudents, studentEmail)
		}
	}
	return mentionedStudents
}

func RemoveDuplicates(registeredStudentEmails []string, nonSuspendedMentionedStudentEmails []string) []string {
	uniqueEmails := make(map[string]bool)
	for _, email := range registeredStudentEmails {
		uniqueEmails[email] = true
	}
	for _, email := range nonSuspendedMentionedStudentEmails {
		uniqueEmails[email] = true
	}
	var result []string
	for email := range uniqueEmails {
		result = append(result, email)
	}

	return result
}
