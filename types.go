package euromail

import "encoding/json"

// ---- Pagination ----

// Pagination holds pagination metadata returned by list endpoints.
type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ListParams are common query parameters for paginated list endpoints.
type ListParams struct {
	Page    *int `json:"page,omitempty"`
	PerPage *int `json:"per_page,omitempty"`
}

// ---- Account ----

// Account represents the authenticated account.
type Account struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Email               string `json:"email"`
	Plan                string `json:"plan"`
	MonthlyQuota        int    `json:"monthly_quota"`
	EmailsSentThisMonth int64  `json:"emails_sent_this_month"`
	QuotaResetAt        string `json:"quota_reset_at"`
	CreatedAt           string `json:"created_at"`
}

// AccountDeleteResponse is returned when an account deletion is accepted.
type AccountDeleteResponse struct {
	Message   string `json:"message"`
	AccountID string `json:"account_id"`
}

// ---- Email ----

// Attachment represents a file attachment on an email.
type Attachment struct {
	Filename    string `json:"filename"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
}

// SendEmailParams are the parameters for sending a single email.
// Recipient accepts a single string or a slice of strings for the "to" field.
// Use euromail.String("user@example.com") for one, or euromail.Recipients("a@example.com", "b@example.com") for many.
type Recipient struct{ v interface{} }

// Recipients creates a Recipient for multiple addresses.
func Recipients(addrs ...string) Recipient {
	if len(addrs) == 1 {
		return Recipient{v: addrs[0]}
	}
	return Recipient{v: addrs}
}

// ToRecipient creates a single-address Recipient (convenience, same as Recipients with one arg).
func ToRecipient(addr string) Recipient { return Recipient{v: addr} }

func (r Recipient) MarshalJSON() ([]byte, error) {
	if r.v == nil {
		return json.Marshal("")
	}
	return json.Marshal(r.v)
}

type SendEmailParams struct {
	From           string                 `json:"from"`
	To             Recipient              `json:"to"`
	Subject        *string                `json:"subject,omitempty"`
	CC             []string               `json:"cc,omitempty"`
	BCC            []string               `json:"bcc,omitempty"`
	ReplyTo        *string                `json:"reply_to,omitempty"`
	HTMLBody       *string                `json:"html_body,omitempty"`
	TextBody       *string                `json:"text_body,omitempty"`
	TemplateAlias  *string                `json:"template_alias,omitempty"`
	TemplateData   map[string]interface{} `json:"template_data,omitempty"`
	Headers        map[string]string      `json:"headers,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	Metadata       map[string]string      `json:"metadata,omitempty"`
	IdempotencyKey *string                `json:"idempotency_key,omitempty"`
	Attachments    []Attachment           `json:"attachments,omitempty"`
}

// SendEmailResponse is returned after successfully sending an email.
type SendEmailResponse struct {
	ID          string  `json:"id"`
	MessageID   string  `json:"message_id"`
	Status      string  `json:"status"`
	To          string  `json:"to"`
	Sandbox     bool    `json:"sandbox"`
	ScheduledAt *string `json:"scheduled_at"`
	CreatedAt   string  `json:"created_at"`
}

// SendBatchParams are the parameters for sending a batch of emails.
type SendBatchParams struct {
	Emails []SendEmailParams `json:"emails"`
}

// BatchError represents an error for a single email in a batch send.
type BatchError struct {
	Index int    `json:"index"`
	Error string `json:"error"`
}

// SendBatchResponse is returned after a batch send.
type SendBatchResponse struct {
	Data   []SendEmailResponse `json:"data"`
	Errors []BatchError        `json:"errors"`
}

// EmailDetail wraps an email with its delivery events (returned by GetEmail).
type EmailDetail struct {
	Email  Email        `json:"email"`
	Events []EmailEvent `json:"events"`
}

// EmailEvent represents a delivery lifecycle event.
type EmailEvent struct {
	ID             string                 `json:"id"`
	EmailID        string                 `json:"email_id"`
	AccountID      string                 `json:"account_id"`
	EventType      string                 `json:"event_type"`
	BounceType     *string                `json:"bounce_type"`
	BounceCategory *string                `json:"bounce_category"`
	RemoteMTA      *string                `json:"remote_mta"`
	DiagnosticCode *string                `json:"diagnostic_code"`
	UserAgent      *string                `json:"user_agent"`
	IPAddress      *string                `json:"ip_address"`
	LinkURL        *string                `json:"link_url"`
	RawPayload     map[string]interface{} `json:"raw_payload"`
	CreatedAt      string                 `json:"created_at"`
}

