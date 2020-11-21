package email

import "github.com/bsladewski/mojito/data"

// getEmailTemplateByTitle retrieves an email template record by its title.
func getEmailTemplateByTitle(templateTitle TemplateTitle) (*emailTemplate, error) {

	var item emailTemplate

	if err := data.DB().Model(&emailTemplate{}).
		Where("title = ?", templateTitle).
		First(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil

}
