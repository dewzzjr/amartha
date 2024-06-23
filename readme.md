# Amartha: reconciliation service (Algorithmic and scaling)

Design and implement a transaction reconciliation service that identifies unmatched and discrepant transactions between internal data (system transactions) and external data (bank statements) for Amartha.

## Problem Statement

Amartha manages multiple bank accounts and requires a service to reconcile transactions occurring within their system against corresponding transactions reflected in bank statements. This process helps identify errors, discrepancies, and missing transactions.

## Data Model

### Transaction

- `trxID` : Unique identifier for the transaction (string)
- `amount` : Transaction amount (decimal)
- `type` : Transaction type (enum: DEBIT, CREDIT)
- `transactionTime` : Date and time of the transaction (datetime)

### Bank Statement

- `unique_identifier` : Unique identifier for the transaction in the bank statement (string) (varies by bank, not necessarily equivalent to trxID)
- `amount` : Transaction amount (decimal) (can be negative for debits)
- `date` : Date of the transaction (date)

## Assumptions

- Both system transactions and bank statements are provided as separate CSV files.
- Discrepancies only occur in amount.

## Functionality

- The service accepts the following input parameters:
  - System transaction CSV file path
  - Bank statement CSV file path (can handle multiple files from different banks)
  - Start date for reconciliation timeframe (date)
  - End date for reconciliation timeframe (date)
- The service performs the reconciliation process by comparing transactions within the specified timeframe across system and bank statement data.
- The service outputs a reconciliation summary containing:
  - Total number of transactions processed
  - Total number of matched transactions
  - Total number of unmatched transactions
    - Details of unmatched transactions:
      - System transaction details if missing in bank statement(s)
      - Bank statement details if missing in system transactions (grouped by bank)
  - Total discrepancies (sum of absolute differences in amount between matched transactions)
