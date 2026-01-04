// Shared Cart Functionality
// Cart state stored in localStorage for persistence across pages
let cart = [];

// Initialize cart from localStorage
function initCart() {
    const savedCart = localStorage.getItem('shopping_cart');
    if (savedCart) {
        try {
            cart = JSON.parse(savedCart);
        } catch (e) {
            cart = [];
        }
    }
    updateCartUI();
}

// Save cart to localStorage
function saveCart() {
    localStorage.setItem('shopping_cart', JSON.stringify(cart));
}

// Update cart count badge (only define if not already defined by page)
// Pages like shop.html define their own updateCartUI that also renders cart items
if (typeof window.updateCartUI === 'undefined') {
    window.updateCartUI = function() {
        const cartCount = document.getElementById('cartCount');
        if (cartCount) {
            const totalItems = cart.reduce((sum, item) => sum + item.quantity, 0);
            cartCount.textContent = totalItems;
        }
    };
}

// Create a reference for easier calling
const updateCartUI = window.updateCartUI;

// Toggle cart (for pages with cart sidebar)
if (typeof window.toggleCart === 'undefined') {
    window.toggleCart = function() {
        console.log('toggleCart called');
        // If on shop page, toggle the cart sidebar
        const cartSidebar = document.getElementById('cartSidebar');
        console.log('cartSidebar element:', cartSidebar);
        if (cartSidebar) {
            console.log('Before toggle, classes:', cartSidebar.className);
            cartSidebar.classList.toggle('open');
            console.log('After toggle, classes:', cartSidebar.className);
        } else {
            console.log('cartSidebar not found, redirecting to shop.html');
            // If not on shop page, navigate to shop page
            window.location.href = 'shop.html';
        }
    };
}
const toggleCart = window.toggleCart;

// Add item to cart
if (typeof window.addToCart === 'undefined') {
    window.addToCart = function(product) {
        const existingItem = cart.find(item => item.id === product.id);
        if (existingItem) {
            existingItem.quantity++;
        } else {
            cart.push({ ...product, quantity: 1 });
        }
        saveCart();
        updateCartUI();
    };
}
const addToCart = window.addToCart;

// Remove item from cart
if (typeof window.removeFromCart === 'undefined') {
    window.removeFromCart = function(productId) {
        cart = cart.filter(item => item.id !== productId);
        saveCart();
        updateCartUI();
    };
}
const removeFromCart = window.removeFromCart;

// Update item quantity
if (typeof window.updateQuantity === 'undefined') {
    window.updateQuantity = function(productId, quantity) {
        const item = cart.find(item => item.id === productId);
        if (item) {
            item.quantity = Math.max(0, quantity);
            if (item.quantity === 0) {
                removeFromCart(productId);
            } else {
                saveCart();
                updateCartUI();
            }
        }
    };
}
const updateQuantity = window.updateQuantity;

// Get cart total
function getCartTotal() {
    return cart.reduce((total, item) => {
        const price = item.finalPrice || item.price || 0;
        return total + (price * item.quantity);
    }, 0);
}

// Initialize cart when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initCart);
} else {
    initCart();
}
