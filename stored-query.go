package cosmos_ai_queries

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (storedQuery StoredQuery) GetUserAddress() (user sdk.AccAddress, err error) {
	user, errBlack := sdk.AccAddressFromBech32(storedQuery.User)
	return user, errors.Wrapf(errBlack, ErrInvalidUser.Error(), storedQuery.User)
}

func (storedQuery StoredQuery) Validate() (err error) {
	_, err = storedQuery.GetUserAddress()
	if err != nil {
		return err
	}
	return err
}
