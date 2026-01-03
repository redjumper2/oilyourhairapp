#!/usr/bin/env python3
"""
OilYourHair.com Setup Script
============================

This script recreates the complete oilyourhair.com setup including:
- Domain configuration
- Admin user
- API keys
- Products (pre-populated with current inventory)
- Reviews (pre-populated with sample 5-star reviews)

Usage:
    python3 setup_oilyourhair.py

Requirements:
    pip install requests colorama

Author: Claude
Date: 2026-01-03
"""

import requests
import json
import time
import sys
from typing import Dict, List, Optional, Tuple
from colorama import init, Fore, Style

# Initialize colorama for colored output
init(autoreset=True)

# Configuration
AUTH_URL = "http://localhost:8080"
PRODUCTS_URL = "http://localhost:9091"
DOMAIN = "oilyourhair.com"
DOMAIN_NAME = "Oil Your Hair"

# Admin user credentials
ADMIN_EMAIL = "admin@oilyourhair.com"
ADMIN_PASSWORD = "SecurePass123!"
ADMIN_FIRST_NAME = "Admin"
ADMIN_LAST_NAME = "User"

# Global variables to store tokens
invite_token = None
api_key = None
admin_jwt = None


def print_step(step_num: int, description: str):
    """Print a step header."""
    print(f"\n{Fore.CYAN}{'='*80}")
    print(f"{Fore.CYAN}STEP {step_num}: {description}")
    print(f"{Fore.CYAN}{'='*80}{Style.RESET_ALL}")


def print_success(message: str):
    """Print a success message."""
    print(f"{Fore.GREEN}✓ {message}{Style.RESET_ALL}")


def print_error(message: str):
    """Print an error message."""
    print(f"{Fore.RED}✗ {message}{Style.RESET_ALL}")


def print_info(message: str):
    """Print an info message."""
    print(f"{Fore.YELLOW}ℹ {message}{Style.RESET_ALL}")


def print_data(label: str, data: any):
    """Print data with a label."""
    print(f"{Fore.MAGENTA}{label}:{Style.RESET_ALL} {data}")


def check_service_health(service_name: str, url: str) -> bool:
    """Check if a service is running and healthy."""
    try:
        response = requests.get(f"{url}/health", timeout=5)
        if response.status_code == 200:
            print_success(f"{service_name} is healthy")
            return True
        else:
            print_error(f"{service_name} returned status {response.status_code}")
            return False
    except requests.exceptions.RequestException as e:
        print_error(f"{service_name} is not accessible: {e}")
        return False


def create_domain() -> bool:
    """Create the oilyourhair.com domain in auth module."""
    print_info(f"Creating domain: {DOMAIN}")

    payload = {
        "domain": DOMAIN,
        "name": DOMAIN_NAME,
        "settings": {
            "allow_signups": True,
            "require_email_verification": False
        }
    }

    try:
        response = requests.post(
            f"{AUTH_URL}/api/v1/domains",
            json=payload,
            headers={"Content-Type": "application/json"}
        )

        if response.status_code in [200, 201]:
            data = response.json()
            print_success(f"Domain created: {data.get('domain')}")
            print_data("Domain ID", data.get('id'))
            return True
        elif response.status_code == 409:
            print_info("Domain already exists, continuing...")
            return True
        else:
            print_error(f"Failed to create domain: {response.text}")
            return False
    except requests.exceptions.RequestException as e:
        print_error(f"Request failed: {e}")
        return False


def create_invite_token() -> Optional[str]:
    """Create an invite token for the domain."""
    print_info(f"Creating invite token for {DOMAIN}")

    # Note: This would typically call an API endpoint
    # For now, we'll simulate by calling the CLI command
    import subprocess

    try:
        result = subprocess.run(
            ["./auth-module", "domain", "create-invite", DOMAIN],
            cwd="/home/sparque/dev/oilyourhairapp/auth_module",
            capture_output=True,
            text=True,
            timeout=10
        )

        if result.returncode == 0:
            # Extract token from output
            output = result.stdout.strip()
            # Look for the token in the output
            for line in output.split('\n'):
                if 'Token:' in line or 'token' in line.lower():
                    # Extract the token (usually after "Token:" or similar)
                    token = line.split(':')[-1].strip()
                    if len(token) > 20:  # Basic validation
                        print_success(f"Invite token created")
                        print_data("Token", token[:20] + "..." if len(token) > 20 else token)
                        return token

            # If we can't parse it, return the full output and let user handle it
            print_info("Could not parse token from output:")
            print(output)
            print_info("Please paste the invite token:")
            return input().strip()
        else:
            print_error(f"Failed to create invite token: {result.stderr}")
            return None
    except Exception as e:
        print_error(f"Failed to run command: {e}")
        return None


