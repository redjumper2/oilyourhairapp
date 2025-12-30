// Auth Module Integration for OilYourHair.com
// Auto-detect environment (production vs development)
const isProduction = window.location.hostname === 'oilyourhair.com' || window.location.hostname === 'www.oilyourhair.com';

const AUTH_CONFIG = {
    portalUrl: isProduction ? 'https://auth.oilyourhair.com' : 'http://localhost:5173',
    apiUrl: isProduction ? 'https://api.oilyourhair.com/api/v1' : 'http://localhost:9090/api/v1',
    domain: 'oilyourhair.com',
    redirectUrl: window.location.origin
};

// Current user state
let currentUser = null;

// Extract token from URL hash and store it
function extractTokenFromHash() {
    const hash = window.location.hash.substring(1);
    const params = new URLSearchParams(hash);
    const token = params.get('token');

    if (token) {
        localStorage.setItem('auth_token', token);
        // Clean up URL
        window.history.replaceState({}, document.title, window.location.pathname + window.location.search);
        return token;
    }

    return null;
}

// Get stored auth token
function getAuthToken() {
    return localStorage.getItem('auth_token');
}

// Remove auth token
function clearAuthToken() {
    localStorage.removeItem('auth_token');
    currentUser = null;
}

// Fetch current user info from API
async function fetchCurrentUser() {
    const token = getAuthToken();
    if (!token) {
        return null;
    }

    try {
        const response = await fetch(`${AUTH_CONFIG.apiUrl}/auth/me`, {
            headers: {
                'Authorization': `Bearer ${token}`,
                'Host': AUTH_CONFIG.domain
            }
        });

        if (!response.ok) {
            if (response.status === 401) {
                // Token expired or invalid
                clearAuthToken();
                return null;
            }
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const user = await response.json();
        currentUser = user;
        return user;
    } catch (error) {
        console.error('Error fetching current user:', error);
        clearAuthToken();
        return null;
    }
}

// Check if user is authenticated
async function isAuthenticated() {
    // First check if there's a token in the URL hash
    extractTokenFromHash();

    // Then get the token from storage
    const token = getAuthToken();
    if (!token) {
        return false;
    }

    // Verify token is valid by fetching user info
    const user = await fetchCurrentUser();
    return user !== null;
}

// Redirect to login portal
function login() {
    const redirectUrl = encodeURIComponent(AUTH_CONFIG.redirectUrl);
    const domain = encodeURIComponent(AUTH_CONFIG.domain);
    window.location.href = `${AUTH_CONFIG.portalUrl}/login?domain=${domain}&redirect=${redirectUrl}`;
}

// Logout user
function logout() {
    clearAuthToken();
    window.location.reload();
}

// Check if user has a specific permission
function hasPermission(permission) {
    if (!currentUser || !currentUser.permissions) {
        return false;
    }
    return currentUser.permissions.includes(permission);
}

// Initialize auth on page load
async function initAuth() {
    const authenticated = await isAuthenticated();

    // Update UI based on auth state
    updateAuthUI(authenticated);

    return authenticated;
}

// Update UI elements based on authentication state
function updateAuthUI(authenticated) {
    // Update auth status display
    const authStatus = document.getElementById('auth-status');
    const loginBtn = document.getElementById('login-btn');
    const logoutBtn = document.getElementById('logout-btn');
    const userEmail = document.getElementById('user-email');

    if (authenticated && currentUser) {
        if (authStatus) {
            authStatus.textContent = `Welcome, ${currentUser.email}`;
            authStatus.style.display = 'inline';
        }
        if (loginBtn) loginBtn.style.display = 'none';
        if (logoutBtn) logoutBtn.style.display = 'inline-block';
        if (userEmail) userEmail.textContent = currentUser.email;

        // Show elements that require authentication
        document.querySelectorAll('[data-auth-required]').forEach(el => {
            el.style.display = '';
        });

        // Hide elements that should only show when not authenticated
        document.querySelectorAll('[data-auth-hidden]').forEach(el => {
            el.style.display = 'none';
        });

        // Show prices
        document.querySelectorAll('.price, .product-price').forEach(el => {
            el.style.display = '';
        });

        // Hide "login to view price" messages
        document.querySelectorAll('.login-to-view-price').forEach(el => {
            el.style.display = 'none';
        });

    } else {
        if (authStatus) {
            authStatus.textContent = '';
            authStatus.style.display = 'none';
        }
        if (loginBtn) loginBtn.style.display = 'inline-block';
        if (logoutBtn) logoutBtn.style.display = 'none';
        if (userEmail) userEmail.textContent = '';

        // Hide elements that require authentication
        document.querySelectorAll('[data-auth-required]').forEach(el => {
            el.style.display = 'none';
        });

        // Show elements that should only show when not authenticated
        document.querySelectorAll('[data-auth-hidden]').forEach(el => {
            el.style.display = '';
        });

        // Hide prices
        document.querySelectorAll('.price, .product-price').forEach(el => {
            el.style.display = 'none';
        });

        // Show "login to view price" messages
        document.querySelectorAll('.login-to-view-price').forEach(el => {
            el.style.display = '';
        });
    }
}

// Auto-initialize when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initAuth);
} else {
    initAuth();
}
