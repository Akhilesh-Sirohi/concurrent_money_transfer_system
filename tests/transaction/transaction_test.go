package transaction

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"

	"concurrent_money_transfer_system/internals/transactions"
	"concurrent_money_transfer_system/internals/wallet"
	"concurrent_money_transfer_system/tests"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	tests.Setup()
	setup()
	code := m.Run()
	os.Exit(code)
}

var walletRepo wallet.WalletRepo
var transactionRepo transactions.TransactionRepo

var testData map[string]tests.TestData

func setup() {
	testData = tests.ReadTestData("test_data.json")
	walletRepo = wallet.NewWalletRepo()
	transactionRepo = transactions.NewTransactionRepo()
	transactions.Reset()
	createWallets()
}

var wallets []wallet.Wallet

func createWallets() {
	wallets = make([]wallet.Wallet, 0)
	ctx := context.Background()
	for i := 1; i <= 10; i++ {
		id := fmt.Sprintf("%d", i)
		wal, _ := walletRepo.CreateWallet(ctx, wallet.Wallet{
			ID:       id,
			UserID:   id,
			Balance:  1000000, // 1 Million
			Currency: "USD",
			Status:   wallet.Active,
		})
		wallets = append(wallets, wal)
	}
}

func TestTransferMoney(t *testing.T) {
	test_data := testData["TestTranferMoney"]
	actual_response, _ := tests.MakeRequestAndGetResponse(t, test_data)
	transactions, err := transactionRepo.GetTransactionsByUserID(context.Background(), "1")
	assert.NoError(t, err)
	expected_response := test_data.Response.Body
	expected_response["id"] = transactions[0].ID
	assert.True(t, tests.SelectiveEqual(expected_response, actual_response))

	// Check if the sender wallet balance is updated correctly
	senderWallet, err := walletRepo.GetWalletByUserID(context.Background(), "1")
	assert.NoError(t, err)
	assert.Equal(t, senderWallet.Balance, float64(999000))

	// Check if the receiver wallet balance is updated correctly
	receiverWallet, err := walletRepo.GetWalletByUserID(context.Background(), "2")
	assert.NoError(t, err)
	assert.Equal(t, receiverWallet.Balance, float64(1001000))
}

func TestGetTransaction(t *testing.T) {
	// first create a transaction
	tests.MakeRequestAndGetResponse(t, testData["TestTranferMoney"])

	fmt.Println("Transaction ID: ", "XXX")
	// then get the transaction
	transactions, err := transactionRepo.GetTransactionsByUserID(context.Background(), "1")
	assert.NoError(t, err)

	test_data := testData["TestGetTransaction"]
	test_data.Request.URL = fmt.Sprintf("api/transaction/%s", transactions[0].ID)
	test_data.Response.Body["id"] = transactions[0].ID

	tests.MakeRequestAndValidateResponse(t, test_data)
}

// TODO: fix this test by handling array of transactions
func TestGetTransactionsByUserID(t *testing.T) {
	return
	// first create a transaction
	tests.MakeRequestAndGetResponse(t, testData["TestTranferMoney"])

	fmt.Println("Transaction ID: ", "XXX")
	// then get the transaction
	transactions, err := transactionRepo.GetTransactionsByUserID(context.Background(), "1")
	assert.NoError(t, err)

	test_data := testData["TestGetTransactionsByUserID"]
	test_data.Response.Body["id"] = transactions[0].ID

	tests.MakeRequestAndValidateResponse(t, test_data)
}

func TestValidateNegativeAmountTransferMoneyRequest(t *testing.T) {
	tests.MakeRequestAndValidateResponse(t, testData["TestValidateNegativeAmountTransferMoneyRequest"])
}

func TestValidateSenderAndReceiverSameTransferMoneyRequest(t *testing.T) {
	tests.MakeRequestAndValidateResponse(t, testData["TestValidateSenderAndReceiverSameTransferMoneyRequest"])
}

func TestValidateSenderDoesNotExistTransferMoneyRequest(t *testing.T) {
	tests.MakeRequestAndValidateResponse(t, testData["TestValidateSenderDoesNotExistTransferMoneyRequest"])
}

func TestCannotTransferMoreThanBalance(t *testing.T) {
	// Now make the API request to test the transfer
	tests.MakeRequestAndValidateResponse(t, testData["TestCannotTransferMoreThanBalance"])
}

func TestConcurrentTransferMoney(t *testing.T) {
	setup()
	wg := sync.WaitGroup{}
	wg.Add(50000)
	for i := 0; i < 50000; i++ {
		go func() {
			defer wg.Done()
			transferMoneyRandomly(t)
		}()
	}
	wg.Wait()
	validateTransactionCount(t)
	validateTotalBalanceAcrossAllWallets(t)
	validateEachWalletBalance(t)
}

func transferMoneyRandomly(t *testing.T) {
	receiver_id_index := rand.Intn(10)
	sender_id_index := (receiver_id_index + 1 + rand.Intn(9)) % 10
	receiver_id := fmt.Sprintf("%d", receiver_id_index+1)
	sender_id := fmt.Sprintf("%d", sender_id_index+1)

	amount := 1 + rand.Intn(500)
	data := tests.TestData{
		Request: tests.Request{
			URL:    "api/transaction/transfer",
			Method: "POST",
			Body: map[string]interface{}{
				"sender_id":   sender_id,
				"receiver_id": receiver_id,
				"amount":      float64(amount),
				"currency":    "USD",
				"description": "Test transfer " + strconv.Itoa(amount),
			},
		},
		Response: tests.Response{
			Status: 200,
			Body: map[string]interface{}{
				"debit_user_id":  sender_id,
				"credit_user_id": receiver_id,
				"amount":         float64(amount),
				"currency":       "USD",
				"description":    "Test transfer " + strconv.Itoa(amount),
			},
		},
	}
	responseBodyMap, response := tests.MakeRequestAndGetResponse(t, data)
	if response.Code != 200 {
		panic(fmt.Sprintf("Failed to transfer money: code:%d responseBodyMap:%v", response.Code, responseBodyMap))
	}

	for key, expectedValue := range data.Response.Body {
		actualValue := responseBodyMap[key]
		assert.Equal(t, expectedValue, actualValue)
	}
}

func validateTotalBalanceAcrossAllWallets(t *testing.T) {
	total_balance := 0
	for i := 1; i <= 10; i++ {
		wallet, err := walletRepo.GetWalletByUserID(context.Background(), fmt.Sprintf("%d", i))
		assert.NoError(t, err)
		total_balance += int(wallet.Balance)
	}
	assert.Equal(t, total_balance, 10000000)
}

func validateTransactionCount(t *testing.T) {
	transactions, err := transactionRepo.GetAllTransactions(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, len(transactions), 50000) // 50000 transactions are created
}

func validateEachWalletBalance(t *testing.T) {
	for _, wallet := range wallets {
		all_transactions, err := transactionRepo.GetTransactionsByUserID(context.Background(), wallet.UserID)
		assert.NoError(t, err)
		balance := float64(1000000) // starting balance
		for _, transaction := range all_transactions {
			if transaction.DebitUserID == wallet.UserID {
				balance -= transaction.Amount
			} else {
				balance += transaction.Amount
			}
			assert.Equal(t, transaction.Status, transactions.Completed)
		}

		updatedWallet, err := walletRepo.GetWalletByUserID(context.Background(), wallet.UserID)
		assert.NoError(t, err)
		assert.Equal(t, updatedWallet.Balance, balance)
	}
}