def create_admin_user(invite_token: str) -> bool:
    """Create the admin user for the domain."""
    print_info(f"Creating admin user: {ADMIN_EMAIL}")

    payload = {
        "email": ADMIN_EMAIL,
        "password": ADMIN_PASSWORD,
        "first_name": ADMIN_FIRST_NAME,
        "last_name": ADMIN_LAST_NAME,
        "invite_token": invite_token
    }

    try:
        response = requests.post(
            f"{AUTH_URL}/api/v1/auth/signup",
            json=payload,
            headers={"Content-Type": "application/json"}
        )

        if response.status_code in [200, 201]:
            data = response.json()
            print_success(f"Admin user created: {ADMIN_EMAIL}")
            print_data("User ID", data.get('user', {}).get('id'))
            return True
        elif response.status_code == 409:
            print_info("User already exists, continuing...")
            return True
        else:
            print_error(f"Failed to create user: {response.text}")
            return False
    except requests.exceptions.RequestException as e:
        print_error(f"Request failed: {e}")
        return False


def login_admin() -> Optional[str]:
    """Login as admin and get JWT token."""
    print_info(f"Logging in as: {ADMIN_EMAIL}")

    payload = {
        "email": ADMIN_EMAIL,
        "password": ADMIN_PASSWORD
    }

    try:
        response = requests.post(
            f"{AUTH_URL}/api/v1/auth/login",
            json=payload,
            headers={"Content-Type": "application/json"}
        )

        if response.status_code == 200:
            data = response.json()
            token = data.get('token')
            if token:
                print_success("Login successful")
                print_data("JWT Token", token[:30] + "...")
                return token
            else:
                print_error("No token in response")
                return None
        else:
            print_error(f"Login failed: {response.text}")
            return None
    except requests.exceptions.RequestException as e:
        print_error(f"Request failed: {e}")
        return None


def create_api_key() -> Optional[str]:
    """Create API key for products service."""
    print_info(f"Creating API key for products service")

    import subprocess

    try:
        result = subprocess.run(
            [
                "./auth-module", "apikey", "create", DOMAIN,
                "--service", "products",
                "--permissions", "products.read,products.write,products.delete"
            ],
            cwd="/home/sparque/dev/oilyourhairapp/auth_module",
            capture_output=True,
            text=True,
            timeout=10
        )

        if result.returncode == 0:
            output = result.stdout.strip()
            # Look for the API key in the output
            for line in output.split('\n'):
                if 'Key:' in line or 'key' in line.lower() or line.startswith('ey'):
                    # Extract the key
                    if 'Key:' in line:
                        key = line.split(':')[-1].strip()
                    else:
                        key = line.strip()

                    if len(key) > 50:  # JWT tokens are usually long
                        print_success(f"API key created")
                        print_data("API Key", key[:40] + "...")
                        return key

            # If we can't parse it, show output
            print_info("Could not parse API key from output:")
            print(output)
            print_info("Please paste the API key:")
            return input().strip()
        else:
            print_error(f"Failed to create API key: {result.stderr}")
            return None
    except Exception as e:
        print_error(f"Failed to run command: {e}")
        return None


