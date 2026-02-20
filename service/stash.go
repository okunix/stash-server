package service

import (
	"gitlab.com/stash-password-manager/stash-server/domain/stash"
)

type StashService struct {
	stashRepository   stash.Repository
	stashSessionStore stash.Repository
}
