package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/internal/core/auth"
	"gitlab.com/stash-password-manager/stash-server/internal/core/crypto"
	"gitlab.com/stash-password-manager/stash-server/internal/core/domain/secret"
	"gitlab.com/stash-password-manager/stash-server/internal/core/domain/stash"
	"gitlab.com/stash-password-manager/stash-server/internal/core/dto"
	"gitlab.com/stash-password-manager/stash-server/internal/core/ports"
)

type stashService struct {
	stashRepo  ports.StashRepository
	secretRepo ports.SecretRepository
}

type StashServiceParams struct {
	StashRepository  ports.StashRepository
	SecretRepository ports.SecretRepository
}

func NewStashService(params StashServiceParams) ports.StashService {
	return &stashService{
		stashRepo:  params.StashRepository,
		secretRepo: params.SecretRepository,
	}
}

func (s *stashService) isStashMaintainer(
	ctx context.Context,
	stashID uuid.UUID,
) (*auth.CurrentUser, error) {
	currentUser, ok := auth.UserFromContext(ctx)
	if !ok {
		return currentUser, ports.UnauthorizedError(nil)
	}

	ok, err := s.stashRepo.IsStashMaintainer(ctx, currentUser.UserID, stashID)
	if err != nil {
		return currentUser, ports.InternalError(err)
	}
	if !ok {
		return currentUser, ports.ForbiddenError(errors.New("you are not stash maintainer"))
	}

	return currentUser, nil
}

func (s *stashService) isStashMemberOrMaintainer(
	ctx context.Context,
	stashID uuid.UUID,
) (*auth.CurrentUser, error) {
	currentUser, ok := auth.UserFromContext(ctx)
	if !ok {
		return currentUser, ports.UnauthorizedError(nil)
	}

	ok, _ = s.stashRepo.IsStashMemberOrMaintainer(ctx, currentUser.UserID, stashID)
	if !ok {
		return currentUser, ports.ForbiddenError(errors.New("you are not a member of this stash"))
	}
	return currentUser, nil
}

func (s *stashService) getStashByID(ctx context.Context, stashID uuid.UUID) (*stash.Stash, error) {
	st, err := s.stashRepo.GetStashByID(ctx, stashID)
	if err != nil {
		return nil, ports.NewError(ports.ErrNotFound, errors.New("stash not found"))
	}
	return st, nil
}

func (s *stashService) getStashIfMemberOrMaintainer(
	ctx context.Context,
	stashID uuid.UUID,
) (*stash.Stash, error) {
	currentUser, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, ports.UnauthorizedError(nil)
	}

	ok, _ = s.stashRepo.IsStashMemberOrMaintainer(ctx, currentUser.UserID, stashID)
	if !ok {
		return nil, ports.ForbiddenError(errors.New("you are not a member of this stash"))
	}

	st, err := s.getStashByID(ctx, stashID)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func (s *stashService) getStashIfMaintainer(
	ctx context.Context,
	stashID uuid.UUID,
) (*stash.Stash, error) {
	currentUser, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, ports.UnauthorizedError(nil)
	}

	st, err := s.getStashByID(ctx, stashID)
	if err != nil {
		return nil, err
	}

	if st.MaintainerID != currentUser.UserID {
		return nil, ports.NewError(
			ports.ErrForbidden,
			errors.New("you are not stash maintainer"),
		)
	}

	return st, nil
}

func (s *stashService) AddStashMember(
	ctx context.Context,
	req dto.AddStashMemberRequest,
) error {
	if _, err := s.isStashMaintainer(ctx, req.StashID); err != nil {
		return err
	}
	params := stash.AddMemberParams{UserID: req.UserID, StashID: req.StashID}
	if err := s.stashRepo.AddMember(ctx, params); err != nil {
		return ports.NotFoundError(nil)
	}
	return nil
}

func (s *stashService) RemoveStashMember(
	ctx context.Context,
	req dto.RemoveStashMemberRequest,
) error {
	if _, err := s.isStashMaintainer(ctx, req.StashID); err != nil {
		return err
	}
	params := stash.RemoveMemberParams{UserID: req.UserID, StashID: req.StashID}
	if err := s.stashRepo.RemoveMember(ctx, params); err != nil {
		return ports.NotFoundError(nil)
	}
	return nil
}