// Email represents a full email record.
type Email struct {
	ID           string                 `json:"id"`
	AccountID    string                 `json:"account_id"`
	DomainID     *string                `json:"domain_id"`
	MessageID    string                 `json:"message_id"`
	FromAddress  string                 `json:"from_address"`
	ToAddress    string                 `json:"to_address"`
	CC           []string               `json:"cc"`
	BCC          []string               `json:"bcc"`
	ReplyTo      *string                `json:"reply_to"`
	Subject      string                 `json:"subject"`
	HTMLBody     *string                `json:"html_body"`
	TextBody     *string                `json:"text_body"`
	TemplateID   *string                `json:"template_id"`
	TemplateData map[string]interface{} `json:"template_data"`
	Headers      map[string]string      `json:"headers"`
	Tags         []string               `json:"tags"`
	Metadata     map[string]string      `json:"metadata"`
	Status       string                 `json:"status"`
	Attempts     int                    `json:"attempts"`
	MaxAttempts  int                    `json:"max_attempts"`
	ErrorMessage *string                `json:"error_message"`
	SMTPResponse *string                `json:"smtp_response"`
	CreatedAt    string                 `json:"created_at"`
	UpdatedAt    string                 `json:"updated_at"`
	SentAt       *string                `json:"sent_at"`
}

// ListEmailsParams are query parameters for listing emails.
type ListEmailsParams struct {
	Page    *int    `json:"page,omitempty"`
	PerPage *int    `json:"per_page,omitempty"`
	Status  *string `json:"status,omitempty"`
}

// ---- Domain ----

// DnsRecord represents a single DNS record for domain verification.
type DnsRecord struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Value    string `json:"value"`
	Priority *int   `json:"priority,omitempty"`
}

// Domain represents a sending domain.
type Domain struct {
	ID                       string                `json:"id"`
	AccountID                string                `json:"account_id"`
	Domain                   string                `json:"domain"`
	DKIMSelector             string                `json:"dkim_selector"`
	DKIMPublicKey            string                `json:"dkim_public_key"`
	SPFVerified              bool                  `json:"spf_verified"`
	DKIMVerified             bool                  `json:"dkim_verified"`
	DMARCVerified            bool                  `json:"dmarc_verified"`
	ReturnPathVerified       bool                  `json:"return_path_verified"`
	MXVerified               bool                  `json:"mx_verified"`
	InboundEnabled           bool                  `json:"inbound_enabled"`
	MXVerifiedAt             *string               `json:"mx_verified_at"`
	DnsRecords               map[string]DnsRecord  `json:"dns_records"`
	VerifiedAt               *string               `json:"verified_at"`
	TrackingDomain           *string               `json:"tracking_domain"`
	TrackingDomainVerified   bool                  `json:"tracking_domain_verified"`
	TrackingDomainVerifiedAt *string               `json:"tracking_domain_verified_at"`
	CreatedAt                string                `json:"created_at"`
	UpdatedAt                string                `json:"updated_at"`
}

// VerificationCheck holds the result of a single DNS verification check.
type VerificationCheck struct {
	Verified bool   `json:"verified"`
	Detail   string `json:"detail"`
}

// DomainVerificationResult holds the result of a domain verification, including the updated domain and individual checks.
type DomainVerificationResult struct {
	Domain Domain                        `json:"domain"`
	Checks map[string]VerificationCheck  `json:"checks"`
}

// ---- Template ----

