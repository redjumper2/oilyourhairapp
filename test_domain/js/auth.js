/**
 * Auth integration for test domain
 * Handles JWT token extraction from URL hash and authentication state
 */

const AUTH_CONFIG = {
  portalUrl: 'http://localhost:5173',
  apiUrl: 'http://localhost:9090/api/v1',
  domain: 'testdomain.com',
  redirectUrl: 'http://localhost:8000'
};

/**
 * Extract JWT token from URL hash
 */
function extractTokenFromHash() {
  const hash = window.location.hash.substring(1);
  const params = new URLSearchParams(hash);
  const token = params.get('token');

  if (token) {
    // Store token
    localStorage.setItem('auth_token', token);

    // Clean URL (remove hash)
    window.history.replaceState({}, document.title, window.location.pathname);

    // Reload to update UI
    window.location.reload();
  }
}

/**
 * Get current auth token
 */
function getToken() {
  return localStorage.getItem('auth_token');
}

/**
 * Check if user is authenticated
 */
function isAuthenticated() {
  return !!getToken();
}

/**
 * Get current user info from API
 */
async function getCurrentUser() {
  const token = getToken();
  if (!token) return null;

  try {
    const response = await fetch(`${AUTH_CONFIG.apiUrl}/auth/me`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Host': AUTH_CONFIG.domain
      }
    });

    if (!response.ok) {
      // Token invalid, clear it
      logout();
      return null;
    }

    return await response.json();
  } catch (error) {
    console.error('Failed to get user:', error);
    return null;
  }
}

/**
 * Logout user
 */
function logout() {
  localStorage.removeItem('auth_token');
  window.location.reload();
}

/**
 * Redirect to login page
 */
function login() {
  const loginUrl = `${AUTH_CONFIG.portalUrl}/login?domain=${AUTH_CONFIG.domain}&redirect=${encodeURIComponent(AUTH_CONFIG.redirectUrl)}`;
  window.location.href = loginUrl;
}

/**
 * Make authenticated API call
 */
async function apiCall(endpoint, options = {}) {
  const token = getToken();
  if (!token) {
    throw new Error('Not authenticated');
  }

  const response = await fetch(`${AUTH_CONFIG.apiUrl}${endpoint}`, {
    ...options,
    headers: {
      'Authorization': `Bearer ${token}`,
      'Host': AUTH_CONFIG.domain,
      'Content-Type': 'application/json',
      ...options.headers
    }
  });

  if (!response.ok) {
    throw new Error(`API call failed: ${response.status}`);
  }

  return response.json();
}

/**
 * Initialize auth on page load
 */
function initAuth() {
  // Extract token from hash if present
  extractTokenFromHash();

  // Update UI based on auth state
  updateAuthUI();
}

/**
 * Update UI based on authentication state
 */
async function updateAuthUI() {
  const authStatus = document.getElementById('auth-status');
  const loginBtn = document.getElementById('login-btn');
  const logoutBtn = document.getElementById('logout-btn');
  const userInfo = document.getElementById('user-info');

  if (!isAuthenticated()) {
    // Not logged in
    if (authStatus) authStatus.textContent = 'Not logged in';
    if (loginBtn) loginBtn.style.display = 'inline-block';
    if (logoutBtn) logoutBtn.style.display = 'none';
    if (userInfo) userInfo.style.display = 'none';

    // Hide prices
    document.querySelectorAll('.product-price').forEach(el => {
      el.style.display = 'none';
    });
    document.querySelectorAll('.login-to-view').forEach(el => {
      el.style.display = 'block';
    });
  } else {
    // Logged in - fetch user info
    const user = await getCurrentUser();

    if (user) {
      if (authStatus) authStatus.textContent = `Logged in as ${user.email}`;
      if (loginBtn) loginBtn.style.display = 'none';
      if (logoutBtn) logoutBtn.style.display = 'inline-block';

      if (userInfo) {
        userInfo.style.display = 'block';
        userInfo.innerHTML = `
          <strong>Email:</strong> ${user.email}<br>
          <strong>Role:</strong> ${user.role}<br>
          <strong>Domain:</strong> ${user.domain}
        `;
      }

      // Show prices
      document.querySelectorAll('.product-price').forEach(el => {
        el.style.display = 'block';
      });
      document.querySelectorAll('.login-to-view').forEach(el => {
        el.style.display = 'none';
      });
    }
  }
}

// Initialize on page load
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initAuth);
} else {
  initAuth();
}