func (s *stashService) CreateStash(ctx context.Context, req dto.CreateStashRequest) error {
	currentUser, ok := auth.UserFromContext(ctx)
	if !ok {
		return ports.UnauthorizedError(nil)
	}

	if problems, ok := req.Validate(); !ok {
		return ports.NewValidationError(problems)
	}

	kdf, _ := crypto.NewArgon2ID()
	key, _ := kdf.DeriveKey([]byte(req.MasterKey))
	masterKeyHashString := key.String()

	initialData := "{}"
	cipher := crypto.AESGCM()
	encryptedData, _ := cipher.Encrypt(key.Bytes(), []byte(initialData))

	_, err := s.stashRepo.CreateStash(ctx, stash.CreateStashParams{
		Name:          req.Name,
		Description:   req.Description,
		MaintainerID:  currentUser.UserID,
		MasterKeyHash: masterKeyHashString,
		EncryptedData: base64.RawStdEncoding.EncodeToString(encryptedData),
	})
	if err != nil {
		return ports.InternalError(err)
	}
	return nil
}

func (s *stashService) DeleteStash(ctx context.Context, stashID uuid.UUID) error {
	if _, err := s.isStashMaintainer(ctx, stashID); err != nil {
		return err
	}
	if err := s.stashRepo.DeleteStash(ctx, stashID); err != nil {
		return ports.InternalError(err)
	}
	return nil
}

func (s *stashService) GetStashByID(
	ctx context.Context,
	stashID uuid.UUID,
) (*dto.StashResponse, error) {
	st, err := s.getStashIfMemberOrMaintainer(ctx, stashID)
	if err != nil {
		return nil, err
	}
	_, err = s.secretRepo.GetSecretByStashID(ctx, stashID)
	return dto.NewStashResponse(st, err != nil), nil
}

func (s *stashService) ListStashMembers(
	ctx context.Context,
	stashID uuid.UUID,
) (*dto.ListStashMemberResponse, error) {
	if _, err := s.isStashMemberOrMaintainer(ctx, stashID); err != nil {
		return nil, err
	}

	members, err := s.stashRepo.GetStashMembers(ctx, stashID)
	if err != nil {
		return nil, ports.NotFoundError(errors.New("stash not found"))
	}

	res := new(dto.ListStashMemberResponse)
	res.Members = make([]dto.StashMemberResponse, 0)
	for _, v := range members {
		res.Members = append(res.Members,
			dto.StashMemberResponse{
				UserID:   v.UserID,
				Username: v.Username,
				Since:    v.Since,
			})
	}

	return res, nil
}

func (s *stashService) ListStashes(
	ctx context.Context,
	req dto.ListStashesRequest,
) (*dto.ListStashResponse, error) {
	currentUser, ok := auth.UserFromContext(ctx)
	if !ok {
		return new(dto.ListStashResponse), ports.UnauthorizedError(nil)
	}
	params := stash.ListStashesParams{
		Limit:        req.Limit,
		Offset:       req.Offset,
		Search:       req.Search,
		MaintainerID: currentUser.UserID,
	}
	stashes, total, err := s.stashRepo.ListStashes(ctx, params)
	if err != nil {
		return new(dto.ListStashResponse), ports.NotFoundError(err)
	}
	resp := &dto.ListStashResponse{
		Page: &dto.Page{
			Limit:  req.Limit,
			Offset: req.Offset,
			Total:  total,
		},
		Result: []dto.StashResponse{},
	}
	for _, v := range stashes {
		_, err = s.secretRepo.GetSecretByStashID(ctx, v.ID)
		resp.Result = append(
			resp.Result,
			*dto.NewStashResponse(v, err != nil),
		)
	}
	return resp, nil
}

func (s *stashService) UpdateStash(ctx context.Context, req dto.UpdateStashRequest) error {
	panic("unimplemented")
}

func (s *stashService) GetSecrets(
	ctx context.Context,
	stashID uuid.UUID,
) (*dto.SecretResponse, error) {
	if _, err := s.isStashMemberOrMaintainer(ctx, stashID); err != nil {
		return nil, err
	}
	secret, err := s.secretRepo.GetSecretByStashID(ctx, stashID)
	if err != nil {
		return nil, ports.BadRequestError(err)
	}
	resp := &dto.SecretResponse{
		Data:       secret.Data,
		UnlockedAt: secret.UnlockedAt,
	}
	return resp, nil
}

func (s *stashService) GetSecretsEntry(
	ctx context.Context,
	stashID uuid.UUID,
	entryKey string,
) (string, error) {
	if _, err := s.isStashMemberOrMaintainer(ctx, stashID); err != nil {
		return "", err
	}
	sec, err := s.secretRepo.GetSecretByStashID(ctx, stashID)
	if err != nil {
		return "", ports.NotFoundError(err)
	}
	entry, ok := sec.GetEntry(entryKey)
	if !ok {
		return "", ports.NotFoundError(errors.New("secrets entry not found"))
	}
	return entry, nil
}

