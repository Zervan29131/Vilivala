package handler

import "testing"

func TestValidateSiteConfig(t *testing.T) {
	cfg := DefaultSiteConfig()
	if err := ValidateSiteConfig(&cfg); err != nil {
		t.Fatalf("expected default site config to be valid, got %v", err)
	}

	badTitle := DefaultSiteConfig()
	badTitle.SiteTitle = "   "
	if err := ValidateSiteConfig(&badTitle); err == nil {
		t.Fatal("expected empty site_title to fail validation")
	}

	badURL := DefaultSiteConfig()
	badURL.Socials = []SocialLink{{Title: "XSS", Link: "javascript:alert(1)"}}
	if err := ValidateSiteConfig(&badURL); err == nil {
		t.Fatal("expected javascript: URL in socials to fail validation")
	}

	badMode := DefaultSiteConfig()
	badMode.NoticeMode = "unknown"
	if err := ValidateSiteConfig(&badMode); err == nil {
		t.Fatal("expected unknown notice_mode to fail validation")
	}
}
