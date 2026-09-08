package swagger

import (
	"strings"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"

	_ "github.com/l0ng7h0r/ecommerce/docs/admin"
	_ "github.com/l0ng7h0r/ecommerce/docs/seller"
	_ "github.com/l0ng7h0r/ecommerce/docs/user"
)

// Setup registers the 3 Swagger portals and adds Authorize-style portal buttons
func Setup(app *fiber.App) {
	// Middleware to inject Swagger Authorize-style buttons next to the Authorize button
	app.Use("/swagger", func(c fiber.Ctx) error {
		err := c.Next()
		if strings.HasSuffix(c.Path(), "index.html") {
			body := string(c.Response().Body())
			buttonsHTML := `
<style>
  .swagger-ui .portal-switch-container {
    display: flex;
    gap: 10px;
    align-items: center;
    margin-right: 15px;
  }
  .swagger-ui .btn.portal-btn {
    border: 2px solid #49cc90;
    color: #49cc90;
    background-color: transparent;
    border-radius: 4px;
    font-family: sans-serif;
    font-weight: 700;
    font-size: 14px;
    padding: 4px 16px;
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    transition: all 0.2s ease-in-out;
  }
  .swagger-ui .btn.portal-btn:hover {
    background-color: rgba(73, 204, 144, 0.1);
  }
  .swagger-ui .btn.portal-btn.user-btn { border-color: #4990e2; color: #4990e2; }
  .swagger-ui .btn.portal-btn.user-btn:hover { background-color: rgba(73, 144, 226, 0.1); }
  .swagger-ui .btn.portal-btn.user-btn.active { background-color: #4990e2; color: #fff; }

  .swagger-ui .btn.portal-btn.seller-btn { border-color: #f5a623; color: #f5a623; }
  .swagger-ui .btn.portal-btn.seller-btn:hover { background-color: rgba(245, 166, 35, 0.1); }
  .swagger-ui .btn.portal-btn.seller-btn.active { background-color: #f5a623; color: #fff; }

  .swagger-ui .btn.portal-btn.admin-btn { border-color: #e54949; color: #e54949; }
  .swagger-ui .btn.portal-btn.admin-btn:hover { background-color: rgba(229, 73, 73, 0.1); }
  .swagger-ui .btn.portal-btn.admin-btn.active { background-color: #e54949; color: #fff; }
</style>
<script>
  window.addEventListener('DOMContentLoaded', function() {
    var checkExist = setInterval(function() {
      var authWrapper = document.querySelector('.swagger-ui .auth-wrapper') || document.querySelector('.swagger-ui .scheme-container');
      if (authWrapper && !document.querySelector('.portal-switch-container')) {
        var container = document.createElement('div');
        container.className = 'portal-switch-container';
        
        var currentPath = window.location.pathname;
        
        var userActive = currentPath.includes('/user/') ? ' active' : '';
        var sellerActive = currentPath.includes('/seller/') ? ' active' : '';
        var adminActive = currentPath.includes('/admin/') ? ' active' : '';

        container.innerHTML = 
          '<a href="/swagger/user/index.html" class="btn portal-btn user-btn' + userActive + '">User Portal</a>' +
          '<a href="/swagger/seller/index.html" class="btn portal-btn seller-btn' + sellerActive + '">Seller Portal</a>' +
          '<a href="/swagger/admin/index.html" class="btn portal-btn admin-btn' + adminActive + '">Admin Portal</a>';
        
        authWrapper.insertBefore(container, authWrapper.firstChild);
        clearInterval(checkExist);
      }
    }, 100);
  });
</script>
</body>`
			newBody := strings.Replace(body, "</body>", buttonsHTML, 1)
			c.Response().SetBodyString(newBody)
		}
		return err
	})

	// 1. User Swagger Portal
	app.Get("/swagger/user/*", swaggo.New(swaggo.Config{
		Title:        "User API Portal",
		URL:          "/swagger/user/doc.json",
		InstanceName: "user",
	}))

	// 2. Seller Swagger Portal
	app.Get("/swagger/seller/*", swaggo.New(swaggo.Config{
		Title:        "Seller API Portal",
		URL:          "/swagger/seller/doc.json",
		InstanceName: "seller",
	}))

	// 3. Admin Swagger Portal
	app.Get("/swagger/admin/*", swaggo.New(swaggo.Config{
		Title:        "Admin API Portal",
		URL:          "/swagger/admin/doc.json",
		InstanceName: "admin",
	}))
}