// Template represents an email template.
type Template struct {
	ID        string  `json:"id"`
	AccountID string  `json:"account_id"`
	Alias     string  `json:"alias"`
	Name      string  `json:"name"`
	Subject   string  `json:"subject"`
	HTMLBody  *string `json:"html_body"`
	TextBody  *string `json:"text_body"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// CreateTemplateParams are the parameters for creating a template.
type CreateTemplateParams struct {
	Alias    string  `json:"alias"`
	Name     string  `json:"name"`
	Subject  string  `json:"subject"`
	HTMLBody *string `json:"html_body,omitempty"`
	TextBody *string `json:"text_body,omitempty"`
}

// UpdateTemplateParams are the parameters for updating a template.
type UpdateTemplateParams struct {
	Name     *string `json:"name,omitempty"`
	Subject  *string `json:"subject,omitempty"`
	HTMLBody *string `json:"html_body,omitempty"`
	TextBody *string `json:"text_body,omitempty"`
}

// ---- Webhook ----

// Webhook represents a webhook endpoint.
type Webhook struct {
	ID                string   `json:"id"`
	AccountID         string   `json:"account_id"`
	URL               string   `json:"url"`
	Events            []string `json:"events"`
	IsActive          bool     `json:"is_active"`
	Secret            *string  `json:"secret,omitempty"`
	FailureCount      *int     `json:"failure_count,omitempty"`
	LastSuccessAt     *string  `json:"last_success_at,omitempty"`
	LastFailureAt     *string  `json:"last_failure_at,omitempty"`
	LastFailureReason *string  `json:"last_failure_reason,omitempty"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
}

// CreateWebhookParams are the parameters for creating a webhook.
type CreateWebhookParams struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

// UpdateWebhookParams are the parameters for updating a webhook.
type UpdateWebhookParams struct {
	URL      string   `json:"url"`
	Events   []string `json:"events"`
	IsActive bool     `json:"is_active"`
}

// WebhookTestResponse is returned after testing a webhook.
type WebhookTestResponse struct {
	Message string                 `json:"message"`
	Payload map[string]interface{} `json:"payload"`
}

// ---- Suppression ----

// Suppression represents a suppressed email address.
type Suppression struct {
	ID            string  `json:"id"`
	AccountID     string  `json:"account_id"`
	EmailAddress  string  `json:"email_address"`
	Reason        string  `json:"reason"`
	SourceEmailID *string `json:"source_email_id"`
	CreatedAt     string  `json:"created_at"`
}

// ---- Contact List ----

