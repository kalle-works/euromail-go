# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Native agent mailbox support: `CreateMailbox`, `ListMailboxes`, `GetMailbox`,
  `DeleteMailbox`, `ListMailboxMessages`, `DeleteMailboxMessage`.
- Long-poll lease/ack/nack methods: `WaitForNextMessage`, `AckMessage`,
  `NackMessage`. `WaitForNextMessage` returns `(nil, nil)` on HTTP 408 so
  polling loops can simply continue on timeout.
- Types: `AgentMailbox`, `MailboxMessage`, `LeasedMessage`,
  `CreateMailboxParams`, `ListMailboxMessagesParams`.
- Agent mailbox parity methods: `ReplyToMessage`, `ListMailboxThreads`,
  `GetMailboxThread`, `SearchMailboxMessages`, `UpdateMessageLabels`,
  `GetMessageAttachmentURLs`, `ListMailboxContacts`, `GetMailboxAnalytics`,
  and `UpdateAutoResponder`.
- Types: `ReplyToMessageParams`, `MailboxReplyResult`, `MailboxAttachmentURL`,
  `MailboxContact`, `MailboxAnalytics`, `UpdateAutoResponderParams`,
  `AutoResponderConfig`. `MailboxMessage` gains threading (`InReplyTo`,
  `ReferencesHeader`), attachment (`AttachmentsStored`, `AttachmentsMetadata`),
  classification (`Classification`, `ClassificationConfidence`, `ClassifiedAt`),
  and lease (`LeasedUntil`, `LeaseToken`) fields; `AgentMailbox` gains
  `WebhookFilters`, `AutoResponderEnabled`, and `AutoResponderRules`.
- README rewritten with native SDK examples in place of the raw `net/http`
  snippet.
- `SendAt`, `Stream`, `Tracking`, and `Transactional` fields on
  `SendEmailParams` (and `Tracking`/`Transactional` on `BroadcastParams`).
  `Transactional` defaults to `true` when left `nil` on `SendEmail`/`SendBatch`
  (unlike `SendBroadcast`, which defaults to the account setting — a broadcast
  targets a contact list and is typically marketing mail).
- `VerifyWebhookSignature` verifies the `X-Euromail-Signature` header on a
  webhook delivery (`t=<unix_ts>,v1=<hex_hmac>` over `"<ts>.<raw_body>"`,
  constant-time compared, 5-minute default tolerance) and returns the
  timestamp it was signed at. Sentinel errors
  (`ErrWebhookSignatureMissingTimestamp`, `ErrWebhookSignatureMissingSignature`,
  `ErrWebhookSignatureExpired`, `ErrWebhookSignatureInvalid`) let a caller
  distinguish failure reasons with `errors.Is`.
- `ImportSuppressions` / `ExportSuppressions`: bulk-import up to 10,000
  addresses in one call and export the full suppression list as CSV. New
  `ImportSuppressionsParams` / `ImportSuppressionsResult` types.

### Changed

- `EuroMailError` gained `Type` and `DocsURL` fields (from the API's error
  envelope) and a `RequestID` (from the `X-Request-Id` response header).
  Validation failures are now classified by error code/type rather than
  HTTP status alone, since the API can return a validation failure as either
  `400` or `422` depending on the endpoint.

## [0.1.0] - 2026-04-13

### Added

- Initial Go SDK for the euromail transactional email API.
- `NewClient` constructor with functional options (`WithBaseURL`, `WithTimeout`, `WithHTTPClient`).
- Email sending, listing, and retrieval (`SendEmail`, `ListEmails`, `GetEmail`, `GetEmailLinks`).
- Domain management (list, get, add, verify, delete).
- Contact list and subscriber management.
- Newsletter sending and listing.
- Template management (create, list, get, update, delete).
- API key management.
- Sub-account support.
- Webhook management.
- Inbound email and routing support.
- Suppression list management.
- Dead letter queue access.
- Analytics and insights (`GetAnalytics`, `GenerateInsights`).
- GDPR data export and deletion.
- Signup form management.
- Billing and account information.
- Audit log access.
- Email validation.
- Comprehensive README with usage examples.
