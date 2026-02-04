// Navigation functionality for mobile hamburger menu
document.addEventListener('DOMContentLoaded', function() {
    const hamburger = document.querySelector('.hamburger');
    const navLinks = document.querySelector('.nav-links');
    const navLinksItems = document.querySelectorAll('.nav-links a');
    const navOverlay = document.querySelector('.nav-overlay');
    const body = document.body;

    if (!hamburger || !navLinks || !navOverlay) {
        console.warn('Navigation elements not found');
        return;
    }

    function toggleMobileNav() {
        hamburger.classList.toggle('active');
        navLinks.classList.toggle('active');
        navOverlay.classList.toggle('active');
        body.style.overflow = navLinks.classList.contains('active') ? 'hidden' : '';
    }

    hamburger.addEventListener('click', toggleMobileNav);
    navOverlay.addEventListener('click', toggleMobileNav);

    // Close mobile menu when a link is clicked
    navLinksItems.forEach(link => {
        link.addEventListener('click', () => {
            hamburger.classList.remove('active');
            navLinks.classList.remove('active');
            navOverlay.classList.remove('active');
            body.style.overflow = '';
        });
    });

    // User profile dropdown menu
    setupUserDropdown();
});

function setupUserDropdown() {
    // Wait a bit for auth to initialize and add user elements
    setTimeout(() => {
        const userInfo = document.querySelector('.auth-user-info');
        if (!userInfo) {
            console.log('User info not found, retrying...');
            setTimeout(setupUserDropdown, 500);
            return;
        }
        console.log('Setting up user dropdown...');

        // Create dropdown menu if it doesn't exist
        let dropdown = userInfo.querySelector('.user-dropdown');
        if (!dropdown) {
            dropdown = document.createElement('div');
            dropdown.className = 'user-dropdown';
            dropdown.innerHTML = `
                <a href="orders.html" class="dropdown-item">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M9 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8"></path>
                        <path d="M9 2v6h6"></path>
                    </svg>
                    My Orders
                </a>
                <a href="#" class="dropdown-item" onclick="event.preventDefault(); if(typeof logout === 'function') logout();">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
                        <polyline points="16 17 21 12 16 7"></polyline>
                        <line x1="21" y1="12" x2="9" y2="12"></line>
                    </svg>
                    Logout
                </a>
            `;
            userInfo.appendChild(dropdown);

            // Add styles for dropdown
            if (!document.querySelector('#user-dropdown-styles')) {
                const style = document.createElement('style');
                style.id = 'user-dropdown-styles';
                style.textContent = `
                    .auth-user-info {
                        position: relative;
                        cursor: pointer;
                    }

                    .user-dropdown {
                        position: absolute;
                        top: calc(100% + 10px);
                        right: 0;
                        background: white;
                        border-radius: 8px;
                        box-shadow: 0 4px 12px rgba(0,0,0,0.15);
                        min-width: 180px;
                        opacity: 0;
                        visibility: hidden;
                        transform: translateY(-10px);
                        transition: all 0.3s ease;
                        z-index: 1000;
                    }

                    .user-dropdown::before {
                        content: '';
                        position: absolute;
                        top: -6px;
                        right: 20px;
                        width: 12px;
                        height: 12px;
                        background: white;
                        transform: rotate(45deg);
                    }

                    .auth-user-info.dropdown-open .user-dropdown {
                        opacity: 1;
                        visibility: visible;
                        transform: translateY(0);
                    }

                    .dropdown-item {
                        display: flex;
                        align-items: center;
                        gap: 0.75rem;
                        padding: 0.75rem 1rem;
                        color: #333;
                        text-decoration: none;
                        transition: background 0.2s;
                        border-bottom: 1px solid #f0f0f0;
                    }

                    .dropdown-item:first-child {
                        border-radius: 8px 8px 0 0;
                    }

                    .dropdown-item:last-child {
                        border-bottom: none;
                        border-radius: 0 0 8px 8px;
                    }

                    .dropdown-item:hover {
                        background: #f5f5f5;
                    }

                    .dropdown-item svg {
                        width: 18px;
                        height: 18px;
                        color: var(--primary-color, #2E7D32);
                    }
                `;
                document.head.appendChild(style);
            }
        }

        // Toggle dropdown on click
        userInfo.addEventListener('click', function(e) {
            e.stopPropagation();
            userInfo.classList.toggle('dropdown-open');
        });

        // Close dropdown when clicking outside
        document.addEventListener('click', function(e) {
            if (!userInfo.contains(e.target)) {
                userInfo.classList.remove('dropdown-open');
            }
        });

        // Close dropdown when clicking a link
        dropdown.querySelectorAll('.dropdown-item').forEach(item => {
            item.addEventListener('click', function() {
                userInfo.classList.remove('dropdown-open');
            });
        });
    }, 500);
}
