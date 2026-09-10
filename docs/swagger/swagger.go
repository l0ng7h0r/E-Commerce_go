package swagger

import (
	"html/template"

	"github.com/gofiber/contrib/v3/swaggo"
)

// NavScript is a custom JavaScript injected into Swagger UI to add Portal Navigation buttons
const NavScript = template.JS(`
window.addEventListener('DOMContentLoaded', function() {
	var interval = setInterval(function() {
		var target = document.querySelector('.auth-wrapper') || document.querySelector('.scheme-container') || document.querySelector('.topbar-wrapper');
		if (target) {
			if (document.getElementById('portal-nav-buttons')) return;

			var container = document.createElement('div');
			container.id = 'portal-nav-buttons';
			container.style.cssText = 'display: inline-flex; gap: 8px; margin-right: 15px; align-items: center;';

			var portals = [
				{ name: 'User Portal', url: '/swagger/user/index.html', activeColor: '#49cc90' },
				{ name: 'Seller Portal', url: '/swagger/seller/index.html', activeColor: '#fca130' },
				{ name: 'Admin Portal', url: '/swagger/admin/index.html', activeColor: '#4990e2' }
			];

			var currentPath = window.location.pathname;

			portals.forEach(function(p) {
				var btn = document.createElement('a');
				btn.href = p.url;
				btn.innerText = p.name;
				var isActive = currentPath.indexOf(p.url.replace('/index.html', '')) !== -1;
				
				var bg = isActive ? p.activeColor : 'transparent';
				var textColor = isActive ? '#ffffff' : p.activeColor;
				
				btn.style.cssText = 'display: inline-block; padding: 4px 14px; font-family: sans-serif; font-size: 13px; font-weight: bold; text-decoration: none; color: ' + textColor + '; border: 2px solid ' + p.activeColor + '; border-radius: 4px; background: ' + bg + '; transition: all 0.2s ease-in-out; cursor: pointer;';
				
				if (!isActive) {
					btn.onmouseover = function() { btn.style.background = p.activeColor; btn.style.color = '#ffffff'; };
					btn.onmouseout = function() { btn.style.background = 'transparent'; btn.style.color = p.activeColor; };
				}

				container.appendChild(btn);
			});

			if (target.firstChild) {
				target.insertBefore(container, target.firstChild);
			} else {
				target.appendChild(container);
			}
		}
	}, 200);
});
`)

// GetUserConfig returns Swaggo Config for User Portal with navigation buttons
func GetUserConfig() swaggo.Config {
	return swaggo.Config{
		Title:        "User API Portal",
		URL:          "/swagger/user/doc.json",
		InstanceName: "user",
		CustomScript: NavScript,
	}
}

// GetSellerConfig returns Swaggo Config for Seller Portal with navigation buttons
func GetSellerConfig() swaggo.Config {
	return swaggo.Config{
		Title:        "Seller API Portal",
		URL:          "/swagger/seller/doc.json",
		InstanceName: "seller",
		CustomScript: NavScript,
	}
}

// GetAdminConfig returns Swaggo Config for Admin Portal with navigation buttons
func GetAdminConfig() swaggo.Config {
	return swaggo.Config{
		Title:        "Admin API Portal",
		URL:          "/swagger/admin/doc.json",
		InstanceName: "admin",
		CustomScript: NavScript,
	}
}