func (s *stashService) AddSecretsEntry(
	ctx context.Context,
	stashID uuid.UUID,
	req dto.AddSecret,
) error {
	_, err := s.isStashMaintainer(ctx, stashID)
	if err != nil {
		return err
	}

	if problems, ok := req.Validate(); !ok {
		return ports.NewValidationError(problems)
	}

	sec, err := s.secretRepo.GetSecretByStashID(ctx, stashID)
	if err != nil {
		return ports.NotFoundError(err)
	}
	sec.AddEntry(req.Name, req.Value)
	err = s.secretRepo.UpdateSecret(ctx, stashID, sec)
	if err != nil {
		return ports.InternalError(err)
	}
	go s.commitDataUpdate(context.Background(), stashID, sec)
	return nil
}

func (s *stashService) RemoveSecretsEntry(
	ctx context.Context,
	stashID uuid.UUID,
	entryKey string,
) error {
	_, err := s.isStashMaintainer(ctx, stashID)
	if err != nil {
		return err
	}
	sec, err := s.secretRepo.GetSecretByStashID(ctx, stashID)
	if err != nil {
		return ports.NotFoundError(err)
	}
	sec.RemoveEntry(entryKey)
	err = s.secretRepo.UpdateSecret(ctx, stashID, sec)
	if err != nil {
		return ports.InternalError(err)
	}
	go s.commitDataUpdate(context.Background(), stashID, sec)
	return nil
}

func (s *stashService) ListUnlockedSecrets(ctx context.Context) ([]*dto.StashResponse, error) {
	return nil, nil
}

func (s *stashService) Unlock(ctx context.Context, stashID uuid.UUID, password string) error {
	st, err := s.getStashIfMaintainer(ctx, stashID)
	if err != nil {
		return err
	}

	kdf, _, _ := crypto.NewArgon2IDFromString(st.MasterKeyHash)
	key, err := kdf.DeriveKey([]byte(password))
	if err != nil {
		return ports.InternalError(err)
	}
	cipher := crypto.AESGCM()
	encryptedData, err := base64.RawStdEncoding.DecodeString(st.EncryptedData)
	if err != nil {
		return ports.InternalError(err)
	}
	plaintext, err := cipher.Decrypt(key.Bytes(), encryptedData)
	if err != nil {
		return ports.BadRequestError(errors.New("failed to decrypt data"))
	}
	var data map[string]string
	if err := json.Unmarshal(plaintext, &data); err != nil {
		return ports.InternalError(err)
	}
	params := secret.AddSecretParams{
		StashID:      st.ID,
		MaintainerID: st.MaintainerID,
		MasterKey:    key.Bytes(),
		Data:         data,
	}
	if _, err := s.secretRepo.AddSecret(ctx, params); err != nil {
		return ports.InternalError(err)
	}
	return nil
}

func (s *stashService) Lock(ctx context.Context, stashID uuid.UUID) error {
	if _, err := s.isStashMaintainer(ctx, stashID); err != nil {
		return err
	}

	sec, err := s.secretRepo.GetSecretByStashID(ctx, stashID)
	if err != nil {
		return ports.NotFoundError(err)
	}

	cipher := crypto.AESGCM()
	dataBytes, _ := json.Marshal(sec.Data)
	ciphertext, err := cipher.Encrypt(sec.MasterKey, dataBytes)
	if err != nil {
		return ports.InternalError(err)
	}
	params := stash.CommitDataParams{
		StashID: stashID,
		Data:    base64.RawStdEncoding.EncodeToString(ciphertext),
	}
	if err := s.stashRepo.CommitData(ctx, params); err != nil {
		return ports.InternalError(err)
	}

	if _, err := s.secretRepo.RemoveSecretByStashID(ctx, stashID); err != nil {
		return ports.BadRequestError(err)
	}
	return nil
}

func (s *stashService) commitDataUpdate(
	ctx context.Context,
	stashID uuid.UUID,
	sec *secret.Secret,
) error {
	cipher := crypto.AESGCM()
	dataBytes, err := sec.Seal(cipher)
	if err != nil {
		return ports.InternalError(errors.New("secrets seal failed"))
	}
	params := stash.CommitDataParams{
		StashID: stashID,
		Data:    base64.RawStdEncoding.EncodeToString(dataBytes),
	}
	err = s.stashRepo.CommitData(ctx, params)
	if err != nil {
		return ports.InternalError(err)
	}
	return nil
}
