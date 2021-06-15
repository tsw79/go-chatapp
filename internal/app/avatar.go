package app

import (
	"chatapp/internal/lib/config"
	"errors"
	"io/ioutil"
	"path"
)

// ErrNoAvatar is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("Chat error: Unable to get an avatar URL.")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client,
	// or returns an error if something goes wrong.
	// ErrNoAvatarURL is returned if the object is unable to get
	// a URL for the specified client.
	GetAvatarURL(ChatUser) (string, error)
}

// `TryAvatars` holds a collection of `Avatar` ---

type AvatarStrategies []Avatar

// TODO should this be here?
var StrategyList Avatar = AvatarStrategies{
	FileSystemAvatarStrategy,
	AuthAvatarStrategy,
	GravatarStrategy,
}

// Loop through all Avatar strategies to find an existing avatar url.
// If found, break and return url, else return error.
func (strategies AvatarStrategies) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatarStrategy := range strategies {
		if url, err := avatarStrategy.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

// AuthAvatar Strategy ---

type AuthAvatar struct{}

var AuthAvatarStrategy AuthAvatar

func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}
	return url, nil
}

// Gravtar  Strategy ---

type GravatarAvatar struct{}

var GravatarStrategy GravatarAvatar

func (GravatarAvatar) GetAvatarURL(usr ChatUser) (string, error) {
	return config.GetInstance().Url.Gravatar + usr.UniqueID(), nil
}

// Filesystem Avatar Strategy ---

type FileSystemAvatar struct{}

var FileSystemAvatarStrategy FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir(config.GetInstance().Dir.Avatars); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := path.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}
