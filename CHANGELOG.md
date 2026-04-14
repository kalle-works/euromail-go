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
- README rewritten with native SDK examples in place of the raw `net/http`
  snippet.

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