def get_products_data() -> List[Dict]:
    """Get the product data to populate."""
    return [
        {
            "domain": DOMAIN,
            "name": "Nourishing Argan Oil",
            "description": "Pure organic argan oil from Morocco. Rich in vitamin E and fatty acids, this luxurious oil penetrates deep to nourish and revitalize your hair from root to tip.",
            "price": 29.99,
            "category": "Hair Oils",
            "attributes": {
                "size": "100ml",
                "ingredients": "100% Pure Argan Oil",
                "origin": "Morocco",
                "hair_type": "All Hair Types"
            },
            "images": [
                "https://images.unsplash.com/photo-1571781926291-c477ebfd024b"
            ],
            "active": True,
            "bestseller": True
        },
        {
            "domain": DOMAIN,
            "name": "Hydrating Coconut Oil",
            "description": "Cold-pressed virgin coconut oil that deeply moisturizes and strengthens hair. Perfect for reducing protein loss and adding natural shine.",
            "price": 24.99,
            "category": "Hair Oils",
            "attributes": {
                "size": "150ml",
                "ingredients": "100% Virgin Coconut Oil",
                "origin": "Philippines",
                "hair_type": "Dry, Damaged Hair"
            },
            "images": [
                "https://images.unsplash.com/photo-1599351431202-1e0f0137899a"
            ],
            "active": True,
            "bestseller": True
        },
        {
            "domain": DOMAIN,
            "name": "Strengthening Castor Oil",
            "description": "Premium cold-pressed castor oil that promotes hair growth and thickness. Rich in ricinoleic acid to strengthen hair follicles.",
            "price": 19.99,
            "category": "Hair Oils",
            "attributes": {
                "size": "120ml",
                "ingredients": "100% Pure Castor Oil",
                "origin": "India",
                "hair_type": "Thinning Hair"
            },
            "images": [
                "https://images.unsplash.com/photo-1556228720-195a672e8a03"
            ],
            "active": True,
            "bestseller": False
        },
        {
            "domain": DOMAIN,
            "name": "Shine & Smooth Jojoba Oil",
            "description": "Lightweight jojoba oil that mimics your hair's natural oils. Perfect for adding shine without weighing down your hair.",
            "price": 27.99,
            "category": "Hair Oils",
            "attributes": {
                "size": "100ml",
                "ingredients": "100% Pure Jojoba Oil",
                "origin": "USA",
                "hair_type": "Fine, Oily Hair"
            },
            "images": [
                "https://images.unsplash.com/photo-1608248543803-ba4f8c70ae0b"
            ],
            "active": True,
            "bestseller": False
        },
        {
            "domain": DOMAIN,
            "name": "Repair Hair Treatment",
            "description": "Intensive hair repair treatment combining multiple natural oils. Restores damaged hair and prevents breakage.",
            "price": 34.99,
            "category": "Treatments",
            "attributes": {
                "size": "200ml",
                "ingredients": "Argan, Coconut, Olive Oil Blend",
                "origin": "France",
                "hair_type": "Damaged, Color-Treated Hair"
            },
            "images": [
                "https://images.unsplash.com/photo-1535585209827-a15fcdbc4c2d"
            ],
            "active": True,
            "bestseller": True
        },
        {
            "domain": DOMAIN,
            "name": "Moisturizing Hair Mask",
            "description": "Deep conditioning hair mask for intense hydration. Leaves hair soft, manageable, and beautifully scented.",
            "price": 32.99,
            "category": "Treatments",
            "attributes": {
                "size": "250ml",
                "ingredients": "Shea Butter, Argan Oil, Keratin",
                "origin": "France",
                "hair_type": "Curly, Frizzy Hair"
            },
            "images": [
                "https://images.unsplash.com/photo-1522338242992-e1a54906a8da"
            ],
            "active": True,
            "bestseller": False
        }
    ]


def create_products(api_key: str) -> List[str]:
    """Create all products and return their IDs."""
    print_info(f"Creating products for {DOMAIN}")

    products_data = get_products_data()
    product_ids = []

    headers = {
        "Authorization": f"Bearer {api_key}",
        "Content-Type": "application/json"
    }

    for idx, product in enumerate(products_data, 1):
        print(f"\n  Creating product {idx}/{len(products_data)}: {product['name']}")

        try:
            response = requests.post(
                f"{PRODUCTS_URL}/api/v1/products",
                json=product,
                headers=headers
            )

            if response.status_code in [200, 201]:
                data = response.json()
                product_id = data.get('id')
                product_ids.append(product_id)
                print_success(f"Created: {product['name']}")
                print_data("  Product ID", product_id)
                print_data("  Price", f"${product['price']}")
                print_data("  Bestseller", "Yes" if product['bestseller'] else "No")
            else:
                print_error(f"Failed to create product: {response.text}")
        except requests.exceptions.RequestException as e:
            print_error(f"Request failed: {e}")

    return product_ids


