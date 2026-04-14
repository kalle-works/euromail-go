// EuroMail Go SDK — comprehensive example exercising every method.
//
// Usage:
//
//	EUROMAIL_API_KEY=em_live_... go run ./examples/all_methods/
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	euromail "github.com/kalle-works/euromail-go"
)

func main() {
	client := euromail.NewClient(os.Getenv("EUROMAIL_API_KEY"))
	ctx := context.Background()

	// ---- Account ----
	account, err := client.GetAccount(ctx)
	must(err)
	fmt.Printf("Account: %s (%s)\n", account.Name, account.Plan)

	// ---- API Keys ----
	apiKey, err := client.CreateApiKey(ctx, euromail.CreateApiKeyParams{
		Name:   "test-key",
		Scopes: []string{"emails:send"},
	})
	must(err)
	fmt.Printf("Created API key: %s... (id: %s)\n", apiKey.KeyPrefix, apiKey.ID)

	keys, err := client.ListApiKeys(ctx)
	must(err)
	fmt.Printf("API keys: %d\n", len(keys))

	must(client.DeleteApiKey(ctx, apiKey.ID))
	fmt.Println("Deleted API key")

	// ---- Domains ----
	domain, err := client.AddDomain(ctx, "test-sdk-example.com")
	must(err)
	fmt.Printf("Added domain: %s (id: %s)\n", domain.Domain, domain.ID)

	domainDetail, err := client.GetDomain(ctx, domain.ID)
	must(err)
	fmt.Printf("Domain DKIM selector: %s\n", domainDetail.DKIMSelector)

	verification, err := client.VerifyDomain(ctx, domain.ID)
	must(err)
	spfCheck, ok := verification.Checks["spf"]
	if ok {
		fmt.Printf("Domain SPF verified: %v\n", spfCheck.Verified)
	}

	domains, _, err := client.ListDomains(ctx, &euromail.ListParams{Page: euromail.Int(1), PerPage: euromail.Int(10)})
	must(err)
	fmt.Printf("Domains: %d total\n", len(domains))

	// Tracking domain
	tracking, err := client.SetTrackingDomain(ctx, domain.ID, "track.test-sdk-example.com")
	if err == nil {
		fmt.Printf("Tracking domain CNAME target: %s\n", tracking.CnameTarget)
		tv, err := client.VerifyTrackingDomain(ctx, domain.ID)
		if err == nil {
			fmt.Printf("Tracking verified: %v\n", tv.TrackingCheck.Verified)
		} else {
			fmt.Printf("Tracking verify: %v\n", err)
		}
		_ = client.RemoveTrackingDomain(ctx, domain.ID)
		fmt.Println("Removed tracking domain")
	} else {
		fmt.Printf("Tracking domain: %v\n", err)
	}

	must(client.DeleteDomain(ctx, domain.ID))
	fmt.Println("Deleted domain")

	// ---- Templates ----
	templateAlias := fmt.Sprintf("test-welcome-%d", time.Now().Unix())
	template, err := client.CreateTemplate(ctx, euromail.CreateTemplateParams{
		Alias:    templateAlias,
		Name:     "Test Welcome",
		Subject:  "Welcome {{ name }}!",
		HTMLBody: euromail.String("<p>Hello {{ name }}</p>"),
	})
	must(err)
	fmt.Printf("Created template: %s (id: %s)\n", template.Alias, template.ID)

	tmpl, err := client.GetTemplate(ctx, template.ID)
	must(err)
	fmt.Printf("Template subject: %s\n", tmpl.Subject)

	updatedTmpl, err := client.UpdateTemplate(ctx, template.ID, euromail.UpdateTemplateParams{
		Name:     euromail.String("Updated Welcome"),
		Subject:  euromail.String(template.Subject),
		HTMLBody: euromail.String("<p>Updated {{ name }}</p>"),
	})
	must(err)
	fmt.Printf("Updated template name: %s\n", updatedTmpl.Name)

	templates, _, err := client.ListTemplates(ctx, &euromail.ListParams{Page: euromail.Int(1), PerPage: euromail.Int(10)})
	must(err)
	fmt.Printf("Templates: %d\n", len(templates))

	must(client.DeleteTemplate(ctx, template.ID))
	fmt.Println("Deleted template")

	// ---- Emails ----
	parts := strings.SplitN(account.Email, "@", 2)
	fromDomain := "example.com"
	if len(parts) == 2 {
		fromDomain = parts[1]
	}
	sent, err := client.SendEmail(ctx, euromail.SendEmailParams{
		From:     fmt.Sprintf("test@%s", fromDomain),
		To:       euromail.ToRecipient(account.Email),
		Subject:  euromail.String("SDK test"),
		TextBody: euromail.String("Hello from the Go SDK example!"),
	})
	must(err)
	fmt.Printf("Sent email: %s (status: %s)\n", sent.ID, sent.Status)

	emailDetail, err := client.GetEmail(ctx, sent.ID)
	must(err)
	fmt.Printf("Email to: %s\n", emailDetail.Email.ToAddress)

	emails, _, err := client.ListEmails(ctx, &euromail.ListEmailsParams{Page: euromail.Int(1), PerPage: euromail.Int(5)})
	must(err)
	fmt.Printf("Emails: %d\n", len(emails))

	// ---- Email Validation ----
	validation, err := client.ValidateEmail(ctx, "test@example.com")
	must(err)
	fmt.Printf("Validation: valid=%v, deliverable=%s\n", validation.Valid, validation.Deliverable)

	// ---- Webhooks ----
	webhook, err := client.CreateWebhook(ctx, euromail.CreateWebhookParams{
		URL:    "https://httpbin.org/post",
		Events: []string{"delivered", "bounced"},
	})
	must(err)
	fmt.Printf("Created webhook: %s\n", webhook.ID)

	wh, err := client.GetWebhook(ctx, webhook.ID)
	must(err)
	fmt.Printf("Webhook events: %s\n", strings.Join(wh.Events, ", "))

	updatedWh, err := client.UpdateWebhook(ctx, webhook.ID, euromail.UpdateWebhookParams{
		URL:      "https://httpbin.org/post",
		Events:   []string{"delivered", "bounced", "opened"},
		IsActive: true,
	})
	must(err)
	fmt.Printf("Updated webhook events: %s\n", strings.Join(updatedWh.Events, ", "))

	webhooks, _, err := client.ListWebhooks(ctx, &euromail.ListParams{Page: euromail.Int(1), PerPage: euromail.Int(10)})
	must(err)
	fmt.Printf("Webhooks: %d\n", len(webhooks))

	testResp, err := client.TestWebhook(ctx, webhook.ID)
	if err == nil {
		fmt.Printf("Webhook test: %s\n", testResp.Message)
	} else {
		fmt.Printf("Webhook test: %v\n", err)
	}

	must(client.DeleteWebhook(ctx, webhook.ID))
	fmt.Println("Deleted webhook")

	// ---- Suppressions ----
	suppression, err := client.AddSuppression(ctx, "blocked@example.com", "manual")
	must(err)
	fmt.Printf("Added suppression: %s\n", suppression.EmailAddress)

	suppressions, _, err := client.ListSuppressions(ctx, &euromail.ListParams{Page: euromail.Int(1), PerPage: euromail.Int(10)})
	must(err)
	fmt.Printf("Suppressions: %d\n", len(suppressions))

	must(client.DeleteSuppression(ctx, "blocked@example.com"))
	fmt.Println("Deleted suppression")

	// ---- Contact Lists ----
	contactList, err := client.CreateContactList(ctx, euromail.CreateContactListParams{Name: "SDK Test List"})
	must(err)
	fmt.Printf("Created list: %s (id: %s)\n", contactList.Name, contactList.ID)

	contact, err := client.AddContact(ctx, contactList.ID, euromail.AddContactParams{Email: "user@example.com"})
	must(err)
	fmt.Printf("Added contact: %s\n", contact.Email)

	bulk, err := client.BulkAddContacts(ctx, contactList.ID, euromail.BulkAddContactsParams{
		Contacts: []euromail.BulkContactEntry{
			{Email: "a@example.com"},
			{Email: "b@example.com"},
		},
	})
	must(err)
	fmt.Printf("Bulk added: %d/%d\n", bulk.Inserted, bulk.TotalRequested)

	contacts, _, err := client.ListContacts(ctx, contactList.ID, &euromail.ListContactsParams{Page: euromail.Int(1), PerPage: euromail.Int(10)})
	must(err)
	fmt.Printf("Contacts: %d\n", len(contacts))

	must(client.RemoveContact(ctx, contactList.ID, "user@example.com"))
	fmt.Println("Removed contact")

	lists, err := client.ListContactLists(ctx)
	must(err)
	fmt.Printf("Contact lists: %d\n", len(lists))

	// ---- Signup Forms ----
	signupForm, err := client.CreateSignupForm(ctx, euromail.CreateSignupFormParams{
		ListID: contactList.ID,
		Title:  "SDK Test Signup",
	})
	must(err)
	fmt.Printf("Created signup form: %s (slug: %s)\n", signupForm.Title, signupForm.Slug)

	sf, err := client.GetSignupForm(ctx, signupForm.ID)
	must(err)
	fmt.Printf("Signup form URL: %s\n", sf.FormURL)

	updatedSf, err := client.UpdateSignupForm(ctx, signupForm.ID, euromail.UpdateSignupFormParams{
		Title: "Updated Signup",
	})
	must(err)
	fmt.Printf("Updated signup form: %s\n", updatedSf.Title)

	toggled, err := client.ToggleSignupForm(ctx, signupForm.ID)
	must(err)
	fmt.Printf("Toggled signup form active: %v\n", toggled.IsActive)

	forms, err := client.ListSignupForms(ctx)
	must(err)
	fmt.Printf("Signup forms: %d\n", len(forms))

	must(client.DeleteSignupForm(ctx, signupForm.ID))
	fmt.Println("Deleted signup form")

	must(client.DeleteContactList(ctx, contactList.ID))
	fmt.Println("Deleted contact list")

	// ---- Analytics ----
	overview, err := client.GetAnalyticsOverview(ctx, &euromail.AnalyticsQuery{Period: euromail.String("30d")})
	must(err)
	fmt.Printf("Analytics: %d sent, %d delivered\n", overview.Data.Sent, overview.Data.Delivered)

	ts, err := client.GetAnalyticsTimeseries(ctx, &euromail.TimeseriesQuery{AnalyticsQuery: euromail.AnalyticsQuery{Period: euromail.String("7d")}})
	must(err)
	fmt.Printf("Timeseries points: %d\n", len(ts.Data))

	domainStats, err := client.GetAnalyticsDomains(ctx, &euromail.DomainAnalyticsQuery{AnalyticsQuery: euromail.AnalyticsQuery{Period: euromail.String("30d")}, Limit: euromail.Int(5)})
	must(err)
	fmt.Printf("Domain analytics: %d domains\n", len(domainStats.Data))

	csv, err := client.ExportAnalyticsCSV(ctx, &euromail.AnalyticsQuery{Period: euromail.String("7d")})
	must(err)
	fmt.Printf("CSV export: %d bytes\n", len(csv))

	// ---- Operations ----
	ops, _, err := client.ListOperations(ctx, &euromail.ListParams{Page: euromail.Int(1), PerPage: euromail.Int(5)})
	must(err)
	fmt.Printf("Operations: %d\n", len(ops))

	// ---- Audit Logs ----
	logs, _, err := client.ListAuditLogs(ctx, &euromail.ListParams{Page: euromail.Int(1), PerPage: euromail.Int(5)})
	must(err)
	fmt.Printf("Audit logs: %d\n", len(logs))

	// ---- Dead Letters ----
	deadLetters, err := client.ListDeadLetters(ctx, euromail.Int(5))
	must(err)
	fmt.Printf("Dead letters: %d\n", len(deadLetters.Data))

	// ---- Inbound ----
	inbound, _, err := client.ListInboundEmails(ctx, &euromail.ListParams{Page: euromail.Int(1), PerPage: euromail.Int(5)})
	must(err)
	fmt.Printf("Inbound emails: %d\n", len(inbound))

	routes, _, err := client.ListInboundRoutes(ctx, &euromail.ListParams{Page: euromail.Int(1), PerPage: euromail.Int(5)})
	must(err)
	fmt.Printf("Inbound routes: %d\n", len(routes))

	// ---- Billing ----
	plans, err := client.ListPlans(ctx)
	must(err)
	names := make([]string, len(plans))
	for i, p := range plans {
		names[i] = p.Plan
	}
	fmt.Printf("Plans: %s\n", strings.Join(names, ", "))

	sub, err := client.GetSubscription(ctx)
	must(err)
	fmt.Printf("Subscription: %s (%s)\n", sub.Plan, sub.SubscriptionStatus)

	// ---- GDPR ----
	gdprExport, err := client.GdprExport(ctx, "test@example.com")
	if err == nil {
		fmt.Printf("GDPR export: %s\n", gdprExport.Data.EmailAddress)
	} else {
		fmt.Printf("GDPR export: %v\n", err)
	}

	// ---- Agent Mailboxes ----
	mailboxes, _, err := client.ListMailboxes(ctx, &euromail.ListParams{Page: euromail.Int(1), PerPage: euromail.Int(5)})
	must(err)
	fmt.Printf("Agent mailboxes: %d\n", len(mailboxes))

	mb, err := client.CreateMailbox(ctx, euromail.CreateMailboxParams{
		DisplayName: euromail.String("Example Agent"),
	})
	if err != nil {
		fmt.Printf("CreateMailbox skipped: %v\n", err)
	} else {
		fmt.Printf("Created mailbox: %s\n", mb.Address)

		// Short long-poll — returns nil on timeout.
		leased, err := client.WaitForNextMessage(ctx, mb.ID, euromail.Int(1))
		must(err)
		if leased == nil {
			fmt.Println("No messages waiting (408 timeout — expected for a fresh mailbox)")
		} else {
			fmt.Printf("Got message %s, acking\n", leased.Data.ID)
			must(client.AckMessage(ctx, mb.ID, leased.Data.ID, leased.LeaseToken))
		}

		must(client.DeleteMailbox(ctx, mb.ID))
		fmt.Println("Mailbox deleted")
	}

	fmt.Println("\nAll methods exercised successfully!")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
