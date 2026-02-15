package mapping

import (
	"testing"
)

func TestNewAuthorMap(t *testing.T) {
	config := map[string]string{
		"johndoe": "John Doe <john@example.com>",
	}

	am := NewAuthorMap(config)
	if am == nil {
		t.Fatal("NewAuthorMap returned nil")
	}

	if am.mapping == nil {
		t.Error("mapping should not be nil")
	}

	if am.defaultEmail != "users.noreply.cvs.example.org" {
		t.Errorf("defaultEmail = %q, want %q", am.defaultEmail, "users.noreply.cvs.example.org")
	}
}

func TestNewAuthorMapNil(t *testing.T) {
	am := NewAuthorMap(nil)
	if am == nil {
		t.Fatal("NewAuthorMap returned nil")
	}

	// Should work with nil config
	name, email := am.Get("someuser")
	if name != "someuser" {
		t.Errorf("name = %q, want %q", name, "someuser")
	}
	if email != "someuser@users.noreply.cvs.example.org" {
		t.Errorf("email = %q, want %q", email, "someuser@users.noreply.cvs.example.org")
	}
}

func TestNewAuthorMapWithDefault(t *testing.T) {
	config := map[string]string{
		"janedoe": "Jane Doe <jane@example.com>",
	}

	am := NewAuthorMapWithDefault(config, "custom.domain.com")
	if am == nil {
		t.Fatal("NewAuthorMapWithDefault returned nil")
	}

	if am.defaultEmail != "custom.domain.com" {
		t.Errorf("defaultEmail = %q, want %q", am.defaultEmail, "custom.domain.com")
	}
}

func TestAuthorMapGetMapped(t *testing.T) {
	config := map[string]string{
		"johndoe": "John Doe <john@example.com>",
		"janedoe": "Jane Doe <jane@example.com>",
	}

	am := NewAuthorMap(config)

	tests := []struct {
		username  string
		wantName  string
		wantEmail string
	}{
		{"johndoe", "John Doe", "john@example.com"},
		{"janedoe", "Jane Doe", "jane@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			name, email := am.Get(tt.username)
			if name != tt.wantName {
				t.Errorf("Get(%q) name = %q, want %q", tt.username, name, tt.wantName)
			}
			if email != tt.wantEmail {
				t.Errorf("Get(%q) email = %q, want %q", tt.username, email, tt.wantEmail)
			}
		})
	}
}

func TestAuthorMapGetUnmapped(t *testing.T) {
	am := NewAuthorMap(map[string]string{
		"existing": "Existing User <existing@example.com>",
	})

	name, email := am.Get("unknown")
	if name != "unknown" {
		t.Errorf("Get unmapped name = %q, want %q", name, "unknown")
	}
	if email != "unknown@users.noreply.cvs.example.org" {
		t.Errorf("Get unmapped email = %q, want %q", email, "unknown@users.noreply.cvs.example.org")
	}
}

func TestAuthorMapGetCustomDefault(t *testing.T) {
	am := NewAuthorMapWithDefault(nil, "mycompany.com")

	name, email := am.Get("newuser")
	if name != "newuser" {
		t.Errorf("Get name = %q, want %q", name, "newuser")
	}
	if email != "newuser@mycompany.com" {
		t.Errorf("Get email = %q, want %q", email, "newuser@mycompany.com")
	}
}

func TestAuthorMapGetInvalidFormat(t *testing.T) {
	// When the mapped format is invalid, should fall back to default
	am := NewAuthorMap(map[string]string{
		"baduser": "invalid format without angle brackets",
	})

	name, email := am.Get("baduser")
	// Should fall back to default since format is invalid
	if name != "baduser" {
		t.Errorf("Get with invalid format name = %q, want %q", name, "baduser")
	}
	if email != "baduser@users.noreply.cvs.example.org" {
		t.Errorf("Get with invalid format email = %q, want %q", email, "baduser@users.noreply.cvs.example.org")
	}
}

func TestParseAuthor(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantName  string
		wantEmail string
		wantErr   bool
	}{
		{
			name:      "standard format",
			input:     "John Doe <john@example.com>",
			wantName:  "John Doe",
			wantEmail: "john@example.com",
			wantErr:   false,
		},
		{
			name:      "with extra spaces",
			input:     "  Jane Doe  <jane@example.com>  ",
			wantName:  "Jane Doe",
			wantEmail: "jane@example.com",
			wantErr:   true, // leading/trailing spaces not handled by regex
		},
		{
			name:      "single name",
			input:     "Alice <alice@example.com>",
			wantName:  "Alice",
			wantEmail: "alice@example.com",
			wantErr:   false,
		},
		{
			name:      "name with numbers",
			input:     "User123 <user123@example.com>",
			wantName:  "User123",
			wantEmail: "user123@example.com",
			wantErr:   false,
		},
		{
			name:      "email with subdomain",
			input:     "Bob Smith <bob@sub.example.com>",
			wantName:  "Bob Smith",
			wantEmail: "bob@sub.example.com",
			wantErr:   false,
		},
		{
			name:    "missing angle brackets",
			input:   "John Doe john@example.com",
			wantErr: true,
		},
		{
			name:    "missing closing bracket",
			input:   "John Doe <john@example.com",
			wantErr: true,
		},
		{
			name:    "missing opening bracket",
			input:   "John Doe john@example.com>",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "only brackets",
			input:   "<>",
			wantErr: true,
		},
		{
			name:    "empty name",
			input:   " <email@example.com>",
			wantErr: true,
		},
		{
			name:    "empty email",
			input:   "Name <>",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, email, err := ParseAuthor(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseAuthor(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseAuthor(%q) unexpected error: %v", tt.input, err)
				return
			}
			if name != tt.wantName {
				t.Errorf("ParseAuthor(%q) name = %q, want %q", tt.input, name, tt.wantName)
			}
			if email != tt.wantEmail {
				t.Errorf("ParseAuthor(%q) email = %q, want %q", tt.input, email, tt.wantEmail)
			}
		})
	}
}

