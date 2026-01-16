#!/usr/bin/env python3

import requests
import json
import webbrowser
from urllib.parse import urlencode

# Configuration
API_URL = "http://localhost:3000"
CLIENT_ID = "u-s4t2ud-3bdeaf20d556eaee3fc3a50aad024c247bfca90037723b2237099c8b889941b7"
REDIRECT_URI = "http://localhost:3000/auth/42/callback"

def main():
    print("="*70)
    print("  42 OAuth Authentication Test")
    print("="*70)
    print()
    
    # Generate authorization URL
    params = {
        'client_id': CLIENT_ID,
        'redirect_uri': REDIRECT_URI,
        'response_type': 'code'
    }
    
    auth_url = f"https://api.intra.42.fr/oauth/authorize?{urlencode(params)}"
    
    print("Step 1: Opening 42 authorization page in your browser...")
    print()
    print(f"URL: {auth_url}")
    print()
    
    # Open browser
    try:
        webbrowser.open(auth_url)
        print("✓ Browser opened")
    except:
        print("⚠ Could not open browser automatically")
        print("Please copy and paste the URL above into your browser")
    
    print()
    print("Step 2: After authorizing, you'll be redirected to:")
    print(f"        {REDIRECT_URI}?code=AUTHORIZATION_CODE")
    print()
    print("Step 3: Copy the 'code' parameter from the redirected URL")
    print("="*70)
    print()
    
    # Get authorization code from user
    code = input("Enter the authorization code: ").strip()
    
    if not code:
        print("✗ No code provided, exiting")
        return
    
    print()
    print("Authenticating with backend...")
    print()
    
    # Send code to backend
    try:
        response = requests.post(
            f"{API_URL}/auth/42",
            json={'code': code},
            headers={'Content-Type': 'application/json'}
        )
        
        data = response.json()
        
        print(f"Status Code: {response.status_code}")
        print()
        print("Response:")
        print(json.dumps(data, indent=2))
        print()
        
        if response.status_code == 200:
            print("="*70)
            print("✓ 42 Authentication Successful!")
            print("="*70)
            print()
            print(f"JWT Token: {data.get('token', 'N/A')}")
            print()
            print("User Information:")
            user = data.get('user', {})
            print(f"  ID:         {user.get('id', 'N/A')}")
            print(f"  Email:      {user.get('email', 'N/A')}")
            print(f"  First Name: {user.get('firstName', 'N/A')}")
            print(f"  Last Name:  {user.get('lastName', 'N/A')}")
            print()
        else:
            print("="*70)
            print("✗ Authentication Failed")
            print("="*70)
            print(f"Error: {data.get('error', 'Unknown error')}")
            print()
            
    except requests.exceptions.ConnectionError:
        print("✗ Error: Cannot connect to backend server")
        print(f"Make sure the server is running at {API_URL}")
    except Exception as e:
        print(f"✗ Error: {str(e)}")

if __name__ == "__main__":
    main()
