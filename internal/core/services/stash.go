package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/okunix/stash-server/internal/core/auth"
	"github.com/okunix/stash-server/internal/core/crypto"
	"github.com/okunix/stash-server/internal/core/domain/secret"
	"github.com/okunix/stash-server/internal/core/domain/stash"
	"github.com/okunix/stash-server/internal/core/dto"
	"github.com/okunix/stash-server/internal/core/ports"
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
	currentUser, err := s.isStashMaintainer(ctx, req.StashID)
	if err != nil {
		return err
	}
	if req.UserID == currentUser.UserID {
		return ports.BadRequestError(errors.New("you can't add yourself as a member"))
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

	innerKDF, _ := crypto.NewArgon2ID()
	masterKey, _ := innerKDF.DeriveKey([]byte(req.Password))
	masterKeySalt := masterKey.Salt()

	outerKDF, _ := crypto.NewArgon2ID()
	masterKeyHash, _ := outerKDF.DeriveKey(masterKey.Bytes())

	initialData := "{}"
	cipher := crypto.AESGCM()
	encryptedData, _ := cipher.Encrypt(masterKey.Bytes(), []byte(initialData))

	_, err := s.stashRepo.CreateStash(ctx, stash.CreateStashParams{
		Name:          req.Name,
		Description:   req.Description,
		MaintainerID:  currentUser.UserID,
		MasterKeyHash: masterKeyHash.String(),
		MasterKeySalt: string(masterKeySalt),
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

	maintainer, err := s.stashRepo.GetStashMaintainer(ctx, stashID)
	if err != nil {
		return nil, ports.InternalError(err)
	}
	res.Maintainer = dto.StashMaintainerResponse{
		UserID:   maintainer.UserID,
		Username: maintainer.Username,
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

func (s *stashService) UpdateStash(
	ctx context.Context,
	stashID uuid.UUID,
	req dto.UpdateStashRequest,
) error {
	st, err := s.getStashIfMemberOrMaintainer(ctx, stashID)
	if err != nil {
		return err
	}
	if problems, ok := req.Validate(); !ok {
		return ports.NewValidationError(problems)
	}
	name := st.Name
	if req.Name != nil {
		name = *req.Name
	}
	desc := st.Description
	if req.Description != nil {
		desc = req.Description
	}
	_, err = s.stashRepo.UpdateStash(ctx,
		stash.UpdateStashParams{
			StashID:     stashID,
			Name:        name,
			Description: desc,
		})
	if err != nil {
		return ports.BadRequestError(err)
	}
	return nil
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
	keys := make([]string, 0, len(secret.Data))
	for k := range secret.Data {
		keys = append(keys, k)
	}
	resp := &dto.SecretResponse{
		Keys:       keys,
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

	innerKDF, _ := crypto.NewArgon2ID(crypto.WithHeader(st.MasterKeySalt))
	masterKey, _ := innerKDF.DeriveKey([]byte(password))
	outerKDF, stashMasterKeyHash, _ := crypto.NewArgon2IDFromString(st.MasterKeyHash)
	masterKeyHash, err := outerKDF.DeriveKey(masterKey.Bytes())
	if err != nil {
		return ports.InternalError(err)
	}
	eq := outerKDF.Compare(stashMasterKeyHash, masterKeyHash.Bytes())
	if !eq {
		return ports.BadRequestError(errors.New("wrong password"))
	}

	cipher := crypto.AESGCM()
	encryptedData, err := base64.RawStdEncoding.DecodeString(st.EncryptedData)
	if err != nil {
		return ports.InternalError(err)
	}
	plaintext, err := cipher.Decrypt(masterKey.Bytes(), encryptedData)
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
		MasterKey:    masterKey.Bytes(),
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

func (s *stashService) GetStashByName(
	ctx context.Context,
	maintainerID uuid.UUID,
	name string,
) (*dto.StashResponse, error) {
	if _, ok := auth.UserFromContext(ctx); !ok {
		return nil, ports.UnauthorizedError(nil)
	}
	st, err := s.stashRepo.GetStashByName(ctx, maintainerID, name)
	if err != nil {
		return nil, ports.NotFoundError(nil)
	}
	if _, err := s.isStashMemberOrMaintainer(ctx, st.ID); err != nil {
		return nil, err
	}
	_, err = s.secretRepo.GetSecretByStashID(ctx, st.ID)
	return dto.NewStashResponse(st, err != nil), nil
}

func (s *stashService) ListMyStashes(ctx context.Context) (*dto.ListMyStashesResponse, error) {
	currentUser, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, ports.UnauthorizedError(nil)
	}
	maintainerStashes, err := s.stashRepo.ListMaintainerStashes(ctx, currentUser.UserID)
	if err != nil {
		return nil, ports.InternalError(err)
	}
	memberStashes, err := s.stashRepo.ListMemberStashes(ctx, currentUser.UserID)
	if err != nil {
		return nil, ports.InternalError(err)
	}

	maintainerStashesResp := make([]*dto.StashResponse, 0, len(maintainerStashes))
	memberStashesResp := make([]*dto.StashResponse, 0, len(memberStashes))
	for _, v := range maintainerStashes {
		_, err = s.secretRepo.GetSecretByStashID(ctx, v.ID)
		maintainerStashesResp = append(
			maintainerStashesResp,
			dto.NewStashResponse(v, err != nil),
		)
	}
	for _, v := range memberStashes {
		_, err := s.secretRepo.GetSecretByStashID(ctx, v.ID)
		memberStashesResp = append(memberStashesResp, dto.NewStashResponse(v, err != nil))
	}
	return &dto.ListMyStashesResponse{
		Maintainer: maintainerStashesResp,
		Member:     memberStashesResp,
	}, nil
}
