# euromail-go

Official Go SDK for the [EuroMail](https://euromail.dev) transactional email service.

[![Go Reference](https://pkg.go.dev/badge/github.com/kalle-works/euromail-go.svg)](https://pkg.go.dev/github.com/kalle-works/euromail-go)

## Installation

```bash
go get github.com/kalle-works/euromail-go
```

Requires Go 1.21+. Zero external dependencies (stdlib only).

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "os"

    euromail "github.com/kalle-works/euromail-go"
)

func main() {
    client := euromail.NewClient(os.Getenv("EUROMAIL_API_KEY"))

    resp, err := client.SendEmail(context.Background(), euromail.SendEmailParams{
        From:     "sender@yourdomain.com",
        To:       euromail.ToRecipient("recipient@example.com"),
        Subject:  euromail.String("Hello from EuroMail"),
        HTMLBody: euromail.String("<h1>Welcome!</h1><p>Your account is ready.</p>"),
    })
    if err != nil {
        panic(err)
    }
    fmt.Printf("Email queued: %s\n", resp.ID)
}
```

## Configuration

```go
client := euromail.NewClient("em_live_...",
    euromail.WithTimeout(10 * time.Second),
    euromail.WithHTTPClient(customHTTPClient),
    euromail.WithBaseURL("https://custom-endpoint.example.com"),
)
```

## Sending Emails

### Direct send

```go
resp, err := client.SendEmail(ctx, euromail.SendEmailParams{
    From:     "noreply@yourdomain.com",
    To:       euromail.ToRecipient("user@example.com"),
    Subject:  euromail.String("Order Confirmation"),
    HTMLBody: euromail.String("<h1>Thanks for your order!</h1>"),
    TextBody: euromail.String("Thanks for your order!"),
    ReplyTo:  euromail.String("support@yourdomain.com"),
    Tags:     []string{"order", "confirmation"},
    Metadata: map[string]string{"order_id": "12345"},
})
```

### Multiple recipients

```go
resp, err := client.SendEmail(ctx, euromail.SendEmailParams{
    From:    "noreply@yourdomain.com",
    To:      euromail.Recipients("user1@example.com", "user2@example.com"),
    Subject: euromail.String("Team Update"),
    TextBody: euromail.String("Here's what happened this week."),
})
```

### Send with template

```go
resp, err := client.SendEmail(ctx, euromail.SendEmailParams{
    From:          "noreply@yourdomain.com",
    To:            euromail.ToRecipient("user@example.com"),
    TemplateAlias: euromail.String("welcome-email"),
    TemplateData: map[string]interface{}{
        "name":           "John",
        "activation_url": "https://example.com/activate/abc123",
    },
})
```

### Send with attachments

```go
resp, err := client.SendEmail(ctx, euromail.SendEmailParams{
    From:     "noreply@yourdomain.com",
    To:       euromail.ToRecipient("user@example.com"),
    Subject:  euromail.String("Your Invoice"),
    HTMLBody: euromail.String("<p>Please find your invoice attached.</p>"),
    Attachments: []euromail.Attachment{
        {
            Filename:    "invoice.pdf",
            Content:     base64EncodedContent,
            ContentType: "application/pdf",
        },
    },
})
```

### Batch send

```go
batch, err := client.SendBatch(ctx, euromail.SendBatchParams{
    Emails: []euromail.SendEmailParams{
        {
            From:     "noreply@yourdomain.com",
            To:       euromail.ToRecipient("user1@example.com"),
            Subject:  euromail.String("Hello User 1"),
            TextBody: euromail.String("Welcome!"),
        },
        {
            From:     "noreply@yourdomain.com",
            To:       euromail.ToRecipient("user2@example.com"),
            Subject:  euromail.String("Hello User 2"),
            TextBody: euromail.String("Welcome!"),
        },
    },
})
fmt.Printf("Sent: %d, Errors: %d\n", len(batch.Data), len(batch.Errors))
```

### Idempotent sends

```go
resp, err := client.SendEmail(ctx, euromail.SendEmailParams{
    From:           "noreply@yourdomain.com",
    To:             euromail.ToRecipient("user@example.com"),
    Subject:        euromail.String("Payment Receipt"),
    HTMLBody:        euromail.String("<p>Payment received.</p>"),
    IdempotencyKey: euromail.String("payment-receipt-12345"),
})
```

### Retrieve and list emails

```go
email, err := client.GetEmail(ctx, "email-uuid")

emails, pagination, err := client.ListEmails(ctx, &euromail.ListEmailsParams{
    ListParams: euromail.ListParams{
        Page:    euromail.Int(1),
        PerPage: euromail.Int(50),
    },
    Status: euromail.String("delivered"),
})
```

## Domains

```go
// Register a sending domain
domain, err := client.AddDomain(ctx, "mail.yourdomain.com")
fmt.Printf("DKIM selector: %s\n", domain.DKIMSelector)

// Trigger DNS verification
verification, err := client.VerifyDomain(ctx, domain.ID)
if spf, ok := verification.Checks["spf"]; ok {
    fmt.Printf("SPF verified: %v\n", spf.Verified)
}

// List all domains
domains, pagination, err := client.ListDomains(ctx, nil)
```

## Templates

```go
// Create a template
tmpl, err := client.CreateTemplate(ctx, euromail.CreateTemplateParams{
    Name:     "welcome",
    Alias:    euromail.String("welcome-email"),
    Subject:  "Welcome, {{ name }}!",
    HTMLBody: "<h1>Hello {{ name }}!</h1>",
})

// List templates
templates, pagination, err := client.ListTemplates(ctx, nil)
```

## Webhooks

```go
// Subscribe to delivery events
webhook, err := client.CreateWebhook(ctx, euromail.CreateWebhookParams{
    URL:    "https://example.com/webhooks/euromail",
    Events: []string{"email.delivered", "email.bounced", "email.complained"},
})

// Send a test event
client.TestWebhook(ctx, webhook.ID)
```

## Contact Lists

```go
// Create a list
list, err := client.CreateContactList(ctx, euromail.CreateContactListParams{
    Name: "Newsletter Subscribers",
})

// Add contacts
client.AddContact(ctx, list.ID, euromail.AddContactParams{
    Email: "user@example.com",
    Name:  euromail.String("Jane Doe"),
})

// Bulk add
client.BulkAddContacts(ctx, list.ID, euromail.BulkAddContactsParams{
    Contacts: []euromail.AddContactParams{
        {Email: "a@example.com"},
        {Email: "b@example.com"},
    },
})
```

## Suppressions

```go
client.AddSuppression(ctx, euromail.AddSuppressionParams{
    Email:  "bounce@example.com",
    Reason: euromail.String("hard_bounce"),
})

suppressions, pagination, err := client.ListSuppressions(ctx, nil)
```

## Analytics

```go
overview, err := client.GetAnalyticsOverview(ctx, &euromail.AnalyticsQuery{
    Period: euromail.String("30d"),
})

timeseries, err := client.GetAnalyticsTimeseries(ctx, &euromail.TimeseriesQuery{
    Period:  euromail.String("7d"),
    Metrics: euromail.String("sent,delivered,bounced"),
})

// Export as CSV
csv, err := client.ExportAnalyticsCSV(ctx, nil)
```

## Inbound Email

```go
inbound, pagination, err := client.ListInboundEmails(ctx, nil)

email, err := client.GetInboundEmail(ctx, "inbound-uuid")

// Set up routing rules
route, err := client.CreateInboundRoute(ctx, euromail.CreateInboundRouteParams{
    Pattern:  "*@yourdomain.com",
    Action:   "webhook",
    Endpoint: "https://example.com/inbound",
})
```

## Account

```go
account, err := client.GetAccount(ctx)
fmt.Printf("Plan: %s, Used: %d/%d\n",
    account.Plan, account.EmailsSentThisMonth, account.MonthlyQuota)

// Export account data (GDPR)
export, err := client.ExportAccount(ctx)

// Delete account permanently
client.DeleteAccount(ctx)
```

## Error Handling

All methods return `error`. Use type assertions or the `Is*Error` helpers:

```go
resp, err := client.SendEmail(ctx, params)
if err != nil {
    var authErr *euromail.AuthenticationError
    var valErr *euromail.ValidationError
    var rateErr *euromail.RateLimitError
    var notFound *euromail.NotFoundError

    switch {
    case errors.As(err, &authErr):
        log.Fatalf("Invalid API key: %s", authErr.Message)
    case errors.As(err, &valErr):
        log.Fatalf("Validation error [%s]: %s", valErr.Code, valErr.Message)
    case errors.As(err, &rateErr):
        log.Printf("Rate limited, retry after %d seconds", rateErr.RetryAfter)
        time.Sleep(time.Duration(rateErr.RetryAfter) * time.Second)
    case errors.As(err, &notFound):
        log.Fatalf("Not found: %s", notFound.Message)
    default:
        var apiErr *euromail.EuroMailError
        if errors.As(err, &apiErr) {
            log.Fatalf("API error [%d] %s: %s", apiErr.Status, apiErr.Code, apiErr.Message)
        }
        log.Fatalf("Network error: %v", err)
    }
}
```

Or use the convenience helpers:

```go
if euromail.IsRateLimitError(err) {
    // handle rate limiting
}
```

| Error Type | HTTP Status | Description |
|---|---|---|
| `AuthenticationError` | 401 | Invalid or missing API key |
| `ValidationError` | 422 | Invalid request parameters |
| `RateLimitError` | 429 | Too many requests (includes `RetryAfter`) |
| `NotFoundError` | 404 | Resource does not exist |
| `EuroMailError` | 4xx/5xx | Base type for all API errors |

## Agent Mailboxes

Agent mailboxes provide persistent email addresses for AI agents with at-least-once message delivery via a lease/ack/nack model. Native SDK support is coming in a future release. In the meantime, use `net/http` directly:

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "os"
)

const api = "https://api.euromail.dev"

func do(method, path string, body any) (*http.Response, error) {
    var buf bytes.Buffer
    if body != nil {
        if err := json.NewEncoder(&buf).Encode(body); err != nil {
            return nil, err
        }
    }
    req, err := http.NewRequest(method, api+path, &buf)
    if err != nil {
        return nil, err
    }
    req.Header.Set("X-EuroMail-Api-Key", os.Getenv("EUROMAIL_API_KEY"))
    req.Header.Set("Content-Type", "application/json")
    return http.DefaultClient.Do(req)
}

func run() error {
    // Create a mailbox
    resp, err := do("POST", "/v1/agent-mailboxes",
        map[string]string{"display_name": "Support Agent"})
    if err != nil {
        return err
    }
    var created struct {
        Data struct{ ID string } `json:"data"`
    }
    json.NewDecoder(resp.Body).Decode(&created)
    resp.Body.Close()

    for {
        // Long-poll for the next message (acquires a 5-minute lease)
        resp, err := do("GET",
            "/v1/agent-mailboxes/"+created.Data.ID+"/messages/next?timeout=30",
            nil)
        if err != nil {
            return err
        }
        if resp.StatusCode == http.StatusRequestTimeout {
            resp.Body.Close()
            continue
        }
        var body struct {
            Data       struct{ ID string } `json:"data"`
            LeaseToken string              `json:"lease_token"`
        }
        json.NewDecoder(resp.Body).Decode(&body)
        resp.Body.Close()

        if err := handle(body.Data); err != nil {
            // Nack to return the message to the queue for retry
            r, _ := do("POST",
                "/v1/agent-mailboxes/"+created.Data.ID+"/messages/"+body.Data.ID+"/nack",
                map[string]string{"lease_token": body.LeaseToken})
            r.Body.Close()
            continue
        }

        // Ack when done — message will not be redelivered
        r, _ := do("POST",
            "/v1/agent-mailboxes/"+created.Data.ID+"/messages/"+body.Data.ID+"/ack",
            map[string]string{"lease_token": body.LeaseToken})
        r.Body.Close()
    }
}
```

See the [Agent Mailboxes guide](https://euromail.dev/docs/guides/agent-mailboxes/) for the full flow, duplicate handling, and horizontal scaling patterns.

## API Reference

| Category | Method | Description |
|---|---|---|
| **Emails** | `SendEmail(ctx, params)` | Send a single email |
| | `SendBatch(ctx, params)` | Send up to 500 emails in one request |
| | `GetEmail(ctx, id)` | Get email details and delivery events |
| | `ListEmails(ctx, params)` | List emails with pagination and filters |
| | `CancelScheduledEmail(ctx, id)` | Cancel a scheduled email |
| | `GetEmailLinks(ctx, id)` | Get per-link click statistics |
| **Templates** | `CreateTemplate(ctx, params)` | Create an email template |
| | `GetTemplate(ctx, id)` | Get template by ID |
| | `UpdateTemplate(ctx, id, params)` | Update template fields |
| | `DeleteTemplate(ctx, id)` | Delete a template |
| | `ListTemplates(ctx, params)` | List templates with pagination |
| **Domains** | `AddDomain(ctx, domain)` | Register a sending domain |
| | `GetDomain(ctx, id)` | Get domain details and DNS records |
| | `VerifyDomain(ctx, id)` | Trigger DNS verification |
| | `DeleteDomain(ctx, id)` | Remove a domain |
| | `ListDomains(ctx, params)` | List domains with pagination |
| **Webhooks** | `CreateWebhook(ctx, params)` | Subscribe to events |
| | `GetWebhook(ctx, id)` | Get webhook details |
| | `UpdateWebhook(ctx, id, params)` | Update URL, events, or status |
| | `TestWebhook(ctx, id)` | Send a test event |
| | `DeleteWebhook(ctx, id)` | Remove a webhook |
| | `ListWebhooks(ctx, params)` | List webhooks with pagination |
| **Suppressions** | `AddSuppression(ctx, params)` | Suppress an email address |
| | `DeleteSuppression(ctx, email)` | Remove a suppression |
| | `ListSuppressions(ctx, params)` | List suppressions with pagination |
| **Contact Lists** | `CreateContactList(ctx, params)` | Create a contact list |
| | `GetContactList(ctx, id)` | Get list details |
| | `UpdateContactList(ctx, id, params)` | Update list settings |
| | `DeleteContactList(ctx, id)` | Delete a list |
| | `ListContactLists(ctx)` | List all contact lists |
| | `AddContact(ctx, listId, params)` | Add a contact to a list |
| | `BulkAddContacts(ctx, listId, params)` | Add multiple contacts |
| | `ListContacts(ctx, listId, params)` | List contacts with filters |
| | `RemoveContact(ctx, listId, email)` | Remove a contact |
| **Inbound** | `ListInboundEmails(ctx, params)` | List received emails |
| | `GetInboundEmail(ctx, id)` | Get inbound email details |
| | `DeleteInboundEmail(ctx, id)` | Delete an inbound email |
| **Inbound Routes** | `CreateInboundRoute(ctx, params)` | Create a routing rule |
| | `GetInboundRoute(ctx, id)` | Get route details |
| | `UpdateInboundRoute(ctx, id, params)` | Update a route |
| | `DeleteInboundRoute(ctx, id)` | Delete a route |
| | `ListInboundRoutes(ctx, params)` | List routes with pagination |
| **Analytics** | `GetAnalyticsOverview(ctx, query)` | Aggregated delivery stats |
| | `GetAnalyticsTimeseries(ctx, query)` | Daily metrics over time |
| | `GetAnalyticsDomains(ctx, query)` | Per-domain breakdown |
| | `ExportAnalyticsCSV(ctx, query)` | Export stats as CSV |
| **Audit Logs** | `ListAuditLogs(ctx, params)` | List account activity |
| **Dead Letters** | `ListDeadLetters(ctx, params)` | List permanently failed emails |
| | `RetryDeadLetter(ctx, id)` | Retry delivery |
| | `DeleteDeadLetter(ctx, id)` | Remove from dead letter queue |
| **Account** | `GetAccount(ctx)` | Get account info and quota |
| | `ExportAccount(ctx)` | Export all account data |
| | `DeleteAccount(ctx)` | Permanently delete account |

## Requirements

- Go 1.21+
- No external dependencies

## License

MIT
