#!/usr/bin/env python3
"""
Domain Setup Script
===================

This script sets up a complete domain environment from configuration files including:
- Domain configuration and settings
- Admin user with full permissions
- Product catalog with comprehensive details
- Customer reviews
- Service configurations

Usage:
    python3 setup.py                    # Normal setup
    python3 setup.py --clean            # Clean database first, then setup
    python3 setup.py --clean-only       # Only clean database, don't setup
    python3 setup.py --verify-only      # Only verify existing setup

Requirements:
    Install dependencies: pip install -r requirements.txt

Configuration Files:
    config/domain.json    - Domain settings and branding
    config/admin.json     - Admin user configuration
    config/products.json  - Product catalog
    config/reviews.json   - Customer reviews
    config/services.json  - Service endpoints

Author: Claude
Date: 2026-01-03
"""

import requests
import json
import sys
import argparse
import subprocess
from pathlib import Path
from typing import Dict, List, Optional
from colorama import init, Fore, Style

# Initialize colorama for colored output
init(autoreset=True)

# Base paths
SCRIPT_DIR = Path(__file__).parent
CONFIG_DIR = SCRIPT_DIR / "config"


class SetupScript:
    """Main setup script class."""

    def __init__(self, config_dir: Path = CONFIG_DIR):
        self.config_dir = config_dir
        self.configs = {}
        self.state = {
            'invite_token': None,
            'api_key': None,
            'admin_jwt': None,
            'product_ids': [],
            'review_ids': []
        }

    # ========== Utility Methods ==========

    def print_step(self, step_num: int, description: str):
        """Print a step header."""
        print(f"\n{Fore.CYAN}{'='*80}")
        print(f"{Fore.CYAN}STEP {step_num}: {description}")
        print(f"{Fore.CYAN}{'='*80}{Style.RESET_ALL}")

    def print_success(self, message: str):
        """Print a success message."""
        print(f"{Fore.GREEN}✓ {message}{Style.RESET_ALL}")

    def print_error(self, message: str):
        """Print an error message."""
        print(f"{Fore.RED}✗ {message}{Style.RESET_ALL}")

    def print_info(self, message: str):
        """Print an info message."""
        print(f"{Fore.YELLOW}ℹ {message}{Style.RESET_ALL}")

    def print_data(self, label: str, data: any):
        """Print data with a label."""
        print(f"{Fore.MAGENTA}{label}:{Style.RESET_ALL} {data}")

    # ========== Configuration Loading ==========

    def load_configs(self) -> bool:
        """Load all configuration files."""
        self.print_info("Loading configuration files...")

        config_files = {
            'domain': 'domain.json',
            'admin': 'admin.json',
            'products': 'products.json',
            'reviews': 'reviews.json',
            'services': 'services.json'
        }

        for key, filename in config_files.items():
            filepath = self.config_dir / filename
            if not filepath.exists():
                self.print_error(f"Missing config file: {filename}")
                return False

            try:
                with open(filepath, 'r') as f:
                    self.configs[key] = json.load(f)
                self.print_success(f"Loaded {filename}")
            except json.JSONDecodeError as e:
                self.print_error(f"Invalid JSON in {filename}: {e}")
                return False
            except Exception as e:
                self.print_error(f"Failed to load {filename}: {e}")
                return False

        return True

    # ========== Service Health Checks ==========

    def check_service_health(self, service_name: str, url: str) -> bool:
        """Check if a service is running and healthy."""
        try:
            response = requests.get(f"{url}/health", timeout=5)
            if response.status_code == 200:
                self.print_success(f"{service_name} is healthy")
                return True
            else:
                self.print_error(f"{service_name} returned status {response.status_code}")
                return False
        except requests.exceptions.RequestException as e:
            self.print_error(f"{service_name} is not accessible: {e}")
            return False

    # ========== Database Cleaning ==========

    def clean_database(self) -> bool:
        """Clean all data for this domain from the database."""
        self.print_info("Cleaning database...")

        domain = self.configs['domain']['domain']
        services = self.configs['services']
        db_config = services['database']

        commands = [
            # Clean auth_module database
            {
                'description': 'domains collection',
                'command': f'db.domains.deleteOne({{domain: "{domain}"}})'
            },
            {
                'description': 'users collection',
                'command': f'db.users.deleteMany({{domain: "{domain}"}})'
            },
            {
                'description': 'invitations collection',
                'command': f'db.invitations.deleteMany({{domain: "{domain}"}})'
            },
            {
                'description': 'api_keys collection',
                'command': f'db.api_keys.deleteMany({{domain: "{domain}"}})'
            },
            # Clean products_module database
            {
                'description': 'products collection',
                'command': f'db.products.deleteMany({{domain: "{domain}"}})',
                'database': 'products_module'
            },
            {
                'description': 'reviews collection',
                'command': f'db.reviews.deleteMany({{domain: "{domain}"}})',
                'database': 'products_module'
            },
            {
                'description': 'contacts collection',
                'command': f'db.contacts.deleteMany({{domain: "{domain}"}})',
                'database': 'products_module'
            }
        ]

        success = True
        for cmd_info in commands:
            database = cmd_info.get('database', 'auth_module')
            mongo_uri = f"mongodb://{db_config['host']}:{db_config['port']}/{database}"

            try:
                result = subprocess.run(
                    ['docker', 'exec', 'auth-mongodb', 'mongosh', mongo_uri,
                     '--eval', cmd_info['command'], '--quiet'],
                    capture_output=True,
                    text=True,
                    timeout=10
                )

                if result.returncode == 0:
                    # Parse result to see how many were deleted
                    if 'deletedCount' in result.stdout:
                        self.print_success(f"Cleaned {cmd_info['description']}")
                    else:
                        self.print_info(f"Checked {cmd_info['description']}")
                else:
                    self.print_error(f"Failed to clean {cmd_info['description']}: {result.stderr}")
                    success = False
            except Exception as e:
                self.print_error(f"Error cleaning {cmd_info['description']}: {e}")
                success = False

        return success

    # ========== Domain Setup ==========

    def create_domain(self) -> bool:
        """Create the domain using CLI and extract invitation token."""
        domain_config = self.configs['domain']
        admin_config = self.configs['admin']
        domain = domain_config['domain']

        self.print_info(f"Creating domain: {domain}")

        try:
            # Use CLI inside Docker container to create domain (with invitation for admin user)
            result = subprocess.run(
                [
                    "docker", "exec", "auth-api",
                    "./auth-module", "domain", "create",
                    "--domain", domain,
                    "--name", domain_config['name'],
                    "--admin-email", admin_config['email']
                ],
                capture_output=True,
                text=True,
                timeout=30
            )

            if result.returncode == 0:
                self.print_success(f"Domain created: {domain}")
                self.print_data("  Name", domain_config['name'])

                # Extract invitation token from output (check both stdout and stderr)
                combined_output = result.stdout + "\n" + result.stderr
                for line in combined_output.split('\n'):
                    if 'token=' in line:
                        # Extract token from URL like: http://localhost:3000/invite?token=ABC123
                        token = line.split('token=')[1].strip()
                        self.state['invite_token'] = token
                        self.print_success("Invitation token extracted")
                        self.print_data("  Token", token[:40] + "...")
                        break

                return True
            elif "already exists" in result.stderr.lower() or "already exists" in result.stdout.lower():
                self.print_info("Domain already exists, continuing...")
                # Domain exists but we need to create an invite token separately
                return True
            else:
                self.print_error(f"Failed to create domain: {result.stderr}")
                return False
        except subprocess.TimeoutExpired:
            self.print_error("Domain creation timed out")
            return False
        except Exception as e:
            self.print_error(f"Failed to create domain: {e}")
            return False

    def create_invite_token(self) -> Optional[str]:
        """Create an invite token for the domain."""
        domain = self.configs['domain']['domain']
        self.print_info(f"Creating invite token for {domain}")

        try:
            result = subprocess.run(
                ["docker", "exec", "auth-api", "./auth-module", "domain", "create-invite", domain],
                capture_output=True,
                text=True,
                timeout=10
            )

            if result.returncode == 0:
                output = result.stdout.strip()
                # Try to extract token from output
                for line in output.split('\n'):
                    line = line.strip()
                    if len(line) > 30 and '=' in line:
                        token = line.split('=')[-1].strip().strip('"\'')
                        if len(token) > 20:
                            self.print_success("Invite token created")
                            self.print_data("  Token", token[:30] + "...")
                            return token

                # If parsing failed, show output
                self.print_info("Could not parse token. Output:")
                print(output)
                self.print_info("Please paste the invite token:")
                return input().strip()
            else:
                self.print_error(f"Failed to create invite token: {result.stderr}")
                return None
        except Exception as e:
            self.print_error(f"Failed to run command: {e}")
            return None

    # ========== User Setup ==========

    def create_admin_user(self) -> bool:
        """Create the admin user via invitation accept."""
        admin_config = self.configs['admin']
        services = self.configs['services']

        self.print_info(f"Creating admin user: {admin_config['email']}")

        url = f"{services['auth_service']['url']}{services['auth_service']['endpoints']['signup']}"

        # Accept invitation to create user (auth system uses magic_link, not passwords)
        payload = {
            'token': self.state['invite_token'],
            'email': admin_config['email'],
            'auth_provider': 'magic_link'
        }

        try:
            response = requests.post(url, json=payload, headers={"Content-Type": "application/json"})

            if response.status_code in [200, 201]:
                data = response.json()
                self.print_success(f"Admin user created: {admin_config['email']}")
                # Store JWT for later use
                if data.get('token'):
                    self.state['admin_jwt'] = data['token']
                    self.print_data("  JWT Token", "received")
                if data.get('user', {}).get('id'):
                    self.print_data("  User ID", data['user']['id'])
                return True
            elif response.status_code == 409 or "already" in response.text.lower():
                self.print_info("User already exists, continuing...")
                return True
            else:
                self.print_error(f"Failed to create user: {response.text}")
                return False
        except requests.exceptions.RequestException as e:
            self.print_error(f"Request failed: {e}")
            return False

    def login_admin(self) -> Optional[str]:
        """Login as admin and get JWT token."""
        admin_config = self.configs['admin']
        services = self.configs['services']

        self.print_info(f"Logging in as: {admin_config['email']}")

        url = f"{services['auth_service']['url']}{services['auth_service']['endpoints']['login']}"

        payload = {
            'email': admin_config['email'],
            'password': admin_config['password']
        }

        try:
            response = requests.post(url, json=payload, headers={"Content-Type": "application/json"})

            if response.status_code == 200:
                data = response.json()
                token = data.get('token')
                if token:
                    self.print_success("Login successful")
                    self.print_data("  JWT Token", token[:30] + "...")
                    return token
                else:
                    self.print_error("No token in response")
                    return None
            else:
                self.print_error(f"Login failed: {response.text}")
                return None
        except requests.exceptions.RequestException as e:
            self.print_error(f"Request failed: {e}")
            return None

    def create_api_key(self) -> Optional[str]:
        """Create API key for products service."""
        domain = self.configs['domain']['domain']
        api_key_config = self.configs['services']['api_keys']['products_service']

        self.print_info("Creating API key for products service")

        permissions = ','.join(api_key_config['permissions'])

        try:
            result = subprocess.run(
                [
                    "docker", "exec", "auth-api",
                    "./auth-module", "apikey", "create",
                    "--domain", domain,
                    "--service", "products",
                    "--permissions", permissions
                ],
                capture_output=True,
                text=True,
                timeout=10
            )

            if result.returncode == 0:
                output = result.stdout + "\n" + result.stderr
                # Try to extract API key from output
                for line in output.split('\n'):
                    line = line.strip()
                    # Look for line with "API Key: " prefix
                    if line.startswith('API Key:'):
                        api_key = line.split('API Key:')[1].strip()
                        self.print_success("API key created")
                        self.print_data("  API Key", api_key[:40] + "...")
                        return api_key
                    # Or JWT tokens that start with 'ey'
                    elif line.startswith('ey') and len(line) > 50 and '.' in line:
                        self.print_success("API key created")
                        self.print_data("  API Key", line[:40] + "...")
                        return line

                # If parsing failed
                self.print_info("Could not parse API key. Output:")
                print(output)
                self.print_info("Please paste the API key:")
                return input().strip()
            else:
                self.print_error(f"Failed to create API key: {result.stderr}")
                return None
        except Exception as e:
            self.print_error(f"Failed to run command: {e}")
            return None

    # ========== Products Setup ==========

    def create_products(self) -> List[str]:
        """Create all products and return their IDs."""
        products = self.configs['products']
        domain = self.configs['domain']['domain']
        services = self.configs['services']

        self.print_info(f"Creating {len(products)} products for {domain}")

        url = f"{services['products_service']['url']}{services['products_service']['endpoints']['products']}"
        headers = {
            "Authorization": f"Bearer {self.state['api_key']}",
            "Content-Type": "application/json"
        }

        product_ids = []

        for idx, product in enumerate(products, 1):
            print(f"\n  [{idx}/{len(products)}] Creating: {product['name']}")

            # Map product to API format (CreateProductRequest)
            # Convert all attribute values to strings (API expects map[string]string)
            attributes = {k: str(v) for k, v in product.get('attributes', {}).items()}

            api_product = {
                'name': product['name'],
                'description': product['description'],
                'base_price': product['price'],
                'images': [img['url'] for img in product.get('images', [])],
                'attributes': attributes,
                'variants': [
                    {
                        'attributes': {k: str(v) for k, v in v.get('attributes', {}).items()},
                        'price': v.get('price', product['price']),
                        'stock': v.get('stock', 0),
                        'sku': v.get('sku', ''),
                        'image_index': 0
                    }
                    for v in product.get('variants', [])
                ],
            }

            try:
                response = requests.post(url, json=api_product, headers=headers)

                if response.status_code in [200, 201]:
                    data = response.json()
                    product_id = data.get('id')
                    product_ids.append(product_id)
                    self.print_success(f"Created: {product['name']}")
                    self.print_data("    ID", product_id)
                    self.print_data("    SKU", product.get('sku'))
                    self.print_data("    Price", f"${product['price']}")
                    if product.get('bestseller'):
                        self.print_data("    Bestseller", "Yes")
                else:
                    self.print_error(f"Failed to create product: {response.text}")
            except requests.exceptions.RequestException as e:
                self.print_error(f"Request failed: {e}")

        return product_ids

    # ========== Reviews Setup ==========

    def create_reviews(self) -> List[str]:
        """Create all reviews and return their IDs."""
        reviews = self.configs['reviews']
        domain = self.configs['domain']['domain']
        services = self.configs['services']

        self.print_info(f"Creating {len(reviews)} reviews for {domain}")

        # Use public endpoint for reviews
        url_template = services['products_service']['endpoints']['reviews']
        url = f"{services['products_service']['url']}{url_template.replace('{domain}', domain)}"

        headers = {"Content-Type": "application/json"}

        review_ids = []

        for idx, review in enumerate(reviews, 1):
            print(f"\n  [{idx}/{len(reviews)}] Creating review by: {review['name']}")

            # Prepare review (only fields accepted by API)
            payload = {
                'product': review['product'],
                'name': review['name'],
                'rating': review['rating'],
                'text': review['text'],
                'highlight': review.get('highlight', '')
            }

            try:
                response = requests.post(url, json=payload, headers=headers)

                if response.status_code in [200, 201]:
                    data = response.json()
                    review_id = data.get('id')
                    review_ids.append(review_id)
                    self.print_success(f"Created review by {review['name']}")
                    self.print_data("    Product", review['product'])
                    self.print_data("    Rating", f"{'⭐' * review['rating']}")
                    self.print_data("    ID", review_id)
                else:
                    self.print_error(f"Failed to create review: {response.text}")
            except requests.exceptions.RequestException as e:
                self.print_error(f"Request failed: {e}")

        return review_ids

    # ========== Frontend Setup ==========

    def create_branding_config(self) -> bool:
        """Create branding.json for frontend."""
        domain_config = self.configs['domain']
        domain = domain_config['domain']

        self.print_info(f"Creating branding configuration for {domain}")

        # Path to frontend sites directory
        frontend_dir = Path("/home/sparque/dev/oilyourhairapp/frontend")
        site_dir = frontend_dir / "sites" / domain
        site_dir.mkdir(parents=True, exist_ok=True)

        # Create branding.json from domain config
        branding = {
            "siteName": domain_config['name'],
            "brandName": domain_config['name'].split()[0],  # First word
            "tagline": domain_config.get('description', 'Your trusted store'),
            "primaryColor": domain_config.get('branding', {}).get('primary_color', '#2E7D32'),
            "primaryColorLight": domain_config.get('branding', {}).get('secondary_color', '#4CAF50'),
            "logo": domain_config.get('branding', {}).get('logo_url', f'https://via.placeholder.com/150x50?text={domain_config["name"]}'),
            "domain": domain,
            "apiUrl": self.configs['services']['products_service']['url'],
            "authUrl": self.configs['services']['auth_service']['url']
        }

        branding_file = site_dir / "branding.json"

        try:
            with open(branding_file, 'w') as f:
                json.dump(branding, f, indent=2)

            self.print_success(f"Branding config created: {branding_file}")
            self.print_data("  Primary Color", branding['primaryColor'])
            self.print_data("  Site Name", branding['siteName'])
            return True
        except Exception as e:
            self.print_error(f"Failed to create branding config: {e}")
            return False

    def build_frontend(self) -> bool:
        """Build the frontend site from template."""
        domain = self.configs['domain']['domain']

        self.print_info(f"Building frontend for {domain}")

        frontend_dir = Path("/home/sparque/dev/oilyourhairapp/frontend")
        build_script = frontend_dir / "build-site.sh"

        if not build_script.exists():
            self.print_error(f"Build script not found: {build_script}")
            return False

        try:
            result = subprocess.run(
                ["./build-site.sh", domain],
                cwd=str(frontend_dir),
                capture_output=True,
                text=True,
                timeout=30
            )

            if result.returncode == 0:
                self.print_success("Frontend built successfully")

                # Check output directory
                output_dir = frontend_dir / "sites" / domain / "public"
                if output_dir.exists():
                    file_count = len(list(output_dir.glob("**/*")))
                    self.print_data("  Output", str(output_dir))
                    self.print_data("  Files", file_count)
                return True
            else:
                self.print_error(f"Build failed: {result.stderr}")
                return False
        except subprocess.TimeoutExpired:
            self.print_error("Build timed out after 30 seconds")
            return False
        except Exception as e:
            self.print_error(f"Failed to build frontend: {e}")
            return False

    # ========== Verification ==========

    def verify_setup(self) -> bool:
        """Verify the complete setup."""
        self.print_info("Verifying setup...")

        domain = self.configs['domain']['domain']
        services = self.configs['services']
        all_ok = True

        # Verify products
        try:
            url_template = services['products_service']['endpoints']['public_products']
            url = f"{services['products_service']['url']}{url_template.replace('{domain}', domain)}"

            response = requests.get(url)
            if response.status_code == 200:
                data = response.json()
                product_count = data.get('count', 0)
                self.print_success(f"Products verified: {product_count} products found")

                products = data.get('products', [])
                bestsellers = [p for p in products if p.get('bestseller')]
                self.print_data("  Bestsellers", len(bestsellers))

                active = [p for p in products if p.get('active')]
                self.print_data("  Active", len(active))
            else:
                self.print_error(f"Failed to verify products: {response.text}")
                all_ok = False
        except Exception as e:
            self.print_error(f"Products verification failed: {e}")
            all_ok = False

        # Verify reviews
        try:
            url_template = services['products_service']['endpoints']['reviews']
            url = f"{services['products_service']['url']}{url_template.replace('{domain}', domain)}"

            response = requests.get(url)
            if response.status_code == 200:
                data = response.json()
                review_count = data.get('count', 0)
                self.print_success(f"Reviews verified: {review_count} reviews found")

                reviews = data.get('reviews', [])
                five_star = [r for r in reviews if r.get('rating') == 5]
                self.print_data("  5-star reviews", len(five_star))

                avg_rating = sum(r.get('rating', 0) for r in reviews) / len(reviews) if reviews else 0
                self.print_data("  Average rating", f"{avg_rating:.1f}")
            else:
                self.print_error(f"Failed to verify reviews: {response.text}")
                all_ok = False
        except Exception as e:
            self.print_error(f"Reviews verification failed: {e}")
            all_ok = False

        # Verify admin login (skip if already authenticated during setup)
        if self.state.get('admin_jwt'):
            self.print_success("Admin authentication verified (JWT from setup)")
        else:
            try:
                jwt = self.login_admin()
                if jwt:
                    self.print_success("Admin authentication verified")
                else:
                    self.print_error("Admin authentication failed")
                    all_ok = False
            except Exception as e:
                self.print_error(f"Auth verification failed: {e}")
                all_ok = False

        return all_ok

    # ========== Main Execution ==========

    def run(self, clean: bool = False, clean_only: bool = False, verify_only: bool = False):
        """Main execution method."""
        print(f"\n{Fore.CYAN}{Style.BRIGHT}")
        print("=" * 80)
        print("  Domain Setup Script")
        print(f"  Domain: {self.configs['domain']['domain']}")
        print("=" * 80)
        print(Style.RESET_ALL)

        # Clean database if requested
        if clean or clean_only:
            self.print_step(0, "Cleaning Database")
            if not self.clean_database():
                self.print_error("Failed to clean database")
                if not clean_only:
                    sys.exit(1)

            if clean_only:
                self.print_success("Database cleaned successfully")
                return

        # Verify only mode
        if verify_only:
            self.print_step(0, "Verification Only")
            if self.verify_setup():
                self.print_success("Verification passed!")
            else:
                self.print_error("Verification failed")
                sys.exit(1)
            return

        # Step 0: Check service health
        self.print_step(0, "Checking Service Health")

        services = self.configs['services']
        auth_ok = self.check_service_health("Auth Service", services['auth_service']['url'])
        products_ok = self.check_service_health("Products Service", services['products_service']['url'])

        if not auth_ok or not products_ok:
            self.print_error("One or more services are not running!")
            sys.exit(1)

        # Step 1: Create domain
        self.print_step(1, "Creating Domain")
        if not self.create_domain():
            self.print_error("Failed to create domain")
            sys.exit(1)

        # Step 2: Create invite token (if not already created during domain creation)
        if not self.state.get('invite_token'):
            self.print_step(2, "Creating Invite Token")
            self.state['invite_token'] = self.create_invite_token()
            if not self.state['invite_token']:
                self.print_error("Failed to create invite token")
                sys.exit(1)
        else:
            self.print_step(2, "Invite Token (already created)")

        # Step 3: Create admin user
        self.print_step(3, "Creating Admin User")
        if not self.create_admin_user():
            self.print_error("Failed to create admin user")
            sys.exit(1)

        # Step 4: Login and get JWT (if not already obtained during user creation)
        if not self.state.get('admin_jwt'):
            self.print_step(4, "Authenticating Admin User")
            self.state['admin_jwt'] = self.login_admin()
            if not self.state['admin_jwt']:
                self.print_error("Failed to login")
                sys.exit(1)
        else:
            self.print_step(4, "JWT Token (already obtained)")

        # Step 5: Create API key
        self.print_step(5, "Creating API Key for Products")
        self.state['api_key'] = self.create_api_key()
        if not self.state['api_key']:
            self.print_error("Failed to create API key")
            sys.exit(1)

        # Step 6: Create products
        self.print_step(6, "Populating Products")
        self.state['product_ids'] = self.create_products()
        self.print_success(f"\nCreated {len(self.state['product_ids'])} products")

        # Step 7: Create reviews
        self.print_step(7, "Populating Reviews")
        self.state['review_ids'] = self.create_reviews()
        self.print_success(f"\nCreated {len(self.state['review_ids'])} reviews")

        # Step 8: Create frontend branding
        self.print_step(8, "Creating Frontend Branding")
        if not self.create_branding_config():
            self.print_error("Failed to create branding config")
            sys.exit(1)

        # Step 9: Build frontend
        self.print_step(9, "Building Frontend from Template")
        if not self.build_frontend():
            self.print_error("Failed to build frontend")
            sys.exit(1)

        # Step 10: Verify everything
        self.print_step(10, "Verification")
        if self.verify_setup():
            print(f"\n{Fore.GREEN}{Style.BRIGHT}")
            print("=" * 80)
            print("  ✓ SETUP COMPLETE!")
            print("=" * 80)
            print(Style.RESET_ALL)

            domain = self.configs['domain']['domain']
            admin = self.configs['admin']

            print(f"\n{Fore.CYAN}Summary:{Style.RESET_ALL}")
            print(f"  Domain: {domain}")
            print(f"  Admin: {admin['email']}")
            print(f"  Password: {admin['password']}")
            print(f"  Products: {len(self.state['product_ids'])}")
            print(f"  Reviews: {len(self.state['review_ids'])}")
            print(f"  Frontend: ✓ Built from template1")

            print(f"\n{Fore.YELLOW}Next Steps:{Style.RESET_ALL}")
            frontend_path = f"/home/sparque/dev/oilyourhairapp/frontend/sites/{domain}/public"
            print(f"  1. Serve site: cd {frontend_path} && python3 -m http.server 8000")
            print(f"  2. Access: http://localhost:8000")
            print(f"  3. Login with: {admin['email']}")
            print()
        else:
            self.print_error("Verification failed")
            sys.exit(1)


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description='Setup domain environment from configuration files',
        formatter_class=argparse.RawDescriptionHelpFormatter
    )

    parser.add_argument(
        '--clean',
        action='store_true',
        help='Clean database before setup'
    )

    parser.add_argument(
        '--clean-only',
        action='store_true',
        help='Only clean database, do not setup'
    )

    parser.add_argument(
        '--verify-only',
        action='store_true',
        help='Only verify existing setup'
    )

    args = parser.parse_args()

    try:
        # Initialize setup script
        setup = SetupScript()

        # Load configurations
        if not setup.load_configs():
            sys.exit(1)

        # Run setup
        setup.run(
            clean=args.clean,
            clean_only=args.clean_only,
            verify_only=args.verify_only
        )

    except KeyboardInterrupt:
        print(f"\n\n{Fore.YELLOW}Setup interrupted by user.{Style.RESET_ALL}")
        sys.exit(1)
    except Exception as e:
        print(f"\n{Fore.RED}Unexpected error: {e}{Style.RESET_ALL}")
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()