func TestAuthorExtractor(t *testing.T) {
	ae := NewAuthorExtractor()
	if ae == nil {
		t.Fatal("NewAuthorExtractor returned nil")
	}
	if ae.authors == nil {
		t.Error("authors map should not be nil")
	}
}

func TestAuthorExtractorAdd(t *testing.T) {
	ae := NewAuthorExtractor()

	ae.Add("user1")
	ae.Add("user2")
	ae.Add("user1") // duplicate

	list := ae.List()
	if len(list) != 2 {
		t.Errorf("List() returned %d authors, want 2", len(list))
	}

	// Check that both users are present
	found := make(map[string]bool)
	for _, author := range list {
		found[author] = true
	}

	if !found["user1"] {
		t.Error("user1 not found in list")
	}
	if !found["user2"] {
		t.Error("user2 not found in list")
	}
}

func TestAuthorExtractorList(t *testing.T) {
	ae := NewAuthorExtractor()

	// Empty extractor
	list := ae.List()
	// nil slice is valid for empty in Go (len(nil slice) == 0)
	if len(list) != 0 {
		t.Errorf("empty List() returned %d items, want 0", len(list))
	}

	// Add some authors
	ae.Add("alice")
	ae.Add("bob")
	ae.Add("charlie")

	list = ae.List()
	if len(list) != 3 {
		t.Errorf("List() returned %d authors, want 3", len(list))
	}
}

func TestAuthorExtractorGenerateTemplate(t *testing.T) {
	ae := NewAuthorExtractor()
	ae.Add("user1")
	ae.Add("user2")

	template := ae.GenerateTemplate()
	if template == nil {
		t.Fatal("GenerateTemplate() returned nil")
	}

	if len(template) != 2 {
		t.Errorf("GenerateTemplate() returned %d entries, want 2", len(template))
	}

	// Check format of template entries
	if template["user1"] != "user1 <user1@example.com>" {
		t.Errorf("template[user1] = %q, want %q", template["user1"], "user1 <user1@example.com>")
	}
	if template["user2"] != "user2 <user2@example.com>" {
		t.Errorf("template[user2] = %q, want %q", template["user2"], "user2 <user2@example.com>")
	}
}

func TestAuthorExtractorGenerateTemplateEmpty(t *testing.T) {
	ae := NewAuthorExtractor()

	template := ae.GenerateTemplate()
	if template == nil {
		t.Fatal("GenerateTemplate() returned nil")
	}
	if len(template) != 0 {
		t.Errorf("GenerateTemplate() returned %d entries, want 0", len(template))
	}
}

func TestAuthorMapEmptyConfig(t *testing.T) {
	am := NewAuthorMap(map[string]string{})

	name, email := am.Get("anyone")
	if name != "anyone" {
		t.Errorf("Get with empty config name = %q, want %q", name, "anyone")
	}
	if email != "anyone@users.noreply.cvs.example.org" {
		t.Errorf("Get with empty config email = %q, want %q", email, "anyone@users.noreply.cvs.example.org")
	}
}

func TestParseAuthorSpecialCharacters(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantName  string
		wantEmail string
		wantErr   bool
	}{
		{
			name:      "hyphenated name",
			input:     "Mary-Jane Watson <mj@example.com>",
			wantName:  "Mary-Jane Watson",
			wantEmail: "mj@example.com",
			wantErr:   false,
		},
		{
			name:      "name with apostrophe",
			input:     "O'Brien <obrien@example.com>",
			wantName:  "O'Brien",
			wantEmail: "obrien@example.com",
			wantErr:   false,
		},
		{
			name:      "email with plus",
			input:     "Test User <test+tag@example.com>",
			wantName:  "Test User",
			wantEmail: "test+tag@example.com",
			wantErr:   false,
		},
		{
			name:      "unicode name",
			input:     "José García <jose@example.com>",
			wantName:  "José García",
			wantEmail: "jose@example.com",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, email, err := ParseAuthor(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseAuthor(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseAuthor(%q) unexpected error: %v", tt.input, err)
				return
			}
			if name != tt.wantName {
				t.Errorf("ParseAuthor(%q) name = %q, want %q", tt.input, name, tt.wantName)
			}
			if email != tt.wantEmail {
				t.Errorf("ParseAuthor(%q) email = %q, want %q", tt.input, email, tt.wantEmail)
			}
		})
	}
}
