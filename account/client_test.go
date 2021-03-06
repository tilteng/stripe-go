package account

import (
	"testing"

	stripe "github.com/tilteng/stripe-go"
	"github.com/tilteng/stripe-go/bankaccount"
	"github.com/tilteng/stripe-go/card"
	"github.com/tilteng/stripe-go/currency"
	"github.com/tilteng/stripe-go/recipient"
	"github.com/tilteng/stripe-go/token"
	. "github.com/tilteng/stripe-go/utils"
)

func init() {
	stripe.Key = GetTestKey()
}

func TestAccountNew(t *testing.T) {
	params := &stripe.AccountParams{
		Managed:              true,
		Country:              "CA",
		BusinessUrl:          "www.stripe.com",
		BusinessName:         "Stripe",
		BusinessPrimaryColor: "#ffffff",
		DebitNegativeBal:     true,
		SupportEmail:         "foo@bar.com",
		SupportUrl:           "www.stripe.com",
		SupportPhone:         "4151234567",
		LegalEntity: &stripe.LegalEntity{
			Type:         stripe.Individual,
			BusinessName: "Stripe Go",
			DOB: stripe.DOB{
				Day:   1,
				Month: 2,
				Year:  1990,
			},
		},
		TOSAcceptance: &stripe.TOSAcceptanceParams{
			IP:        "127.0.0.1",
			Date:      1437578361,
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/600.7.12 (KHTML, like Gecko) Version/8.0.7 Safari/600.7.12",
		},
	}

	_, err := New(params)
	if err != nil {
		t.Error(err)
	}
}

func TestAccountLegalEntity(t *testing.T) {
	params := &stripe.AccountParams{
		Managed: true,
		Country: "US",
		LegalEntity: &stripe.LegalEntity{
			Type:          stripe.Company,
			BusinessTaxID: "111111",
			SSN:           "1111",
			PersonalID:    "111111111",
			DOB: stripe.DOB{
				Day:   1,
				Month: 2,
				Year:  1990,
			},
		},
	}

	target, err := New(params)
	if err != nil {
		t.Error(err)
	}

	if !target.LegalEntity.BusinessTaxIDProvided {
		t.Errorf("Account is missing BusinessTaxIDProvided even though we submitted the value.\n")
	}

	if !target.LegalEntity.SSNProvided {
		t.Errorf("Account is missing SSNProvided even though we submitted the value.\n")
	}

	if !target.LegalEntity.PersonalIDProvided {
		t.Errorf("Account is missing PersonalIDProvided even though we submitted the value.\n")
	}
}

func TestAccountDelete(t *testing.T) {
	params := &stripe.AccountParams{
		Managed:              true,
		Country:              "CA",
		BusinessUrl:          "www.stripe.com",
		BusinessName:         "Stripe",
		BusinessPrimaryColor: "#ffffff",
		SupportEmail:         "foo@bar.com",
		SupportUrl:           "www.stripe.com",
		SupportPhone:         "4151234567",
		LegalEntity: &stripe.LegalEntity{
			Type:         stripe.Individual,
			BusinessName: "Stripe Go",
			DOB: stripe.DOB{
				Day:   1,
				Month: 2,
				Year:  1990,
			},
		},
	}

	acct, err := New(params)
	if err != nil {
		t.Error(err)
	}

	acctDel, err := Del(acct.ID)
	if err != nil {
		t.Error(err)
	}

	if !acctDel.Deleted {
		t.Errorf("Account id %q expected to be marked as deleted on the returned resource\n", acctDel.ID)
	}
}

func TestAccountReject(t *testing.T) {
	params := &stripe.AccountParams{
		Managed:              true,
		Country:              "CA",
		BusinessUrl:          "www.stripe.com",
		BusinessName:         "Stripe",
		BusinessPrimaryColor: "#ffffff",
		SupportEmail:         "foo@bar.com",
		SupportUrl:           "www.stripe.com",
		SupportPhone:         "4151234567",
		LegalEntity: &stripe.LegalEntity{
			Type:         stripe.Individual,
			BusinessName: "Stripe Go",
			DOB: stripe.DOB{
				Day:   1,
				Month: 2,
				Year:  1990,
			},
		},
	}

	acct, err := New(params)
	if err != nil {
		t.Error(err)
	}

	rejectedAcct, err := Reject(acct.ID, &stripe.AccountRejectParams{Reason: "fraud"})
	if err != nil {
		t.Error(err)
	}

	if rejectedAcct.Verification.DisabledReason != "rejected.fraud" {
		t.Error("Account DisabledReason did not change to rejected.fraud.")
	}
}

