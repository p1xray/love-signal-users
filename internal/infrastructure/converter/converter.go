package converter

import (
	"github.com/guregu/null/v6"
	"love-signal-users/internal/dto"
	"love-signal-users/internal/entity"
	"love-signal-users/internal/enum"
	"love-signal-users/internal/infrastructure"
	"love-signal-users/internal/infrastructure/kafka/data"
	"love-signal-users/internal/infrastructure/storage/models"
)

func ToUserDTO(user models.User) dto.User {
	return dto.User{
		ID:            user.ID,
		ExternalID:    user.ExternalID,
		FullName:      user.FullName,
		Gender:        enum.GenderFromNullInt16(user.Gender),
		DateOfBirth:   user.DateOfBirth.Ptr(),
		AvatarFileKey: user.AvatarFileKey.Ptr(),
	}
}

func ToFollowDTO(follow models.Follow, users []models.User) (dto.Follow, error) {
	followingUser, err := findUserByID(users, follow.FollowingUserID)
	if err != nil {
		return dto.Follow{}, err
	}

	followedUser, err := findUserByID(users, follow.FollowedUserID)
	if err != nil {
		return dto.Follow{}, err
	}

	return dto.Follow{
		ID:            follow.ID,
		FollowingUser: ToUserDTO(followingUser),
		FollowedUser:  ToUserDTO(followedUser),
		NumberOfLikes: follow.NumberOfLikes,
	}, nil
}

func ToFollowStorage(follow *entity.Follow, setters ...models.FollowOption) models.Follow {
	followStorage := models.Follow{
		ID:              follow.ID,
		FollowingUserID: follow.FollowingUser.ID,
		FollowedUserID:  follow.FollowedUser.ID,
		NumberOfLikes:   follow.NumberOfLikes,
	}

	for _, setter := range setters {
		setter(&followStorage)
	}

	return followStorage
}

func ToUserStorage(user data.User, setters ...models.UserOption) models.User {
	userStorage := models.User{
		ExternalID:    user.ID,
		FullName:      user.FullName,
		DateOfBirth:   null.TimeFromPtr(user.DateOfBirth),
		Gender:        user.Gender.ToNullInt16(),
		AvatarFileKey: null.StringFromPtr(user.AvatarFileKey),
	}

	for _, setter := range setters {
		setter(&userStorage)
	}

	return userStorage
}

func findUserByID(users []models.User, id int64) (models.User, error) {
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return models.User{}, infrastructure.ErrEntityNotFound
}