// ContactList represents a mailing list.
type ContactList struct {
	ID           string  `json:"id"`
	AccountID    string  `json:"account_id"`
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	DoubleOptIn  bool    `json:"double_opt_in"`
	ContactCount int     `json:"contact_count"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// CreateContactListParams are the parameters for creating a contact list.
type CreateContactListParams struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	DoubleOptIn *bool   `json:"double_opt_in,omitempty"`
}

// UpdateContactListParams are the parameters for updating a contact list.
type UpdateContactListParams struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	DoubleOptIn bool    `json:"double_opt_in"`
}

// Contact represents a contact in a contact list.
type Contact struct {
	ID        string            `json:"id"`
	ListID    string            `json:"list_id"`
	Email     string            `json:"email"`
	Metadata  map[string]string `json:"metadata"`
	Status    string            `json:"status"`
	CreatedAt string            `json:"created_at"`
}

// AddContactParams are the parameters for adding a contact to a list.
type AddContactParams struct {
	Email    string            `json:"email"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// BulkContactEntry represents a single contact in a bulk add operation.
type BulkContactEntry struct {
	Email    string            `json:"email"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// BulkAddContactsParams are the parameters for bulk adding contacts.
type BulkAddContactsParams struct {
	Contacts []BulkContactEntry `json:"contacts"`
}

// BulkAddContactsResponse is returned after a bulk add operation.
type BulkAddContactsResponse struct {
	Inserted       int `json:"inserted"`
	TotalRequested int `json:"total_requested"`
}

// ListContactsParams are query parameters for listing contacts.
type ListContactsParams struct {
	Page    *int    `json:"page,omitempty"`
	PerPage *int    `json:"per_page,omitempty"`
	Status  *string `json:"status,omitempty"`
}

// ---- Analytics ----

// AnalyticsQuery holds query parameters for analytics endpoints.
type AnalyticsQuery struct {
	Period *string `json:"period,omitempty"` // "7d", "30d", "90d"
	From   *string `json:"from,omitempty"`
	To     *string `json:"to,omitempty"`
}

// TimeseriesQuery extends AnalyticsQuery with a metrics filter.
type TimeseriesQuery struct {
	AnalyticsQuery
	Metrics *string `json:"metrics,omitempty"`
}

// DomainAnalyticsQuery extends AnalyticsQuery with a limit.
type DomainAnalyticsQuery struct {
	AnalyticsQuery
	Limit *int `json:"limit,omitempty"`
}

// AnalyticsPeriod describes the time range for analytics data.
type AnalyticsPeriod struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Period string `json:"period"`
}

// AnalyticsData holds the summary metrics.
type AnalyticsData struct {
	Sent         int     `json:"sent"`
	Delivered    int     `json:"delivered"`
	Bounced      int     `json:"bounced"`
	Opens        int     `json:"opens"`
	Clicks       int     `json:"clicks"`
	Complaints   int     `json:"complaints"`
	DeliveryRate float64 `json:"delivery_rate"`
	OpenRate     float64 `json:"open_rate"`
	ClickRate    float64 `json:"click_rate"`
}

// AnalyticsSummary is the response from the analytics overview endpoint.
type AnalyticsSummary struct {
	Data   AnalyticsData   `json:"data"`
	Period AnalyticsPeriod `json:"period"`
}

// TimeseriesPoint represents a single data point in a time series.
type TimeseriesPoint struct {
	Date      string `json:"date"`
	Sent      *int   `json:"sent,omitempty"`
	Delivered *int   `json:"delivered,omitempty"`
	Bounced   *int   `json:"bounced,omitempty"`
	Opens     *int   `json:"opens,omitempty"`
	Clicks    *int   `json:"clicks,omitempty"`
}

// TimeseriesResponse is the response from the analytics timeseries endpoint.
type TimeseriesResponse struct {
	Data   []TimeseriesPoint `json:"data"`
	Period AnalyticsPeriod   `json:"period"`
}

// DomainAnalytics holds analytics data for a single domain.
type DomainAnalytics struct {
	Domain    string      `json:"domain"`
	Sent      int         `json:"sent"`
	Delivered int         `json:"delivered"`
	Bounced   int         `json:"bounced"`
	OpenRate  FlexFloat64 `json:"open_rate"`
	ClickRate FlexFloat64 `json:"click_rate"`
}

// DomainAnalyticsResponse is the response from the analytics domains endpoint.
type DomainAnalyticsResponse struct {
	Data   []DomainAnalytics `json:"data"`
	Period AnalyticsPeriod   `json:"period"`
}

// ---- Audit Log ----

// AuditLog represents an audit log entry.
type AuditLog struct {
	ID           string                 `json:"id"`
	AccountID    string                 `json:"account_id"`
	Action       string                 `json:"action"`
	ResourceType string                 `json:"resource_type"`
	ResourceID   *string                `json:"resource_id"`
	IPAddress    *string                `json:"ip_address"`
	UserAgent    *string                `json:"user_agent"`
	Details      map[string]interface{} `json:"details"`
	CreatedAt    string                 `json:"created_at"`
}

// ---- Dead Letter ----

// DeadLetter represents a dead letter message.
type DeadLetter struct {
	StreamID       string                 `json:"stream_id"`
	OriginalStream string                 `json:"original_stream"`
	EmailID        string                 `json:"email_id"`
	AccountID      string                 `json:"account_id"`
	FailureReason  string                 `json:"failure_reason"`
	AttemptCount   int                    `json:"attempt_count"`
	LastError      string                 `json:"last_error"`
	FailedAt       string                 `json:"failed_at"`
	Payload        map[string]interface{} `json:"payload"`
}

// DeadLetterListResponse is the response from listing dead letters.
type DeadLetterListResponse struct {
	Data  []DeadLetter `json:"data"`
	Total int          `json:"total"`
}

// MessageResponse is a generic response containing a message string.
type MessageResponse struct {
	Message string `json:"message"`
}

// ---- Inbound Email ----

// InboundEmail represents an inbound email.
type InboundEmail struct {
	ID                 string                 `json:"id"`
	AccountID          string                 `json:"account_id"`
	DomainID           string                 `json:"domain_id"`
	RouteID            *string                `json:"route_id"`
	MessageID          *string                `json:"message_id"`
	MailFrom           string                 `json:"mail_from"`
	RcptTo             []string               `json:"rcpt_to"`
	FromHeader         *string                `json:"from_header"`
	ToHeader           *string                `json:"to_header"`
	CCHeader           *string                `json:"cc_header"`
	Subject            *string                `json:"subject"`
	TextBody           *string                `json:"text_body"`
	HTMLBody           *string                `json:"html_body"`
	RawHeaders         map[string]interface{} `json:"raw_headers"`
	Attachments        interface{}            `json:"attachments"`
	SourceIP           *string                `json:"source_ip"`
	SPFResult          *string                `json:"spf_result"`
	SizeBytes          int                    `json:"size_bytes"`
	Status             string                 `json:"status"`
	WebhookDeliveredAt *string                `json:"webhook_delivered_at"`
	ErrorMessage       *string                `json:"error_message"`
	CreatedAt          string                 `json:"created_at"`
	UpdatedAt          string                 `json:"updated_at"`
}

// ---- Inbound Route ----

// InboundRoute represents an inbound routing rule.
type InboundRoute struct {
	ID         string  `json:"id"`
	AccountID  string  `json:"account_id"`
	DomainID   string  `json:"domain_id"`
	Pattern    string  `json:"pattern"`
	MatchType  string  `json:"match_type"`
	Priority   int     `json:"priority"`
	WebhookURL *string `json:"webhook_url"`
	IsActive   bool    `json:"is_active"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

// CreateInboundRouteParams are the parameters for creating an inbound route.
type CreateInboundRouteParams struct {
	DomainID   string  `json:"domain_id"`
	Pattern    string  `json:"pattern"`
	MatchType  string  `json:"match_type"`
	Priority   *int    `json:"priority,omitempty"`
	WebhookURL *string `json:"webhook_url,omitempty"`
}

// UpdateInboundRouteParams are the parameters for updating an inbound route.
type UpdateInboundRouteParams struct {
	Pattern    string  `json:"pattern"`
	MatchType  string  `json:"match_type"`
	Priority   int     `json:"priority"`
	WebhookURL *string `json:"webhook_url,omitempty"`
	IsActive   bool    `json:"is_active"`
}

// ---- API Key ----

// ApiKey represents an API key (without the secret portion).
type ApiKey struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	KeyPrefix  string   `json:"key_prefix"`
	Scopes     []string `json:"scopes"`
	IsActive   bool     `json:"is_active"`
	LastUsedAt *string  `json:"last_used_at"`
	CreatedAt  string   `json:"created_at"`
}

// ApiKeyCreated is returned when creating an API key, includes the full key.
type ApiKeyCreated struct {
	ApiKey
	Key string `json:"key"`
}

// CreateApiKeyParams are the parameters for creating an API key.
type CreateApiKeyParams struct {
	Name   string   `json:"name"`
	Scopes []string `json:"scopes,omitempty"`
}

// ---- Newsletter ----

// Newsletter represents a newsletter draft or sent newsletter.
type Newsletter struct {
	ID              string                 `json:"id"`
	AccountID       string                 `json:"account_id"`
	ListID          *string                `json:"list_id"`
	Subject         string                 `json:"subject"`
	FromAddress     string                 `json:"from_address"`
	HTMLBody        *string                `json:"html_body"`
	TextBody        *string                `json:"text_body"`
	TemplateID      *string                `json:"template_id"`
	TemplateData    map[string]interface{} `json:"template_data"`
	ReplyTo         *string                `json:"reply_to"`
	Status          string                 `json:"status"`
	OperationID     *string                `json:"operation_id"`
	ScheduledAt     *string                `json:"scheduled_at"`
	SentAt          *string                `json:"sent_at"`
	TotalRecipients *int                   `json:"total_recipients"`
	CreatedAt       string                 `json:"created_at"`
	UpdatedAt       string                 `json:"updated_at"`
}

// CreateNewsletterParams are the parameters for creating a newsletter.
type CreateNewsletterParams struct {
	ListID       string                 `json:"list_id"`
	Subject      string                 `json:"subject"`
	FromAddress  string                 `json:"from_address"`
	HTMLBody     *string                `json:"html_body,omitempty"`
	TextBody     *string                `json:"text_body,omitempty"`
	TemplateID   *string                `json:"template_id,omitempty"`
	TemplateData map[string]interface{} `json:"template_data,omitempty"`
	ReplyTo      *string                `json:"reply_to,omitempty"`
}

// UpdateNewsletterParams are the parameters for updating a newsletter.
type UpdateNewsletterParams struct {
	ListID       *string                `json:"list_id,omitempty"`
	Subject      *string                `json:"subject,omitempty"`
	FromAddress  *string                `json:"from_address,omitempty"`
	HTMLBody     *string                `json:"html_body,omitempty"`
	TextBody     *string                `json:"text_body,omitempty"`
	TemplateID   *string                `json:"template_id,omitempty"`
	TemplateData map[string]interface{} `json:"template_data,omitempty"`
	ReplyTo      *string                `json:"reply_to,omitempty"`
}

// NewsletterSendResponse is returned after sending a newsletter.
type NewsletterSendResponse struct {
	OperationID     string `json:"operation_id"`
	TotalRecipients int    `json:"total_recipients"`
	Message         string `json:"message"`
}

// ---- Broadcast ----

// BroadcastParams are the parameters for broadcasting to a contact list.
type BroadcastParams struct {
	ContactListID string                 `json:"contact_list_id"`
	FromAddress   string                 `json:"from_address"`
	Subject       *string                `json:"subject,omitempty"`
	HTMLBody      *string                `json:"html_body,omitempty"`
	TextBody      *string                `json:"text_body,omitempty"`
	TemplateAlias *string                `json:"template_alias,omitempty"`
	TemplateData  map[string]interface{} `json:"template_data,omitempty"`
	ReplyTo       *string                `json:"reply_to,omitempty"`
	Headers       map[string]string      `json:"headers,omitempty"`
	Tags          []string               `json:"tags,omitempty"`
	SendAt        *string                `json:"send_at,omitempty"`
}

// BroadcastResponse is returned after a broadcast send.
type BroadcastResponse struct {
	OperationID     string `json:"operation_id"`
	TotalRecipients int    `json:"total_recipients"`
	Message         string `json:"message"`
}

// ---- Signup Form ----

// SignupForm represents a signup form for collecting contacts.
type SignupForm struct {
	ID             string          `json:"id"`
	AccountID      string          `json:"account_id"`
	ListID         string          `json:"list_id"`
	Slug           string          `json:"slug"`
	Title          string          `json:"title"`
	Description    *string         `json:"description"`
	SuccessMessage *string         `json:"success_message"`
	RedirectURL    *string         `json:"redirect_url"`
	CustomFields   json.RawMessage `json:"custom_fields"`
	Theme          json.RawMessage `json:"theme"`
	IsActive       bool            `json:"is_active"`
	FormURL        string          `json:"form_url"`
	EmbedCode      string          `json:"embed_code"`
	CreatedAt      string          `json:"created_at"`
	UpdatedAt      string          `json:"updated_at"`
}

// CreateSignupFormParams are the parameters for creating a signup form.
type CreateSignupFormParams struct {
	ListID         string          `json:"list_id"`
	Title          string          `json:"title"`
	Description    *string         `json:"description,omitempty"`
	SuccessMessage *string         `json:"success_message,omitempty"`
	RedirectURL    *string         `json:"redirect_url,omitempty"`
	CustomFields   json.RawMessage `json:"custom_fields,omitempty"`
	Theme          json.RawMessage `json:"theme,omitempty"`
}

// UpdateSignupFormParams are the parameters for updating a signup form.
type UpdateSignupFormParams struct {
	Title          string          `json:"title"`
	Description    *string         `json:"description,omitempty"`
	SuccessMessage *string         `json:"success_message,omitempty"`
	RedirectURL    *string         `json:"redirect_url,omitempty"`
	CustomFields   json.RawMessage `json:"custom_fields,omitempty"`
	Theme          json.RawMessage `json:"theme,omitempty"`
}

// ---- Email Validation ----

// EmailValidation is the result of validating an email address.
type EmailValidation struct {
	Email        string  `json:"email"`
	Valid        bool    `json:"valid"`
	Deliverable  string  `json:"deliverable"`
	IsDisposable bool    `json:"is_disposable"`
	IsRole       bool    `json:"is_role"`
	IsFree       bool    `json:"is_free"`
	MXFound      bool    `json:"mx_found"`
	Reason       *string `json:"reason"`
}

// ---- Operation ----

// Operation represents an async operation.
type Operation struct {
	ID             string                 `json:"id"`
	AccountID      string                 `json:"account_id"`
	OperationType  string                 `json:"operation_type"`
	Status         string                 `json:"status"`
	TotalItems     int                    `json:"total_items"`
	CompletedItems int                    `json:"completed_items"`
	FailedItems    int                    `json:"failed_items"`
	ErrorSummary   map[string]interface{} `json:"error_summary"`
	Metadata       map[string]interface{} `json:"metadata"`
	CreatedAt      string                 `json:"created_at"`
	UpdatedAt      string                 `json:"updated_at"`
	CompletedAt    *string                `json:"completed_at"`
	ExpiresAt      string                 `json:"expires_at"`
}

// ---- Billing ----

// BillingPlan represents a billing plan.
type BillingPlan struct {
	Plan            string  `json:"plan"`
	MonthlyQuota    int     `json:"monthly_quota"`
	MaxDomains      int     `json:"max_domains"`
	MaxTemplates    int     `json:"max_templates"`
	MaxWebhooks     int     `json:"max_webhooks"`
	MaxContactLists int     `json:"max_contact_lists"`
	MaxSubAccounts  int     `json:"max_sub_accounts"`
	TrackingEnabled bool    `json:"tracking_enabled"`
	PriceCents      int     `json:"price_cents"`
	StripePriceID   *string `json:"stripe_price_id"`
}

// SubscriptionLimits are the limits for a subscription plan.
type SubscriptionLimits struct {
	MaxDomains      int  `json:"max_domains"`
	MaxTemplates    int  `json:"max_templates"`
	MaxWebhooks     int  `json:"max_webhooks"`
	TrackingEnabled bool `json:"tracking_enabled"`
	PriceCents      int  `json:"price_cents"`
}

// Subscription represents the current billing subscription.
type Subscription struct {
	Plan                 string             `json:"plan"`
	SubscriptionStatus   string             `json:"subscription_status"`
	StripeSubscriptionID *string            `json:"stripe_subscription_id"`
	BillingEmail         *string            `json:"billing_email"`
	TrialEndsAt          *string            `json:"trial_ends_at"`
	MonthlyQuota         int                `json:"monthly_quota"`
	EmailsSentThisMonth  int64              `json:"emails_sent_this_month"`
	Limits               SubscriptionLimits `json:"limits"`
}

// CheckoutParams are the parameters for creating a Stripe checkout session.
type CheckoutParams struct {
	Plan       string `json:"plan"`
	SuccessURL string `json:"success_url"`
	CancelURL  string `json:"cancel_url"`
}

// CheckoutResponse contains the Stripe checkout URL.
type CheckoutResponse struct {
	CheckoutURL string `json:"checkout_url"`
}

// PortalParams are the parameters for creating a Stripe billing portal session.
type PortalParams struct {
	ReturnURL string `json:"return_url"`
}

// PortalResponse contains the Stripe billing portal URL.
type PortalResponse struct {
	PortalURL string `json:"portal_url"`
}

// ---- GDPR ----

// GdprExportResponse is the response from a GDPR data export.
type GdprExportResponse struct {
	Data       GdprExportData `json:"data"`
	ExportedAt string         `json:"exported_at"`
}

// GdprExportData contains all personal data for an email address.
type GdprExportData struct {
	EmailAddress      string                   `json:"email_address"`
	Emails            []map[string]interface{} `json:"emails"`
	Events            []map[string]interface{} `json:"events"`
	Suppressions      []map[string]interface{} `json:"suppressions"`
	UnsubscribeEvents []map[string]interface{} `json:"unsubscribe_events"`
	InboundEmails     []map[string]interface{} `json:"inbound_emails"`
}

// GdprEraseResponse is the response from a GDPR data erasure.
type GdprEraseResponse struct {
	Data        GdprEraseData `json:"data"`
	OperationID string        `json:"operation_id"`
}

// GdprEraseData contains the result of erasing data for an email address.
type GdprEraseData struct {
	EmailAddress string `json:"email_address"`
	RowsDeleted  int64  `json:"rows_deleted"`
	Message      string `json:"message"`
}

// ---- Tracking Domain ----

// TrackingDomainResponse is returned when setting a tracking domain.
type TrackingDomainResponse struct {
	Data        Domain `json:"data"`
	CnameTarget string `json:"cname_target"`
}

// TrackingCheck holds the result of a tracking domain verification.
type TrackingCheck struct {
	Verified bool   `json:"verified"`
	Detail   string `json:"detail"`
}

// TrackingDomainVerification is returned when verifying a tracking domain.
type TrackingDomainVerification struct {
	Data          Domain        `json:"data"`
	TrackingCheck TrackingCheck `json:"tracking_check"`
}