func TestAccountMigrateFromRecipients(t *testing.T) {
	recipientParams := &stripe.RecipientParams{
		Name:  "Recipient Name",
		Type:  "individual",
		TaxID: "000000000",
		Email: "a@b.com",
		Desc:  "Recipient Desc",
		Bank: &stripe.BankAccountParams{
			Country: "US",
			Routing: "110000000",
			Account: "000123456789",
		},
		Card: &stripe.CardParams{
			Name:   "Test Debit",
			Number: "4000056655665556",
			Month:  "10",
			Year:   "20",
		},
	}

	target, err := recipient.New(recipientParams)
	if err != nil {
		t.Error(err)
	}

	target2, err := New(&stripe.AccountParams{FromRecipient: target.ID})
	if err != nil {
		t.Error(err)
	}

	target, err = recipient.Get(target.ID, nil)
	if err != nil {
		t.Error(err)
	}

	if target2.ID != target.MigratedTo.ID {
		t.Errorf("The new account ID %v does not match the MigratedTo property %v", target2.ID, target.MigratedTo.ID)
	}
}

func TestAccountGetByID(t *testing.T) {
	params := &stripe.AccountParams{
		Managed: true,
		Country: "CA",
	}

	acct, _ := New(params)

	_, err := GetByID(acct.ID, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestAccountUpdate(t *testing.T) {
	params := &stripe.AccountParams{
		Managed:          true,
		Country:          "CA",
		DebitNegativeBal: true,
	}

	acct, _ := New(params)

	if acct.DebitNegativeBal != true {
		t.Error("debit_negative_balance was not set to true")
	}

	params = &stripe.AccountParams{
		Statement:          "Stripe Go",
		NoDebitNegativeBal: true,
	}

	acct, err := Update(acct.ID, params)
	if err != nil {
		t.Error(err)
	}

	if acct.DebitNegativeBal != false {
		t.Error("debit_negative_balance was not set to false")
	}
}

func TestAccountUpdateLegalEntity(t *testing.T) {
	params := &stripe.AccountParams{
		Managed: true,
		Country: "CA",
		LegalEntity: &stripe.LegalEntity{
			Address: stripe.Address{
				Country: "CA",
				City:    "Montreal",
				Zip:     "H2Y 1C6",
				Line1:   "275, rue Notre-Dame Est",
				State:   "QC",
			},
		},
	}

	acct, err := New(params)

	if err != nil {
		t.Error(err)
	}

	params = &stripe.AccountParams{
		LegalEntity: &stripe.LegalEntity{
			Address: stripe.Address{
				Line1: "321, rue Notre-Dame Est",
			},
		},
	}

	acct, err = Update(acct.ID, params)
	if err != nil {
		t.Error(err)
	}

	if acct.LegalEntity.Address.Line1 != params.LegalEntity.Address.Line1 {
		t.Errorf("The account address line1 %v does not match the params address line1: %v", acct.LegalEntity.Address.Line1, params.LegalEntity.Address.Line1)
	}
}

func TestAccountUpdateWithBankAccount(t *testing.T) {
	params := &stripe.AccountParams{
		Managed: true,
		Country: "CA",
	}

	acct, _ := New(params)

	params = &stripe.AccountParams{
		ExternalAccount: &stripe.AccountExternalAccountParams{
			Country:  "US",
			Currency: "usd",
			Routing:  "110000000",
			Account:  "000123456789",
		},
	}

	_, err := Update(acct.ID, params)
	if err != nil {
		t.Error(err)
	}
}

func TestAccountAddExternalAccountsDefault(t *testing.T) {
	params := &stripe.AccountParams{
		Managed: true,
		Country: "CA",
		ExternalAccount: &stripe.AccountExternalAccountParams{
			Country:  "US",
			Currency: "usd",
			Routing:  "110000000",
			Account:  "000123456789",
		},
	}

	acct, _ := New(params)

	ba, err := bankaccount.New(&stripe.BankAccountParams{
		AccountID: acct.ID,
		Country:   "US",
		Currency:  "usd",
		Routing:   "110000000",
		Account:   "000111111116",
		Default:   true,
	})

	if err != nil {
		t.Error(err)
	}

	if ba.Default == false {
		t.Error("The new external account should be the default but isn't.")
	}

	baTok, err := token.New(&stripe.TokenParams{
		Bank: &stripe.BankAccountParams{
			Country:  "US",
			Currency: "usd",
			Routing:  "110000000",
			Account:  "000333333335",
		},
	})
	if err != nil {
		t.Error(err)
	}

	ba2, err := bankaccount.New(&stripe.BankAccountParams{
		AccountID: acct.ID,
		Token:     baTok.ID,
		Default:   true,
	})

	if err != nil {
		t.Error(err)
	}

	if ba2.Default == false {
		t.Error("The third external account should be the default but isn't.")
	}
}

func TestAccountUpdateWithToken(t *testing.T) {
	params := &stripe.AccountParams{
		Managed: true,
		Country: "CA",
	}

	acct, _ := New(params)

	tokenParams := &stripe.TokenParams{
		Bank: &stripe.BankAccountParams{
			Country: "US",
			Routing: "110000000",
			Account: "000123456789",
		},
	}

	tok, _ := token.New(tokenParams)

	params = &stripe.AccountParams{
		ExternalAccount: &stripe.AccountExternalAccountParams{
			Token: tok.ID,
		},
	}

	_, err := Update(acct.ID, params)
	if err != nil {
		t.Error(err)
	}
}

func TestAccountUpdateWithCardToken(t *testing.T) {
	params := &stripe.AccountParams{
		Managed: true,
		Country: "US",
	}

	acct, _ := New(params)

	tokenParams := &stripe.TokenParams{
		Card: &stripe.CardParams{
			Number:   "4000056655665556",
			Month:    "06",
			Year:     "20",
			Currency: "usd",
		},
	}

	tok, _ := token.New(tokenParams)

	cardParams := &stripe.CardParams{
		Account: acct.ID,
		Token:   tok.ID,
	}

	c, err := card.New(cardParams)

	if err != nil {
		t.Error(err)
	}

	if c.Currency != currency.USD {
		t.Errorf("Currency %v does not match expected value %v\n", c.Currency, currency.USD)
	}
}

func TestAccountGet(t *testing.T) {
	target, err := Get()

	if err != nil {
		t.Error(err)
	}

	if len(target.ID) == 0 {
		t.Errorf("Account is missing id\n")
	}

	if len(target.Country) == 0 {
		t.Errorf("Account is missing country\n")
	}

	if len(target.DefaultCurrency) == 0 {
		t.Errorf("Account is missing default currency\n")
	}

	if len(target.Name) == 0 {
		t.Errorf("Account is missing name\n")
	}

	if len(target.Email) == 0 {
		t.Errorf("Account is missing email\n")
	}

	if len(target.Timezone) == 0 {
		t.Errorf("Account is missing timezone\n")
	}

	if len(target.Statement) == 0 {
		t.Errorf("Account is missing Statement\n")
	}

	if len(target.BusinessName) == 0 {
		t.Errorf("Account is missing business name\n")
	}

	if len(target.BusinessPrimaryColor) == 0 {
		t.Errorf("Account is missing business primary color\n")
	}

	if len(target.BusinessUrl) == 0 {
		t.Errorf("Account is missing business URL\n")
	}

	if len(target.SupportPhone) == 0 {
		t.Errorf("Account is missing support phone\n")
	}

	if len(target.SupportEmail) == 0 {
		t.Errorf("Account is missing support email\n")
	}

	if len(target.SupportUrl) == 0 {
		t.Errorf("Account is missing support URL\n")
	}

	if len(target.DefaultCurrency) == 0 {
		t.Errorf("Account is missing default currency\n")
	}

	if len(target.Name) == 0 {
		t.Errorf("Account is missing name\n")
	}

	if len(target.Email) == 0 {
		t.Errorf("Account is missing email\n")
	}

	if len(target.Timezone) == 0 {
		t.Errorf("Account is missing timezone\n")
	}
}
