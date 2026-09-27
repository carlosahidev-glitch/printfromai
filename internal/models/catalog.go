package models

// QuantityPrice holds the pricing for a specific quantity.
type QuantityPrice struct {
	Qty       int     `json:"qty"`
	FrontOnly float64 `json:"frontOnly"` // 4/0 pricing
	FrontBack float64 `json:"frontBack"` // 4/4 pricing (Set to 0 if N/A)
}

// SizeOption represents a physical dimension with an optional preview image.
type SizeOption struct {
	Label           string `json:"label"`
	Value           string `json:"value"`
	PreviewImageURL string `json:"previewImageUrl"` // For visual size comparisons
}

// OptionChoice represents a single selectable choice for a variant option.
type OptionChoice struct {
	Label     string          `json:"label"`
	Value     string          `json:"value"`
	Upcharges map[int]float64 `json:"upcharges"` // Quantity -> Upcharge mapping
}

// VariantOption represents a dynamic finishing option (e.g., "Coating", "Corner Cut").
type VariantOption struct {
	Title   string         `json:"title"`
	Choices []OptionChoice `json:"choices"`
}

// Variant represents a specific material or stock.
type Variant struct {
	Name             string                     `json:"name"`
	ShortDescription string                     `json:"shortDescription"`
	Features         []string                   `json:"features"`
	PreviewImageURL  string                     `json:"previewImageUrl"`
	ImageURLs        []string                   `json:"imageUrls"`
	Options          []VariantOption            `json:"options"`
	Pricing          map[string][]QuantityPrice `json:"pricing"` // Size Value -> Pricing Array
}

// FAQ holds a single question and answer pair for GEO and user help.
type FAQ struct {
	Question string
	Answer   string
}

// Review represents a customer review for social proof.
type Review struct {
	Author string
	Rating int // out of 5
	Date   string
	Text   string
}

// Product contains all the SEO, specification, pricing, and social proof data needed
type Product struct {
	Title           string
	MetaTitle       string
	MetaDescription string
	Slug            string
	Description     string
	LongDescription string
	PrimaryImageURL string
	JSONSchema      string
	SEOKeywords     []string
	Features        []string
	TurnaroundDays  string
	VariantLabel    string       // e.g., "Material" or "Stock"
	Sizes           []SizeOption // Upgraded to struct for previews
	ReviewScore     float64
	ReviewCount     int
	CustomerImages  []string
	CustomerReviews []Review
	Variants        []Variant
	FAQs            []FAQ
}

