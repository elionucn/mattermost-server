// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package oauthelion

import (
	"encoding/json"
	"io"
	"strconv"
	"strings"

	"github.com/mattermost/mattermost-server/einterfaces"
	"github.com/mattermost/mattermost-server/model"
)

type ElionProvider struct{}

type ElionUser struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

func init() {
	provider := &ElionProvider{}
	einterfaces.RegisterOauthProvider(model.USER_AUTH_SERVICE_ELION, provider)
}

func userFromElionUser(glu *ElionUser) *model.User {
	user := &model.User{}
	username := glu.Username
	if username == "" {
		username = glu.Login
	}
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
	userId := strconv.FormatInt(glu.Id, 10)
	user.AuthData = &userId
	user.AuthService = model.USER_AUTH_SERVICE_ELION

	return user
}

func elionUserFromJson(data io.Reader) *ElionUser {
	decoder := json.NewDecoder(data)
	var glu ElionUser
	err := decoder.Decode(&glu)
	if err == nil {
		return &glu
	} else {
		return nil
	}
}

func (glu *ElionUser) ToJson() string {
	b, err := json.Marshal(glu)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func (glu *ElionUser) IsValid() bool {
	if glu.Id == 0 {
		return false
	}

	if len(glu.Email) == 0 {
		return false
	}

	return true
}

func (glu *ElionUser) getAuthData() string {
	return strconv.FormatInt(glu.Id, 10)
}

func (m *ElionProvider) GetIdentifier() string {
	return model.USER_AUTH_SERVICE_ELION
}

func (m *ElionProvider) GetUserFromJson(data io.Reader) *model.User {
	glu := elionUserFromJson(data)
	if glu.IsValid() {
		return userFromElionUser(glu)
	}

	return &model.User{}
}

func (m *ElionProvider) GetAuthDataFromJson(data io.Reader) string {
	glu := elionUserFromJson(data)

	if glu.IsValid() {
		return glu.getAuthData()
	}

	return ""
}