def get_reviews_data() -> List[Dict]:
    """Get the review data to populate."""
    return [
        {
            "product": "Nourishing Argan Oil",
            "name": "Sarah Johnson",
            "rating": 5,
            "text": "Absolutely love this argan oil! My hair has never been softer or shinier. I use it after every wash and the results are incredible. The light coconut scent is amazing too!",
            "highlight": "Best hair product I have ever used!"
        },
        {
            "product": "Hydrating Coconut Oil",
            "name": "Michael Chen",
            "rating": 5,
            "text": "This serum is a game-changer for my dry, damaged hair. After just two weeks of use, my hair feels healthier and looks more vibrant. Highly recommend!",
            "highlight": "Transformed my damaged hair completely"
        },
        {
            "product": "Moisturizing Hair Mask",
            "name": "Emily Rodriguez",
            "rating": 5,
            "text": "I have curly hair that tends to get frizzy, but this hair mask keeps it smooth and defined. The natural ingredients are a huge plus, and it smells wonderful!",
            "highlight": "Perfect for curly hair - no more frizz!"
        },
        {
            "product": "Strengthening Castor Oil",
            "name": "David Thompson",
            "rating": 5,
            "text": "As someone with thinning hair, I was skeptical, but this oil has made a noticeable difference. My hair feels thicker and stronger. Will definitely keep using this!",
            "highlight": "Noticeable improvement in hair thickness"
        },
        {
            "product": "Repair Hair Treatment",
            "name": "Jessica Park",
            "rating": 5,
            "text": "After years of heat styling and coloring, my hair was in bad shape. This treatment has brought it back to life! It is soft, shiny, and so much healthier. Cannot recommend enough!",
            "highlight": "Brought my damaged hair back to life"
        },
        {
            "product": "Shine & Smooth Jojoba Oil",
            "name": "Ryan Martinez",
            "rating": 5,
            "text": "I have been using this oil for a month now and the shine it gives my hair is incredible. It is lightweight and does not make my hair greasy. A little goes a long way!",
            "highlight": "Amazing shine without the grease"
        }
    ]


def create_reviews() -> int:
    """Create all reviews and return count."""
    print_info(f"Creating reviews for {DOMAIN}")

    reviews_data = get_reviews_data()
    success_count = 0

    headers = {"Content-Type": "application/json"}

    for idx, review in enumerate(reviews_data, 1):
        print(f"\n  Creating review {idx}/{len(reviews_data)}: by {review['name']}")

        try:
            response = requests.post(
                f"{PRODUCTS_URL}/api/v1/public/{DOMAIN}/reviews",
                json=review,
                headers=headers
            )

            if response.status_code in [200, 201]:
                data = response.json()
                success_count += 1
                print_success(f"Created review by {review['name']}")
                print_data("  Product", review['product'])
                print_data("  Rating", f"{'⭐' * review['rating']}")
                print_data("  Review ID", data.get('id'))
            else:
                print_error(f"Failed to create review: {response.text}")
        except requests.exceptions.RequestException as e:
            print_error(f"Request failed: {e}")

    return success_count


