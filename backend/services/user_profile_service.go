package services

import (
	"context"
	"crypto-wallet/config"
	"errors"

	"cloud.google.com/go/firestore"
)

func GetUserProfileService(walletId string) (map[string]interface{}, error) {
	ctx := context.Background()

	// find user WHERE user.wallet_id == walletId
	iter := config.Firestore.Collection("users").
		Where("wallet_id", "==", walletId).
		Limit(1).
		Documents(ctx)

	doc, err := iter.Next()
	if err != nil {
		return nil, errors.New("profile not found")
	}

	user := doc.Data()

	// fetch wallet for keys + beneficiaries
	walletDoc, err := config.Firestore.Collection("wallets").
		Doc(walletId).
		Get(ctx)

	if err == nil && walletDoc.Exists() {
		user["public_key"] = walletDoc.Data()["public_key"]
		user["private_key_enc"] = walletDoc.Data()["private_key_enc"]
		user["beneficiaries"] = walletDoc.Data()["beneficiaries"]
	}

	return user, nil
}

func UpdateUserProfileService(walletId, name, email, cnic string) error {

	ctx := context.Background()

	// locate user by walletID
	iter := config.Firestore.Collection("users").
		Where("wallet_id", "==", walletId).
		Limit(1).
		Documents(ctx)

	doc, err := iter.Next()
	if err != nil {
		return errors.New("user not found")
	}

	_, err = doc.Ref.Update(ctx, []firestore.Update{
		{Path: "name", Value: name},
		{Path: "email", Value: email},
		{Path: "cnic", Value: cnic},
	})

	return err
}
