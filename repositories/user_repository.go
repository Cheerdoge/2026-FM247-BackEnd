package repository

import (
	"2026-FM247-BackEnd/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	result := r.db.Create(user)
	if result != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) UpdateUserInfo(userID uint, username, telenum, gender string) error {
	var user *models.User
	user, err := r.GetUserByID(userID)
	if err != nil {
		return err
	}

	if username != user.Username && username != "" {
		user.Username = username
	}
	if telenum != user.Telenum && telenum != "" {
		user.Telenum = telenum
	}
	if gender != user.Gender && gender != "" {
		user.Gender = gender
	}

	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepository) UpdateUserEmail(userid uint, newEmail string) error {
	var user *models.User
	user, err := r.GetUserByID(userid)
	if err != nil {
		return err
	}
	user.Email = newEmail
	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepository) UpdatePassword(userid uint, newpassword string) error {
	var user *models.User
	user, err := r.GetUserByID(userid)
	if err != nil {
		return err
	}
	user.Password = newpassword
	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepository) DeleteUser(userid uint) error {
	error := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userid).Delete(&models.DailyStudyData{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userid).Delete(&models.MonthlyStudyData{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userid).Delete(&models.TotalStudyData{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userid).Delete(&models.Todo{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userid).Delete(&models.Note{}).Error; err != nil {
			return err
		}
		if err := tx.Where("uploader_id = ?", userid).Delete(&models.Music{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", userid).Delete(&models.User{}).Error; err != nil {
			return err
		}
		return nil
	})
	if error != nil {
		return error
	}
	return nil
}

// UpdateAvatarURL 更新用户头像URL
func (r *UserRepository) UpdateAvatarURL(userID uint, avatarURL string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("avatar", avatarURL).
		Error
}