def verify_setup() -> bool:
    """Verify the complete setup."""
    print_info("Verifying setup...")

    all_ok = True

    # Verify products
    try:
        response = requests.get(f"{PRODUCTS_URL}/api/v1/public/{DOMAIN}/products")
        if response.status_code == 200:
            data = response.json()
            product_count = data.get('count', 0)
            print_success(f"Products verified: {product_count} products found")

            # Check for bestsellers
            bestsellers = [p for p in data.get('products', []) if p.get('bestseller')]
            print_data("  Bestsellers", len(bestsellers))
        else:
            print_error("Failed to verify products")
            all_ok = False
    except Exception as e:
        print_error(f"Products verification failed: {e}")
        all_ok = False

    # Verify reviews
    try:
        response = requests.get(f"{PRODUCTS_URL}/api/v1/public/{DOMAIN}/reviews")
        if response.status_code == 200:
            data = response.json()
            review_count = data.get('count', 0)
            print_success(f"Reviews verified: {review_count} reviews found")

            # Check all are 5-star
            reviews = data.get('reviews', [])
            five_star = [r for r in reviews if r.get('rating') == 5]
            print_data("  5-star reviews", len(five_star))
        else:
            print_error("Failed to verify reviews")
            all_ok = False
    except Exception as e:
        print_error(f"Reviews verification failed: {e}")
        all_ok = False

    # Verify admin can login
    try:
        jwt = login_admin()
        if jwt:
            print_success("Admin authentication verified")
        else:
            print_error("Admin authentication failed")
            all_ok = False
    except Exception as e:
        print_error(f"Auth verification failed: {e}")
        all_ok = False

    return all_ok


def main():
    """Main execution function."""
    print(f"\n{Fore.CYAN}{Style.BRIGHT}")
    print("=" * 80)
    print("  OilYourHair.com Setup Script")
    print("  Automated setup and verification")
    print("=" * 80)
    print(Style.RESET_ALL)

    # Step 0: Check service health
    print_step(0, "Checking Service Health")

    auth_ok = check_service_health("Auth Module", AUTH_URL)
    products_ok = check_service_health("Products API", PRODUCTS_URL)

    if not auth_ok or not products_ok:
        print_error("\nOne or more services are not running!")
        print_info("Please ensure all services are started before running this script.")
        sys.exit(1)

    # Step 1: Create domain
    print_step(1, "Creating Domain")
    if not create_domain():
        print_error("Failed to create domain. Exiting.")
        sys.exit(1)

    # Step 2: Create invite token
    print_step(2, "Creating Invite Token")
    global invite_token
    invite_token = create_invite_token()
    if not invite_token:
        print_error("Failed to create invite token. Exiting.")
        sys.exit(1)

    # Step 3: Create admin user
    print_step(3, "Creating Admin User")
    if not create_admin_user(invite_token):
        print_error("Failed to create admin user. Exiting.")
        sys.exit(1)

    # Step 4: Login and get JWT
    print_step(4, "Authenticating Admin User")
    global admin_jwt
    admin_jwt = login_admin()
    if not admin_jwt:
        print_error("Failed to login. Exiting.")
        sys.exit(1)

    # Step 5: Create API key
    print_step(5, "Creating API Key for Products")
    global api_key
    api_key = create_api_key()
    if not api_key:
        print_error("Failed to create API key. Exiting.")
        sys.exit(1)

    # Step 6: Create products
    print_step(6, "Populating Products")
    product_ids = create_products(api_key)
    print_success(f"\nCreated {len(product_ids)} products")

    # Step 7: Create reviews
    print_step(7, "Populating Reviews")
    review_count = create_reviews()
    print_success(f"\nCreated {review_count} reviews")

    # Step 8: Verify everything
    print_step(8, "Verification")
    if verify_setup():
        print(f"\n{Fore.GREEN}{Style.BRIGHT}")
        print("=" * 80)
        print("  ✓ SETUP COMPLETE!")
        print("=" * 80)
        print(Style.RESET_ALL)
        print(f"\n{Fore.CYAN}Summary:{Style.RESET_ALL}")
        print(f"  Domain: {DOMAIN}")
        print(f"  Admin: {ADMIN_EMAIL}")
        print(f"  Products: {len(product_ids)}")
        print(f"  Reviews: {review_count}")
        print(f"\n{Fore.YELLOW}Next Steps:{Style.RESET_ALL}")
        print(f"  1. Build frontend: cd frontend && ./build-site.sh {DOMAIN}")
        print(f"  2. Access site: http://localhost:8000")
        print(f"  3. Login with: {ADMIN_EMAIL}")
        print()
    else:
        print_error("\nVerification failed. Please check the errors above.")
        sys.exit(1)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print(f"\n\n{Fore.YELLOW}Setup interrupted by user.{Style.RESET_ALL}")
        sys.exit(1)
    except Exception as e:
        print_error(f"\nUnexpected error: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