// ProductCatalog is our in-memory "database" compiled directly into the Go binary.
var ProductCatalog = map[string]Product{
	"business-cards": {
		Title:           "Business Cards",
		MetaTitle:       "Premium AI-Ready Business Card Printing | PrintFromAI",
		MetaDescription: "Print your AI-generated business card designs on premium 16pt, suede, and plastic stocks. Fast turnaround and 300 DPI CMYK conversion included.",
		Slug:            "business-cards",
		Description:     "Make a lasting impression with press-ready, ultra-high-quality cardstock tailored perfectly for your AI-generated art.",
		LongDescription: "Our business cards are explicitly engineered for the AI era. When you generate a design using Midjourney, ChatGPT, or Claude, it often lacks the strict color profiles (CMYK) and DPI requirements needed for commercial printing. We bridge that gap. Every design uploaded is automatically upscaled to a crystal-clear 300 DPI and mapped to standard bleed dimensions, ensuring your physical cards look exactly as vibrant and sharp as your digital prompt. Choose from industry-leading 16pt standard stock, luxurious suede textures, or futuristic plastic finishes.",
		PrimaryImageURL: "https://printfromai.com/static/images/og-business-cards.jpg",
		JSONSchema:      `{"@context": "https://schema.org/", "@type": "Product", "name": "Premium AI-Ready Business Cards", "description": "High-quality business card printing specifically optimized for AI-generated designs.", "brand": {"@type": "Brand", "name": "PrintFromAI"}}`,
		SEOKeywords:     []string{"AI business cards", "custom print", "premium cardstock", "suede business cards", "plastic business cards", "print AI art"},
		Features: []string{
			"Free automated 300 DPI upscaling",
			"Automatic RGB to CMYK press-ready conversion",
			"Standard 2x3.5\" sizing with full bleed",
			"Optional $25 human-expert prepress check",
		},
		TurnaroundDays: "2-4",
		VariantLabel:   "Material",
		Sizes: []SizeOption{
			{Label: "2 x 3.5\" (Standard)", Value: "2x3.5"},
			{Label: "2 x 2\" (Square)", Value: "2x2"},
			{Label: "1.75 x 3.5\" (Slim)", Value: "1.75x3.5"},
		},
		ReviewScore: 4.9,
		ReviewCount: 1304,
		CustomerImages: []string{
			"https://placehold.co/600x800/1A1A1A/EDEDED?text=Customer+Card+1",
			"https://placehold.co/800x600/1A1A1A/EDEDED?text=Customer+Card+2",
			"https://placehold.co/600x600/1A1A1A/EDEDED?text=Customer+Card+3",
			"https://placehold.co/800x1000/1A1A1A/EDEDED?text=Customer+Card+4",
			"https://placehold.co/800x600/1A1A1A/EDEDED?text=Customer+Card+5",
		},
		CustomerReviews: []Review{
			{Author: "Sarah Jenkins", Rating: 5, Date: "October 12, 2026", Text: "I generated a crazy cyberpunk design in Midjourney and was worried it would print dark. PrintFromAI nailed the colors. The Suede finish feels incredibly premium."},
			{Author: "David Chen", Rating: 5, Date: "September 28, 2026", Text: "Fastest turnaround ever. The 16pt stock is thick and the automated CMYK conversion saved me a huge headache. Highly recommend."},
			{Author: "Elena R.", Rating: 4, Date: "September 15, 2026", Text: "Great quality cards. The frosted plastic looks so futuristic with my minimalist AI logo. Only docked a star because I wish they had overnight shipping."},
			{Author: "Marcus T.", Rating: 5, Date: "August 30, 2026", Text: "Flawless execution. Uploaded my PDF, saw the crop marks instantly, and the final printed result was exactly as shown on screen."},
		},
		Variants: []Variant{
			{
				Name:             "Standard (16pt)",
				ShortDescription: "Our classic, rigid 16pt stock provides a heavy, premium feel. Available in standard Matte or High-Gloss UV coating.",
				Features:         []string{"16pt thickness", "Choice of Matte or Glossy finish", "Cost-effective premium feel"},
				PreviewImageURL:  "/static/images/variants/bc-standard-preview.webp",
				ImageURLs:        []string{"/static/images/variants/bc-standard-1.webp", "/static/images/variants/bc-standard-2.webp"},
				Options: []VariantOption{
					{
						Title: "Coating Options",
						Choices: []OptionChoice{
							{Label: "UV Gloss (Both Sides)", Value: "gloss-both", Upcharges: map[int]float64{}},
							{Label: "Gloss Front / Matte Back", Value: "gloss-front-matte-back", Upcharges: map[int]float64{}},
							{Label: "Matte (Both Sides)", Value: "matte-both", Upcharges: map[int]float64{}},
						},
					},
					{
						Title: "Corner Cut",
						Choices: []OptionChoice{
							{Label: "Square", Value: "square", Upcharges: map[int]float64{}},
							{Label: "Rounded 1/8\"", Value: "round-1-8", Upcharges: map[int]float64{100: 5.0, 250: 10.0, 500: 15.0, 1000: 25.0, 2500: 45.0, 5000: 80.0}},
							{Label: "Rounded 1/4\"", Value: "round-1-4", Upcharges: map[int]float64{100: 5.0, 250: 10.0, 500: 15.0, 1000: 25.0, 2500: 45.0, 5000: 80.0}},
						},
					},
				},
				Pricing: map[string][]QuantityPrice{
					"2x3.5": {
						{Qty: 100, FrontOnly: 15.00, FrontBack: 20.00},
						{Qty: 250, FrontOnly: 20.00, FrontBack: 30.00},
						{Qty: 500, FrontOnly: 25.00, FrontBack: 40.00},
						{Qty: 1000, FrontOnly: 40.00, FrontBack: 65.00},
						{Qty: 2500, FrontOnly: 70.00, FrontBack: 110.00},
						{Qty: 5000, FrontOnly: 120.00, FrontBack: 190.00},
					},
					"2x2": {
						{Qty: 100, FrontOnly: 15.00, FrontBack: 20.00},
						{Qty: 250, FrontOnly: 20.00, FrontBack: 30.00},
						{Qty: 500, FrontOnly: 25.00, FrontBack: 40.00},
						{Qty: 1000, FrontOnly: 40.00, FrontBack: 65.00},
						{Qty: 2500, FrontOnly: 70.00, FrontBack: 110.00},
						{Qty: 5000, FrontOnly: 120.00, FrontBack: 190.00},
					},
					"1.75x3.5": {
						{Qty: 100, FrontOnly: 15.00, FrontBack: 20.00},
						{Qty: 250, FrontOnly: 20.00, FrontBack: 30.00},
						{Qty: 500, FrontOnly: 25.00, FrontBack: 40.00},
						{Qty: 1000, FrontOnly: 40.00, FrontBack: 65.00},
						{Qty: 2500, FrontOnly: 70.00, FrontBack: 110.00},
						{Qty: 5000, FrontOnly: 120.00, FrontBack: 190.00},
					},
				},
			},
			{
				Name:             "Premium Suede",
				ShortDescription: "Luxurious soft-touch velvet lamination. This finish adds a tactile, high-end suede feel that elevates digital art.",
				Features:         []string{"Soft-touch velvet feel", "Scuff resistant", "Enhances dark colors beautifully"},
				PreviewImageURL:  "/static/images/variants/bc-suede-preview.webp",
				ImageURLs:        []string{"/static/images/variants/bc-suede-1.webp"},
				Options: []VariantOption{
					{
						Title: "Corner Cut",
						Choices: []OptionChoice{
							{Label: "Square", Value: "square", Upcharges: map[int]float64{}},
							{Label: "Rounded 1/8\"", Value: "round-1-8", Upcharges: map[int]float64{100: 5.0, 250: 10.0, 500: 15.0, 1000: 25.0, 2500: 45.0, 5000: 80.0}},
							{Label: "Rounded 1/4\"", Value: "round-1-4", Upcharges: map[int]float64{100: 5.0, 250: 10.0, 500: 15.0, 1000: 25.0, 2500: 45.0, 5000: 80.0}},
						},
					},
				},
				Pricing: map[string][]QuantityPrice{
					"2x3.5": {
						{Qty: 100, FrontOnly: 35.00, FrontBack: 45.00},
						{Qty: 250, FrontOnly: 45.00, FrontBack: 60.00},
						{Qty: 500, FrontOnly: 65.00, FrontBack: 80.00},
						{Qty: 1000, FrontOnly: 95.00, FrontBack: 120.00},
						{Qty: 2500, FrontOnly: 190.00, FrontBack: 240.00},
						{Qty: 5000, FrontOnly: 350.00, FrontBack: 450.00},
					},
					"2x2": {
						{Qty: 100, FrontOnly: 35.00, FrontBack: 45.00},
						{Qty: 250, FrontOnly: 45.00, FrontBack: 60.00},
						{Qty: 500, FrontOnly: 65.00, FrontBack: 80.00},
						{Qty: 1000, FrontOnly: 95.00, FrontBack: 120.00},
						{Qty: 2500, FrontOnly: 190.00, FrontBack: 240.00},
						{Qty: 5000, FrontOnly: 350.00, FrontBack: 450.00},
					},
					"1.75x3.5": {
						{Qty: 100, FrontOnly: 35.00, FrontBack: 45.00},
						{Qty: 250, FrontOnly: 45.00, FrontBack: 60.00},
						{Qty: 500, FrontOnly: 65.00, FrontBack: 80.00},
						{Qty: 1000, FrontOnly: 95.00, FrontBack: 120.00},
						{Qty: 2500, FrontOnly: 190.00, FrontBack: 240.00},
						{Qty: 5000, FrontOnly: 350.00, FrontBack: 450.00},
					},
				},
			},
			{
				Name:             "Smooth Silk",
				ShortDescription: "Elegant silky texture that mutes harsh reflections while keeping colors vibrant. Water and tear resistant.",
				Features:         []string{"Silky matte lamination", "Water and tear resistant", "Smooth, professional finish"},
				PreviewImageURL:  "/static/images/variants/bc-silk-preview.webp",
				ImageURLs:        []string{"/static/images/variants/bc-silk-1.webp"},
				Options: []VariantOption{
					{
						Title: "Corner Cut",
						Choices: []OptionChoice{
							{Label: "Square", Value: "square", Upcharges: map[int]float64{}},
							{Label: "Rounded 1/8\"", Value: "round-1-8", Upcharges: map[int]float64{100: 5.0, 250: 10.0, 500: 15.0, 1000: 25.0, 2500: 45.0, 5000: 80.0}},
							{Label: "Rounded 1/4\"", Value: "round-1-4", Upcharges: map[int]float64{100: 5.0, 250: 10.0, 500: 15.0, 1000: 25.0, 2500: 45.0, 5000: 80.0}},
						},
					},
				},
				Pricing: map[string][]QuantityPrice{
					"2x3.5": {
						{Qty: 100, FrontOnly: 30.00, FrontBack: 40.00},
						{Qty: 250, FrontOnly: 40.00, FrontBack: 55.00},
						{Qty: 500, FrontOnly: 60.00, FrontBack: 75.00},
						{Qty: 1000, FrontOnly: 90.00, FrontBack: 110.00},
						{Qty: 2500, FrontOnly: 180.00, FrontBack: 220.00},
						{Qty: 5000, FrontOnly: 320.00, FrontBack: 400.00},
					},
					"2x2": {
						{Qty: 100, FrontOnly: 30.00, FrontBack: 40.00},
						{Qty: 250, FrontOnly: 40.00, FrontBack: 55.00},
						{Qty: 500, FrontOnly: 60.00, FrontBack: 75.00},
						{Qty: 1000, FrontOnly: 90.00, FrontBack: 110.00},
						{Qty: 2500, FrontOnly: 180.00, FrontBack: 220.00},
						{Qty: 5000, FrontOnly: 320.00, FrontBack: 400.00},
					},
					"1.75x3.5": {
						{Qty: 100, FrontOnly: 30.00, FrontBack: 40.00},
						{Qty: 250, FrontOnly: 40.00, FrontBack: 55.00},
						{Qty: 500, FrontOnly: 60.00, FrontBack: 75.00},
						{Qty: 1000, FrontOnly: 90.00, FrontBack: 110.00},
						{Qty: 2500, FrontOnly: 180.00, FrontBack: 220.00},
						{Qty: 5000, FrontOnly: 320.00, FrontBack: 400.00},
					},
				},
			},
			{
				Name:             "Solid Plastic",
				ShortDescription: "Ultra-durable, solid opaque white plastic. 100% waterproof and impossible to tear. Perfect for VIP cards.",
				Features:         []string{"20pt plastic thickness", "100% waterproof", "Tear-proof opaque white base"},
				PreviewImageURL:  "/static/images/variants/bc-solid-plastic-preview.webp",
				ImageURLs:        []string{"/static/images/variants/bc-solid-plastic-1.webp"},
				Pricing: map[string][]QuantityPrice{
					"2x3.5": {
						{Qty: 100, FrontOnly: 50.00, FrontBack: 65.00},
						{Qty: 250, FrontOnly: 75.00, FrontBack: 95.00},
						{Qty: 500, FrontOnly: 110.00, FrontBack: 140.00},
						{Qty: 1000, FrontOnly: 180.00, FrontBack: 220.00},
						{Qty: 2500, FrontOnly: 400.00, FrontBack: 480.00},
						{Qty: 5000, FrontOnly: 750.00, FrontBack: 900.00},
					},
					"2x2": {
						{Qty: 100, FrontOnly: 50.00, FrontBack: 65.00},
						{Qty: 250, FrontOnly: 75.00, FrontBack: 95.00},
						{Qty: 500, FrontOnly: 110.00, FrontBack: 140.00},
						{Qty: 1000, FrontOnly: 180.00, FrontBack: 220.00},
						{Qty: 2500, FrontOnly: 400.00, FrontBack: 480.00},
						{Qty: 5000, FrontOnly: 750.00, FrontBack: 900.00},
					},
					"1.75x3.5": {
						{Qty: 100, FrontOnly: 50.00, FrontBack: 65.00},
						{Qty: 250, FrontOnly: 75.00, FrontBack: 95.00},
						{Qty: 500, FrontOnly: 110.00, FrontBack: 140.00},
						{Qty: 1000, FrontOnly: 180.00, FrontBack: 220.00},
						{Qty: 2500, FrontOnly: 400.00, FrontBack: 480.00},
						{Qty: 5000, FrontOnly: 750.00, FrontBack: 900.00},
					},
				},
			},
			{
				Name:             "Frosted Plastic",
				ShortDescription: "Semi-transparent frosted plastic that diffuses light. Offers a highly modern, tech-forward aesthetic.",
				Features:         []string{"Semi-transparent matte finish", "Waterproof and durable", "Modern aesthetic"},
				PreviewImageURL:  "/static/images/variants/bc-frosted-plastic-preview.webp",
				ImageURLs:        []string{"/static/images/variants/bc-frosted-plastic-1.webp"},
				Pricing: map[string][]QuantityPrice{
					"2x3.5": {
						{Qty: 100, FrontOnly: 55.00, FrontBack: 70.00},
						{Qty: 250, FrontOnly: 85.00, FrontBack: 105.00},
						{Qty: 500, FrontOnly: 130.00, FrontBack: 160.00},
						{Qty: 1000, FrontOnly: 210.00, FrontBack: 250.00},
						{Qty: 2500, FrontOnly: 450.00, FrontBack: 550.00},
						{Qty: 5000, FrontOnly: 850.00, FrontBack: 1050.00},
					},
					"2x2": {
						{Qty: 100, FrontOnly: 55.00, FrontBack: 70.00},
						{Qty: 250, FrontOnly: 85.00, FrontBack: 105.00},
						{Qty: 500, FrontOnly: 130.00, FrontBack: 160.00},
						{Qty: 1000, FrontOnly: 210.00, FrontBack: 250.00},
						{Qty: 2500, FrontOnly: 450.00, FrontBack: 550.00},
						{Qty: 5000, FrontOnly: 850.00, FrontBack: 1050.00},
					},
					"1.75x3.5": {
						{Qty: 100, FrontOnly: 55.00, FrontBack: 70.00},
						{Qty: 250, FrontOnly: 85.00, FrontBack: 105.00},
						{Qty: 500, FrontOnly: 130.00, FrontBack: 160.00},
						{Qty: 1000, FrontOnly: 210.00, FrontBack: 250.00},
						{Qty: 2500, FrontOnly: 450.00, FrontBack: 550.00},
						{Qty: 5000, FrontOnly: 850.00, FrontBack: 1050.00},
					},
				},
			},
			{
				Name:             "Clear Plastic",
				ShortDescription: "100% transparent clear plastic. Colors are printed semi-opaque for a striking stained-glass effect.",
				Features:         []string{"Fully transparent base", "Striking visual effects", "Waterproof"},
				PreviewImageURL:  "/static/images/variants/bc-clear-plastic-preview.webp",
				ImageURLs:        []string{"/static/images/variants/bc-clear-plastic-1.webp"},
				Pricing: map[string][]QuantityPrice{
					"2x3.5": {
						{Qty: 100, FrontOnly: 55.00, FrontBack: 70.00},
						{Qty: 250, FrontOnly: 85.00, FrontBack: 105.00},
						{Qty: 500, FrontOnly: 130.00, FrontBack: 160.00},
						{Qty: 1000, FrontOnly: 210.00, FrontBack: 250.00},
						{Qty: 2500, FrontOnly: 450.00, FrontBack: 550.00},
						{Qty: 5000, FrontOnly: 850.00, FrontBack: 1050.00},
					},
					"2x2": {
						{Qty: 100, FrontOnly: 55.00, FrontBack: 70.00},
						{Qty: 250, FrontOnly: 85.00, FrontBack: 105.00},
						{Qty: 500, FrontOnly: 130.00, FrontBack: 160.00},
						{Qty: 1000, FrontOnly: 210.00, FrontBack: 250.00},
						{Qty: 2500, FrontOnly: 450.00, FrontBack: 550.00},
						{Qty: 5000, FrontOnly: 850.00, FrontBack: 1050.00},
					},
					"1.75x3.5": {
						{Qty: 100, FrontOnly: 55.00, FrontBack: 70.00},
						{Qty: 250, FrontOnly: 85.00, FrontBack: 105.00},
						{Qty: 500, FrontOnly: 130.00, FrontBack: 160.00},
						{Qty: 1000, FrontOnly: 210.00, FrontBack: 250.00},
						{Qty: 2500, FrontOnly: 450.00, FrontBack: 550.00},
						{Qty: 5000, FrontOnly: 850.00, FrontBack: 1050.00},
					},
				},
			},
		},
		FAQs: []FAQ{
			{Question: "What size should my AI prompt generate?", Answer: "Generate your image in a 2:3.5 ratio. Don't worry about the exact pixels; our system will upscale it to the required 300 DPI automatically."},
			{Question: "Can I print on both sides?", Answer: "Yes! You can upload a separate generated image for the front and the back. Note: Clear plastic does not support double-sided printing effectively."},
			{Question: "Will my RGB image look different when printed?", Answer: "Screens use RGB light, while printers use CMYK ink. We use advanced color-mapping to ensure the printed result matches your digital design as closely as physically possible."},
		},
	},
	"letterhead": {
		Title:           "Letterhead",
		MetaTitle:       "Custom Corporate Letterhead Printing | PrintFromAI",
		MetaDescription: "Print your custom 8.5x11 corporate letterhead on premium 70lb uncoated text. Perfect for automated, AI-generated stationery.",
		Slug:            "letterhead",
		Description:     "Professional Custom Corporate Letterhead. Premium uncoated text perfect for your business correspondence.",
		LongDescription: "Establish absolute authority with press-printed corporate letterhead. Whether you have generated a minimalist letterhead via AI or designed a complex watermark, our 70lb premium uncoated text stock ensures a luxurious feel in the hands of your clients. This stock is fully compatible with your home or office laser/inkjet printers, allowing you to feed these pre-printed shells through your own machines to add variable text later.",
		PrimaryImageURL: "https://printfromai.com/static/images/og-letterhead.jpg",
		JSONSchema:      `{"@context": "https://schema.org/", "@type": "Product", "name": "Custom Corporate Letterhead", "description": "8.5x11 Premium 70lb Letterhead printing.", "brand": {"@type": "Brand", "name": "PrintFromAI"}}`,
		SEOKeywords:     []string{"custom letterhead", "corporate stationery", "70lb text letterhead", "8.5x11 print"},
		Features: []string{
			"Standard 8.5x11\" size",
			"Laser and inkjet printer safe",
			"Premium 70lb uncoated text stock",
			"Full color printing available on one or both sides",
		},
		TurnaroundDays: "3-5",
		VariantLabel:   "Stock",
		Sizes: []SizeOption{
			{Label: "8.5 x 11\"", Value: "8.5x11"},
		},
		ReviewScore: 4.8,
		ReviewCount: 215,
		CustomerImages: []string{
			"https://placehold.co/800x1131/1A1A1A/EDEDED?text=Letterhead+1",
			"https://placehold.co/800x1131/1A1A1A/EDEDED?text=Letterhead+2",
		},
		CustomerReviews: []Review{
			{Author: "Valeria Gomez", Rating: 5, Date: "October 02, 2026", Text: "Feeds perfectly through our office laser printer without smudging. The AI upscaling made our generated logo look crisp."},
		},
		Variants: []Variant{
			{
				Name:             "Standard 70lb",
				ShortDescription: "Industry-standard 70lb premium uncoated text. Excellent writability and fully compatible with office printers.",
				Features:         []string{"70lb Uncoated Text", "Laser/Inkjet safe", "Bright white finish"},
				PreviewImageURL:  "/static/images/variants/lh-standard-preview.webp",
				ImageURLs:        []string{"/static/images/variants/lh-standard-1.webp"},
				Pricing: map[string][]QuantityPrice{
					"8.5x11": {
						{Qty: 250, FrontOnly: 80.00, FrontBack: 110.00},
						{Qty: 500, FrontOnly: 110.00, FrontBack: 150.00},
						{Qty: 1000, FrontOnly: 160.00, FrontBack: 210.00},
						{Qty: 2500, FrontOnly: 290.00, FrontBack: 370.00},
						{Qty: 5000, FrontOnly: 520.00, FrontBack: 680.00},
					},
				},
			},
		},
		FAQs: []FAQ{
			{Question: "Can I run this letterhead through my office printer?", Answer: "Yes. Our 70lb uncoated stock is guaranteed to be laser and inkjet safe."},
			{Question: "Does the pricing include full bleed?", Answer: "Yes, your design can extend all the way to the edge of the 8.5x11 sheet at no extra cost."},
		},
	},
	"door-hangers": {
		Title:           "Door Hangers",
		MetaTitle:       "Custom Door Hanger Printing for AI Designs | PrintFromAI",
		MetaDescription: "High-impact custom door hangers printed on premium 16pt semi-gloss stock. Available in 3.5x8.5 and 4.25x11 sizes.",
		Slug:            "door-hangers",
		Description:     "Take your local marketing to the next level with heavy-duty, die-cut door hangers.",
		LongDescription: "Door hangers offer guaranteed visibility. By printing your AI-generated marketing assets on our ultra-thick 16pt semi-gloss stock, you ensure your message withstands the outdoors and catches the eye. We include standard die-cutting for the door knob hole and slit at no extra charge. With both standard (3.5x8.5) and jumbo (4.25x11) sizes available, your artwork will scale perfectly to dominate local advertising.",
		PrimaryImageURL: "https://printfromai.com/static/images/og-door-hangers.jpg",
		JSONSchema:      `{"@context": "https://schema.org/", "@type": "Product", "name": "Custom Door Hangers", "description": "16pt Semi-Gloss Door Hanger printing.", "brand": {"@type": "Brand", "name": "PrintFromAI"}}`,
		SEOKeywords:     []string{"door hanger printing", "custom door hangers", "die cut print", "local marketing print"},
		Features: []string{
			"Ultra-thick 16pt paper stock",
			"Vibrant Semi-Gloss coating",
			"Die-cut hole and slit included",
			"Available in standard and jumbo sizes",
		},
		TurnaroundDays: "4-6",
		VariantLabel:   "Material",
		Sizes: []SizeOption{
			{Label: "3.5 x 8.5\"", Value: "3.5x8.5"},
			{Label: "4.25 x 11\" (Jumbo)", Value: "4.25x11"},
		},
		ReviewScore: 4.7,
		ReviewCount: 412,
		CustomerImages: []string{
			"https://placehold.co/600x1200/1A1A1A/EDEDED?text=Hanger+1",
			"https://placehold.co/600x1200/1A1A1A/EDEDED?text=Hanger+2",
		},
		CustomerReviews: []Review{
			{Author: "Jim Bradley", Rating: 5, Date: "July 14, 2026", Text: "The colors popped perfectly and the die-cut was clean. Very happy with the 16pt thickness, they don't blow away in the wind easily."},
		},
		Variants: []Variant{
			{
				Name:             "16pt Semi-Gloss",
				ShortDescription: "Ultra-thick cardstock coated in a protective semi-gloss finish.",
				Features:         []string{"Standard die-cut hole", "16pt semi-gloss"},
				PreviewImageURL:  "/static/images/variants/dh-standard-preview.webp",
				ImageURLs:        []string{"/static/images/variants/dh-standard-1.webp"},
				Pricing: map[string][]QuantityPrice{
					"3.5x8.5": {
						{Qty: 250, FrontOnly: 95.00, FrontBack: 125.00},
						{Qty: 500, FrontOnly: 125.00, FrontBack: 155.00},
						{Qty: 1000, FrontOnly: 180.00, FrontBack: 220.00},
						{Qty: 2500, FrontOnly: 340.00, FrontBack: 410.00},
						{Qty: 5000, FrontOnly: 620.00, FrontBack: 780.00},
					},
					"4.25x11": {
						{Qty: 250, FrontOnly: 130.00, FrontBack: 170.00},
						{Qty: 500, FrontOnly: 170.00, FrontBack: 210.00},
						{Qty: 1000, FrontOnly: 240.00, FrontBack: 290.00},
						{Qty: 2500, FrontOnly: 450.00, FrontBack: 550.00},
						{Qty: 5000, FrontOnly: 850.00, FrontBack: 1050.00},
					},
				},
			},
		},
		FAQs: []FAQ{
			{Question: "Do I need to add the hole cutout to my design?", Answer: "No! Just provide your full rectangular artwork. We have standard die-cut templates and our system handles cutting the hole automatically."},
		},
	},
	"envelopes": {
		Title:           "Envelopes",
		MetaTitle:       "Custom #10 Envelopes Printed | PrintFromAI",
		MetaDescription: "Print your custom #10 business envelopes on premium 70lb opaque stock. Full color front printing.",
		Slug:            "envelopes",
		Description:     "Standard #10 business envelopes printed on high-quality 70lb opaque stock.",
		LongDescription: "Complete your corporate identity package with full-color custom #10 envelopes. Whether you are adding a sophisticated AI-generated logo or a full-bleed background pattern, our 70lb opaque stock ensures the contents of the envelope remain private while presenting a premium exterior. These standard-sized, no-window envelopes pair perfectly with our custom letterhead.",
		PrimaryImageURL: "https://printfromai.com/static/images/og-envelopes.jpg",
		JSONSchema:      `{"@context": "https://schema.org/", "@type": "Product", "name": "Custom #10 Envelopes", "description": "#10 No Window Envelopes on 70lb opaque stock.", "brand": {"@type": "Brand", "name": "PrintFromAI"}}`,
		SEOKeywords:     []string{"#10 envelopes", "custom envelope print", "70lb opaque envelopes", "business envelopes"},
		Features: []string{
			"Standard #10 Size (4.125 x 9.5\")",
			"No Window design",
			"Premium 70lb Opaque stock for privacy",
			"Full color printing (Front only)",
		},
		TurnaroundDays: "4-6",
		VariantLabel:   "Stock",
		Sizes: []SizeOption{
			{Label: "#10 No Window", Value: "#10"},
		},
		ReviewScore: 4.9,
		ReviewCount: 180,
		CustomerImages: []string{
			"https://placehold.co/800x400/1A1A1A/EDEDED?text=Envelope+1",
		},
		CustomerReviews: []Review{
			{Author: "TechCorp Inc.", Rating: 5, Date: "June 03, 2026", Text: "Opaque paper holds true. You cannot see the contents through the envelope at all. Beautiful front print quality."},
		},
		Variants: []Variant{
			{
				Name:             "70lb Opaque",
				ShortDescription: "Standard 4.125\" x 9.5\" business envelope. Made with 70lb opaque text to prevent contents from showing through.",
				Features:         []string{"4.125\" x 9.5\" dimension", "Opaque privacy layer", "Full color front print"},
				PreviewImageURL:  "/static/images/variants/env-standard-preview.webp",
				ImageURLs:        []string{"/static/images/variants/env-standard-1.webp"},
				Pricing: map[string][]QuantityPrice{
					"#10": {
						{Qty: 250, FrontOnly: 120.00, FrontBack: 0.00},
						{Qty: 500, FrontOnly: 160.00, FrontBack: 0.00},
						{Qty: 1000, FrontOnly: 240.00, FrontBack: 0.00},
						{Qty: 2500, FrontOnly: 450.00, FrontBack: 0.00},
						{Qty: 5000, FrontOnly: 850.00, FrontBack: 0.00},
					},
				},
			},
		},
		FAQs: []FAQ{
			{Question: "Can I print on the back flap?", Answer: "Currently, we only support full-color printing on the front side (4/0) of the envelopes."},
			{Question: "Is the paper thick enough to hide checks?", Answer: "Yes, our 70lb opaque stock is specifically chosen for its security and privacy, preventing text from showing through."},
		},
	},
	"postcards": {
		Title:           "Postcards",
		MetaTitle:       "AI-Generated Art Postcards | PrintFromAI",
		MetaDescription: "Turn your AI-generated art into stunning physical postcards. Premium 16pt cardstock with various coatings.",
		Slug:            "postcards",
		Description:     "Vibrant, heavy-duty postcards perfect for direct mail, handouts, or selling your digital art physically.",
		LongDescription: "Postcards are the perfect canvas for AI art. Printed on rigid 16pt cardstock. What makes our postcard pipeline unique is our flexible coating options at no extra charge. Choose high-gloss for vibrant, punchy colors, elegant matte for a sophisticated look, or a mix of both (Gloss Front / Matte Back) so you can easily write handwritten notes on the reverse side.",
		PrimaryImageURL: "https://printfromai.com/static/images/og-postcards.jpg",
		JSONSchema:      `{"@context": "https://schema.org/", "@type": "Product", "name": "Premium Postcards", "description": "16pt Postcards printed from AI designs."}`,
		SEOKeywords:     []string{"custom postcards", "print AI art", "16pt postcards"},
		Features: []string{
			"Heavy-duty 16pt cardstock",
			"Mix and match Gloss or Matte coatings for free",
			"Writeable matte back option available",
			"Free AI upscaling to press-ready 300 DPI",
		},
		TurnaroundDays: "2-4",
		VariantLabel:   "Stock",
		Sizes: []SizeOption{
			{Label: "4 x 6\"", Value: "4x6"},
			{Label: "5 x 7\"", Value: "5x7"},
			{Label: "6 x 9\" (Jumbo)", Value: "6x9"},
		},
		ReviewScore: 5.0,
		ReviewCount: 892,
		CustomerImages: []string{
			"https://placehold.co/600x400/1A1A1A/EDEDED?text=Postcard+1",
			"https://placehold.co/400x600/1A1A1A/EDEDED?text=Postcard+2",
		},
		CustomerReviews: []Review{
			{Author: "Alice Wonderland", Rating: 5, Date: "November 01, 2026", Text: "Selling my Midjourney art at conventions using these 5x7s. The glossy front makes the colors look radioactive!"},
		},
		Variants: []Variant{
			{
				Name:             "16pt Cardstock",
				ShortDescription: "Ultra-thick, premium 16pt stock available in multiple protective finishes.",
				Features:         []string{"Rigid 16pt thickness", "USPS EDDM Eligible", "Premium feel"},
				PreviewImageURL:  "/static/images/variants/pc-preview.webp",
				ImageURLs:        []string{"/static/images/variants/pc-1.webp"},
				Options: []VariantOption{
					{
						Title: "Coating Options",
						Choices: []OptionChoice{
							{Label: "UV Gloss (Both Sides)", Value: "gloss-both", Upcharges: map[int]float64{}},
							{Label: "Gloss Front / Matte Back", Value: "gloss-front-matte-back", Upcharges: map[int]float64{}},
							{Label: "Matte (Both Sides)", Value: "matte-both", Upcharges: map[int]float64{}},
						},
					},
				},
				Pricing: map[string][]QuantityPrice{
					"4x6": {
						{Qty: 100, FrontOnly: 25.00, FrontBack: 35.00},
						{Qty: 250, FrontOnly: 35.00, FrontBack: 45.00},
						{Qty: 500, FrontOnly: 45.00, FrontBack: 60.00},
						{Qty: 1000, FrontOnly: 65.00, FrontBack: 85.00},
						{Qty: 2500, FrontOnly: 120.00, FrontBack: 160.00},
						{Qty: 5000, FrontOnly: 210.00, FrontBack: 280.00},
					},
					"5x7": {
						{Qty: 100, FrontOnly: 35.00, FrontBack: 50.00},
						{Qty: 250, FrontOnly: 50.00, FrontBack: 70.00},
						{Qty: 500, FrontOnly: 70.00, FrontBack: 95.00},
						{Qty: 1000, FrontOnly: 110.00, FrontBack: 140.00},
						{Qty: 2500, FrontOnly: 210.00, FrontBack: 280.00},
						{Qty: 5000, FrontOnly: 380.00, FrontBack: 500.00},
					},
					"6x9": {
						{Qty: 100, FrontOnly: 45.00, FrontBack: 65.00},
						{Qty: 250, FrontOnly: 70.00, FrontBack: 95.00},
						{Qty: 500, FrontOnly: 100.00, FrontBack: 130.00},
						{Qty: 1000, FrontOnly: 150.00, FrontBack: 195.00},
						{Qty: 2500, FrontOnly: 290.00, FrontBack: 390.00},
						{Qty: 5000, FrontOnly: 520.00, FrontBack: 700.00},
					},
				},
			},
		},
		FAQs: []FAQ{
			{Question: "Which coating should I choose?", Answer: "If you want to write on the postcard, select Matte for that side. A popular choice is Gloss Front and Matte Back."},
		},
	},
}
