package base

import (
	"../../response"

	"github.com/labstack/echo"
)

// InitHelp initialze help api
func InitHelp(parentRoute *echo.Group) {
	parentRoute.GET("/public/help/:lang", readHelp)
}

func readHelp(c echo.Context) error {
	type menu struct {
		Title    string `json:"title"`
		Image    string `json:"image"`
		URL      string `json:"url"`
		SubMenus []menu `json:"subMenus"`
	}
	helps := []menu{
		{"Past Orders", "helps/icon.png", "helps/hello.html", nil},
		{"Account and Payment Options", "helps/icon.png", "", []menu{
			{"Account Settings", "", "", []menu{
				{"I forgot my password", "", "", nil},
				{"How do I review and download a receipt", "", "helps/hello.html", nil},
				{"Update account information", "", "helps/hello.html", nil},
				{"Verifying your account", "", "helps/hello.html", nil},
				{"Updating saved places", "", "helps/hello.html", nil},
				{"Receipts and order history", "", "helps/hello.html", nil},
			}},
			{"Payment", "", "", []menu{
				{"How do I set up or update a payment method?", "", "helps/hello.html", nil},
				{"I have a different payment issue or question", "", "helps/hello.html", nil},
			}},
			{"Promotions and Referrals", "", "", []menu{
				{"How do promo codes and credits work?", "", "helps/hello.html", nil},
				{"Can I choose when to apply a promo code?", "", "helps/hello.html", nil},
				{"I like free fodd. How can I get more of it?", "helps/hello.html", "helps/hello.html", nil},
				{"I had an issue with my invite code", "", "helps/hello.html", nil},
			}},
		}},
	}
	return response.SuccessInterface(c, helps)
}
