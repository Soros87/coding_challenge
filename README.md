# Switcheo coding challenge

**Problem 6: Transaction Broadcaster Service**

**Objective** of a transaction broadcast service: Relay transaction information across decentralised network in a timely and secure manner to all relevant nodes.

Specifications:

1. **System Architecture:**

Define components - (see example diagram below). Client interacts with the smart contract and triggers an event that broadcast the transaction to the distributed ledger which is a **network of multiple nodes.** There are mainly 2 types of nodes.

   1. broadcaster nodes responsible for relaying transactions and storing information of the blockchain
   2. Validating nodes that is responsible for validating transactions and generating new blocks on the Blockchain

![Untitled](https://prod-files-secure.s3.us-west-2.amazonaws.com/1fa197a5-27ce-4f9e-8827-47d2efcd1a80/a525ce58-0c39-4cf8-9e58-7bc206edc8af/Untitled.png)

In our application layer, we need to create the following components.

**TransactionBroadcaster:**

1. Determine the structure of the transaction broadcaster which consist of a broadcastManager field of **Type** BroadcastManager.
2. Declare the function broadcastTransaction that is part of the receiver (tb \* TransactionBroadcaster).
3. Takes in argument called requestData with a map of string keys and values: interface
4. Extracts the value associated with the key "message_type" from the requestData map and asserts its type as a string. Assign to messageType
5. Extracts the value associated with the key "data" from the requestData map and asserts its type as a string.Assign to transactionData
6. Next you sign a transaction with messageType and transactionData as arguments and store results to signedTransaction
7. Access the broadcastManager field of tb (the receiver) and execute the broadcast method.

```go
// TransactionBroadcaster handles broadcasting transactions
type TransactionBroadcaster struct {
broadcastManager *BroadcastManager
}

func (tb *TransactionBroadcaster) broadcastTransaction(requestData map[string]interface{}) string {
	// Validate and extract necessary information from the request_data
	messageType := requestData["message_type"].(string)
	transactionData := requestData["data"].(string)

	// Sign the transaction data (assuming there is a signTransaction function)
	signedTransaction := signTransaction(messageType, transactionData)

	// Initiate the broadcasting process
	result := tb.broadcastManager.broadcast(signedTransaction)

	return result
}
```

**Broadcast Manager:**

- Manages the process of broadcasting signed transactions to the EVM-compatible blockchain network.
- Handles automatic retry in case of failure.
- Logs the transaction details and results for monitoring.

1. First determine the structure of Broadcast Manager which consists of a
   1. blockchainRPCClient which is an entity that facilitates communication between the application and the blockchain network through Remote Procedure Call (RPC) requests. It sends signed transactions to the blockchain network.
   2. transactionStatusTracker which is a TransactionStatusTracker object to store as the name suggests transaction statuses such as success, failure, pending etc.
   3. maxRetries which is a counter
2. define the function broadcast - which accepts the signedTransaction object as an argument.
3. the broadcastTransaction method is called from the instance of blockchainRPCClient where it handles the various scenarios
   1. 1% of the time, it does not respond earlier than 30 seconds.
   2. 95% of the time it responds with a success code within 20-30 seconds.
   3. The rest of the time it returns a failure code.
4. Depending on the result in pt3, the handleResponse method is called and assigned to the result.

```go
// BroadcastManager manages the process of broadcasting transactions
type BroadcastManager struct {
	blockchainRPCClient      *BlockchainRPCClient,
	transactionStatusTracker *TransactionStatusTracker,
	maxRetries:               int,
}

func (bm *BroadcastManager) broadcast(signedTransaction string) string {
// Try broadcasting the signed transaction with automatic retry on failure
	for attempt := 1; attempt <= bm.maxRetries; attempt++ {
		result := bm.blockchainRPCClient.broadcastTransaction(signedTransaction)

		// Log transaction details and status
		bm.transactionStatusTracker.logTransaction(signedTransaction, result)

		// Check if the transaction was successful
		if result == "success" {
			return result
		}

		// Sleep before the next retry, you may implement an exponential backoff strategy
		time.Sleep(2 * time.Second)
	}

	return "failure"
}

// BlockchainRPCClient handles communication with the blockchain node
type BlockchainRPCClient struct {
	// Add any necessary configuration or dependencies
}

func (rpcClient *BlockchainRPCClient) broadcastTransaction(signedTransaction string) string {
	// Make an RPC request to the blockchain node to broadcast the signed transaction
	response := makeRPCRequest(signedTransaction)

	// Handle response time and codes (success, failure)
	result := handleResponse(response)

	return result
}

func makeRPCRequest(signedTransaction string) string {
	// Simulate the RPC request to the blockchain node
	// In a real implementation, this function would send an actual RPC request to the blockchain node
	// and return the response.
	// For the sake of the example, it simulates success, failure, and a delayed response.

	// Simulate 1% of cases where response time is greater than 30 seconds
 //generates a random number between 0 and 99.
	if time.Now().Unix()%100 == 0 {
		time.Sleep(31 * time.Second)
	}

	// Simulate 95% success response within 20-30 seconds
	if time.Now().Unix()%100 > 1 && time.Now().Unix()%100 <= 96 {
		return "success"
	}

	// Simulate the rest of the cases as failure
	return "failure"
}

func handleResponse(response string) string {
	// Placeholder for handling response
	if response == "failure" {
		return "failure"
	}

	// Simulate 1% of cases where response time is greater than 30 seconds
	if time.Now().Unix()%100 == 0 {
		time.Sleep(31 * time.Second)
	}

	return "success"
}
```

1. **Transaction Status Tracker -** Display list of transactions that passed or fail.

- TransactionStatusTracker struct is a slice where each element is a map with keys strings and value of any type (interface {})
- logTransaction method of the receiver tst takes the signedTransaction and status as an argument and returns a slice of transactions appending the new transactionInfo object created.
- Define the getTransactionStatuses function that returns the tst.transaction slice
- getFailedTransactions function filters the tst.transaction slice for “status” that is equal to “failure” and returns a new slice called failedTransactions.

```go
// TransactionStatusTracker keeps track of the status of each transaction
type TransactionStatusTracker struct {
	transactions []map[string]interface{}
}

func (tst *TransactionStatusTracker) logTransaction(signedTransaction string, status string) {
	// Log transaction details and status
	transactionInfo := map[string]interface{}{
		"timestamp":          time.Now().Unix(),
		"signed_transaction": signedTransaction,
		"status":             status,
	}

	tst.transactions = append(tst.transactions, transactionInfo)
}

func (tst *TransactionStatusTracker) getTransactionStatuses() []map[string]interface{} {
	// Return the list of all transactions and their statuses
	return tst.transactions
}

func (tst *TransactionStatusTracker) getFailedTransactions() []map[string]interface{} {
	// Return the list of failed transactions
	var failedTransactions []map[string]interface{}
	for _, transaction := range tst.transactions {
		if transaction["status"].(string) == "failure" {
			failedTransactions = append(failedTransactions, transaction)
		}
	}
	return failedTransactions
}
```

## **Additional requirements**

1. **Update Broadcast Manager Struct**

To do this, we need to update the BroadcastManager struct to track failedTransactions.

1. **Update Broadcast method**

Update the broadcast method to include a retry method for failed transactions.

1. Update b**roadcastTransaction** method to **broadcastWithRetry** method.

This step entails looping through the range of maxRetries to perform the broadcast transaction and returning the result and then logging the failed transactions by appending a newly created FailedTransaction struct to the failedTransactions slice.

1. Create new method called **RetryFailedTransactions** method in ‘Broadcast Manager’ to retry all failed transactions.

```go
// BroadcastManager manages the process of broadcasting transactions
type BroadcastManager struct {
	blockchainRPCClient      *BlockchainRPCClient,
	transactionStatusTracker *TransactionStatusTracker,
	maxRetries:               int,
	failedTransactions       []*FailedTransaction

}

// FailedTransaction struct to store information about failed transactions
type FailedTransaction struct {
    signedTransaction string
    retryCount         int
}

func (bm *BroadcastManager) broadcast(signedTransaction string) string {
    // Try broadcasting the signed transaction with automatic retry on failure
    result := bm.broadcastWithRetry(signedTransaction)

    // Log transaction details and status
    bm.transactionStatusTracker.logTransaction(signedTransaction, result)

    return result
}

func (bm *BroadcastManager) broadcastWithRetry(signedTransaction string) string {
    // Implement retry logic for failed transactions

    for attempt := 1; attempt <= bm.maxRetries; attempt++ {
        result := bm.blockchainRPCClient.broadcastTransaction(signedTransaction)

        if result == "success" {
            return result
        }

        // Log the failed transaction
        bm.failedTransactions = append(bm.failedTransactions, &FailedTransaction{
            signedTransaction: signedTransaction,
            retryCount:         attempt,
        })

        // Sleep before the next retry, you may implement an exponential backoff strategy
        time.Sleep(2 * time.Second)
    }

    return "failure"
}

func (bm *BroadcastManager) RetryFailedTransactions() {
    for _, failedTransaction := range bm.failedTransactions {
        // Retry the failed transaction
        result := bm.blockchainRPCClient.broadcastWithRetry(failedTransaction.signedTransaction)

        // Update transaction status
        bm.transactionStatusTracker.logTransaction(failedTransaction.signedTransaction, result)

				// Check if the retry was successful
        if result == "success" {
            // The transaction was successfully retried, do not include it in the remainingFailedTransactions
            continue
        }

				// If the retry was not successful, keep the transaction in the list for potential future retries
        remainingFailedTransactions = append(remainingFailedTransactions, failedTransaction)
    }

		// Update the list of failed transactions to only include those that remain after retry
    bm.failedTransactions = remainingFailedTransactions
}
```

\***\*Handling Unexpected Service Restarts:\*\***

In order to handle such a situation, we will need to persist the record of all failed transactions to a database.

**Persist Failed Transactions to Database:**

- Integrate a database (e.g., MySQL, PostgreSQL) to store information about failed transactions.

**Modify FailedTransaction to Include Database ID:**

- Update the FailedTransaction struct to include an identifier that can be used for database storage.

```go
// FailedTransaction struct to store information about failed transactions
type FailedTransaction struct {
		id                int
    signedTransaction string
    retryCount         int

}
```

**Update broadcastWithRetry method:**

- Save failed transactions to the database.

```go
func (bm *BroadcastManager) broadcastWithRetry(signedTransaction string) string {
    // Implement retry logic for failed transactions

    for attempt := 1; attempt <= bm.maxRetries; attempt++ {
        result := bm.blockchainRPCClient.broadcastTransaction(signedTransaction)

        if result == "success" {
            return result
        }

				// Save the failed transaction to the database
        dbID := saveFailedTransactionToDatabase(signedTransaction, attempt)

        // Log the failed transaction
        bm.failedTransactions = append(bm.failedTransactions, &FailedTransaction{
            signedTransaction: signedTransaction,
            retryCount:         attempt,
						id:                 dbID
        })

        // Sleep before the next retry, you may implement an exponential backoff strategy
        time.Sleep(2 * time.Second)
    }

    return "failure"
}
```

**Implement** saveFailedTransactionToDatabase **Function:**

- Implement a function that saves information about failed transactions to the database.

```go
func saveFailedTransactionToDatabase(signedTransaction string, retryCount int) int {
    // Implement logic to save the failed transaction to the database
    // Return the database ID assigned to the failed transaction
    // ...
    return 123 // Replace with the actual database ID
}
```

\***\*Admin Retry Functionality:\*\***

**Add an Admin Retry Endpoint:**

- Introduce an admin endpoint, for example, **`POST /admin/retry-failed-transactions`**, to trigger the retry of all failed transactions.

```go
// Handler for admin retry endpoint
func (bm *BroadcastManager) AdminRetryFailedTransactions() {
    bm.RetryFailedTransactions()
}
```

**Expose the Admin Endpoint:**

- Expose the admin retry endpoint through your HTTP server.

```go
// Example HTTP server setup
func main() {
    bm := BroadcastManager(/*...*/)

    // Register the admin retry endpoint
    http.HandleFunc("/admin/retry-failed-transactions", bm.AdminRetryFailedTransactions)

    // ... Other HTTP handlers ...

    // Start the HTTP server
    http.ListenAndServe(":8080", nil)
}
```
