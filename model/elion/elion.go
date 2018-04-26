// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package oauthelion

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	l4g "github.com/alecthomas/log4go"

	"github.com/mattermost/mattermost-server/einterfaces"
	"github.com/mattermost/mattermost-server/model"
)

type ElionProvider struct{}

type ElionUserProfileImage struct {
	Full     string `json:"image_url_full"`
	Large    string `json:"image_url_large"`
	Medium   string `json:"image_url_medium"`
	Small    string `json:"image_url_small"`
	HasImage bool   `json:"has_image"`
}

type ElionUser struct {
	Username       string                `json:"username"`
	Bio            string                `json:"bio"`
	Name           string                `json:"name"`
	Email          string                `json:"email"`
	Country        string                `json:"country"`
	ProfileImage   ElionUserProfileImage `json:"profile_image"`
	YearOfBirth    int                   `json:"year_of_birth"`
	EducationLevel string                `json:"level_of_education"`
	Languages      []map[string]string   `json:"language_proeficiencies"`
	Gender         string                `json:"gender"`
	DateJoined     time.Time             `json:"date_joined"`
}

func init() {
	provider := &ElionProvider{}
	einterfaces.RegisterOauthProvider(model.USER_AUTH_SERVICE_ELION, provider)
}

func userFromElionUser(glu *ElionUser) *model.User {
	user := &model.User{}
	username := glu.Username
	user.Username = model.CleanUsername(username)
	splitName := strings.Split(glu.Name, " ")
	if len(splitName) == 2 {
		user.FirstName = splitName[0]
		user.LastName = splitName[1]
	} else if len(splitName) >= 2 {
		user.FirstName = splitName[0]
		user.LastName = strings.Join(splitName[1:], " ")
	} else {
		user.FirstName = glu.Name
	}
	user.Email = glu.Email
	user.AuthData = &user.Username
	user.AuthService = model.USER_AUTH_SERVICE_ELION
	if len(glu.Languages) > 0 {
		user.Locale = glu.Languages[0]["code"]
	}

	return user
}

func elionUserFromJson(data io.Reader) *ElionUser {
	// Uncomment to take a peek into the JSON response.
	// buf := new(bytes.Buffer)
	// buf.ReadFrom(data)
	// s := buf.String()
	// l4g.Info(s)

	decoder := json.NewDecoder(data)
	var glu []ElionUser
	err := decoder.Decode(&glu)
	if err == nil {
		return &glu[0]
	}

	l4g.Error("failed to decode user data: %v", err)

	return nil
}

func (glu *ElionUser) ToJson() string {
	b, err := json.Marshal(glu)
	if err != nil {
		return ""
	}
	return string(b)
}

func (glu *ElionUser) IsValid() bool {
	if len(glu.Username) == 0 {
		return false
	}

	if len(glu.Email) == 0 {
		return false
	}

	return true
}

func (glu *ElionUser) getAuthData() string {
	return glu.Username
}

func (m *ElionProvider) GetIdentifier() string {
	return model.USER_AUTH_SERVICE_ELION
}

func (m *ElionProvider) GetUserFromJson(data io.Reader) *model.User {
	glu := elionUserFromJson(data)
	if glu != nil && glu.IsValid() {
		return userFromElionUser(glu)
	}

	return &model.User{}
}

func (m *ElionProvider) GetAuthDataFromJson(data io.Reader) string {
	glu := elionUserFromJson(data)

	if glu != nil && glu.IsValid() {
		return glu.getAuthData()
	}

	return ""
}
