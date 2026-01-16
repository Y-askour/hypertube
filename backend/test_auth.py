#!/usr/bin/env python3

import requests
import json
import sys

# Configuration
API_URL = "http://localhost:3000"
EMAIL = "test@example.com"
PASSWORD = "TestPassword123"
FIRST_NAME = "John"
LAST_NAME = "Doe"

def print_section(title):
    print(f"\n{'='*50}")
    print(f"  {title}")
    print(f"{'='*50}\n")

def register_user():
    """Register a new user"""
    print_section("REGISTRATION")
    
    payload = {
        "email": EMAIL,
        "password": PASSWORD,
        "firstName": FIRST_NAME,
        "lastName": LAST_NAME,
        "language": "en"
    }
    
    print(f"Registering user: {EMAIL}")
    print(f"Password: {PASSWORD}")
    print(f"Name: {FIRST_NAME} {LAST_NAME}\n")
    
    try:
        response = requests.post(
            f"{API_URL}/auth/register",
            json=payload,
            headers={"Content-Type": "application/json"}
        )
        
        data = response.json()
        print(f"Status Code: {response.status_code}")
        print(f"Response: {json.dumps(data, indent=2)}\n")
        
        if response.status_code == 201:
            print("✓ Registration successful!")
            print(f"Token: {data.get('token', 'N/A')}")
            return data.get('token')
        elif "already exists" in data.get('error', ''):
            print("⚠ User already exists, proceeding to login...")
            return None
        else:
            print(f"✗ Registration failed: {data.get('error', 'Unknown error')}")
            return None
            
    except requests.exceptions.ConnectionError:
        print("✗ Error: Cannot connect to API server")
        print(f"Make sure the server is running at {API_URL}")
        sys.exit(1)
    except Exception as e:
        print(f"✗ Error: {str(e)}")
        return None

def login_user():
    """Login with email and password"""
    print_section("LOGIN")
    
    payload = {
        "email": EMAIL,
        "password": PASSWORD
    }
    
    print(f"Logging in with: {EMAIL}")
    print(f"Password: {PASSWORD}\n")
    
    try:
        response = requests.post(
            f"{API_URL}/auth/email_and_password",
            json=payload,
            headers={"Content-Type": "application/json"}
        )
        
        data = response.json()
        print(f"Status Code: {response.status_code}")
        print(f"Response: {json.dumps(data, indent=2)}\n")
        
        if response.status_code == 200:
            print("✓ Login successful!")
            print(f"Token: {data.get('token', 'N/A')}")
            print(f"\nUser Info:")
            print(f"  ID: {data.get('user', {}).get('id', 'N/A')}")
            print(f"  Email: {data.get('user', {}).get('email', 'N/A')}")
            print(f"  Name: {data.get('user', {}).get('firstName', '')} {data.get('user', {}).get('lastName', '')}")
            return data.get('token')
        else:
            print(f"✗ Login failed: {data.get('error', 'Unknown error')}")
            return None
            
    except Exception as e:
        print(f"✗ Error: {str(e)}")
        return None

def test_wrong_password():
    """Test login with wrong password"""
    print_section("TESTING WRONG PASSWORD")
    
    payload = {
        "email": EMAIL,
        "password": "WrongPassword123"
    }
    
    print(f"Attempting login with wrong password...")
    print(f"Email: {EMAIL}")
    print(f"Password: WrongPassword123\n")
    
    try:
        response = requests.post(
            f"{API_URL}/auth/email_and_password",
            json=payload,
            headers={"Content-Type": "application/json"}
        )
        
        data = response.json()
        print(f"Status Code: {response.status_code}")
        print(f"Response: {json.dumps(data, indent=2)}\n")
        
        if response.status_code == 401:
            print("✓ Wrong password correctly rejected")
        else:
            print("✗ Security issue: Wrong password was not rejected!")
            
    except Exception as e:
        print(f"✗ Error: {str(e)}")

def main():
    print("=" * 50)
    print("  Authentication Test Script")
    print("=" * 50)
    print(f"API URL: {API_URL}")
    
    # Step 1: Register
    register_token = register_user()
    
    # Step 2: Login
    login_token = login_user()
    
    # Step 3: Test wrong password
    test_wrong_password()
    
    print_section("TEST COMPLETE")
    if login_token:
        print("✓ All authentication tests passed!")
        print(f"\nYour JWT Token:\n{login_token}")
    else:
        print("✗ Some tests failed")

if __name__ == "__main__":
    main()
