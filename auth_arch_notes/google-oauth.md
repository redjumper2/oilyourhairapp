 Update(auth_module/API.md)
  ⎿  Added 57 lines                                                                                                                                 
      114  
      115  ---
      116  
      117 +### Google OAuth Login                                                                                                                   
      118 +                                                                                                                                         
      119 +Initiate Google OAuth flow.                                                                                                              
      120 +                                                                                                                                         
      121 +**Endpoint:** `GET /auth/google`                                                                                                         
      122 +                                                                                                                                         
      123 +**Headers:**                                                                                                                             
      124 +- `Host: oilyourhair.com` (domain context)                                                                                               
      125 +                                                                                                                                         
      126 +**Flow:**                                                                                                                                
      127 +1. User clicks "Sign in with Google"                                                                                                     
      128 +2. Frontend redirects to: `GET /api/v1/auth/google`                                                                                      
      129 +3. User is redirected to Google for authentication                                                                                       
      130 +4. After approval, Google redirects back to callback URL                                                                                 
      131 +                                                                                                                                         
      132 +---                                                                                                                                      
      133 +                                                                                                                                         
      134 +### Google OAuth Callback                                                                                                                
      135 +                                                                                                                                         
      136 +Handle Google OAuth callback (browser redirect).                                                                                         
      137 +                                                                                                                                         
      138 +**Endpoint:** `GET /auth/google/callback`                                                                                                
      139 +                                                                                                                                         
      140 +**Note:** This is called by Google, not your frontend directly.                                                                          
      141 +                                                                                                                                         
      142 +**Response:** Redirects to frontend with JWT                                                                                             
      143 +```                                                                                                                                      
      144 +HTTP 302 Redirect                                                                                                                        
      145 +Location: http://localhost:3000/auth/callback?token=eyJhbGc...                                                                           
      146 +```                                                                                                                                      
      147 +                                                                                                                                         
      148 +Frontend should extract token from URL and store it.                                                                                     
      149 +                                                                                                                                         
      150 +---                                                                                                                                      
      151 +                                                                                                                                         
      152 +### Google OAuth Callback (JSON)                                                                                                         
      153 +                                                                                                                                         
      154 +Handle Google OAuth callback and return JSON (for API clients).                                                                          
      155 +                                                                                                                                         
      156 +**Endpoint:** `GET /auth/google/callback/json?domain=<domain>`                                                                           
      157 +                                                                                                                                         
      158 +**Response:** `200 OK`                                                                                                                   
      159 +```json                                                                                                                                  
      160 +{                                                                                                                                        
      161 +  "token": "eyJhbGc...",                                                                                                                 
      162 +  "user": {                                                                                                                              
      163 +    "id": "507f1f77bcf86cd799439011",                                                                                                    
      164 +    "email": "user@gmail.com",                                                                                                           
      165 +    "domain": "oilyourhair.com",                                                                                                         
      166 +    "role": "customer",                                                                                                                  
      167 +    "permissions": ["products.read", "cart.read", "cart.write"]                                                                          
      168 +  }                                                                                                                                      
      169 +}                                                                                                                                        
      170 +```                                                                                                                                      
      171 +                                                                                                                                         
      172 +---                                                                                                                                      
      173 +                                                                                                                                         
      174  ### Get Current User
      175  
      176  Get authenticated user information.

● Update(auth_module/API.md)
  ⎿  Added 16 lines                                                                                                                                 
      423  5. User accepts → Frontend: `POST /auth/invitation/accept`
      424  6. Frontend receives JWT, user is logged in
      425  
      426 +### Google OAuth Flow                                                                                                                    
      427 +                                                                                                                                         
      428 +1. User clicks "Sign in with Google" on `oilyourhair.com`                                                                                
      429 +2. Frontend: `window.location.href = "http://localhost:8080/api/v1/auth/google"` with `Host: oilyourhair.com` header                     
      430 +3. Backend stores domain in session                                                                                                      
      431 +4. User redirected to Google for authentication                                                                                          
      432 +5. User approves permissions                                                                                                             
      433 +6. Google redirects to: `GET /auth/google/callback`                                                                                      
      434 +7. Backend verifies OAuth, creates/logs in user                                                                                          
      435 +8. Backend redirects to: `http://localhost:3000/auth/callback?token=<jwt>`                                                               
      436 +9. Frontend extracts token from URL, stores it                                                                                           
      437 +10. All subsequent requests include `Authorization: Bearer <jwt>`                                                                        
      438 +                                                                                                                                         
      439 +**Alternative (API clients):**                                                                                                           
      440 +Use `/auth/google/callback/json` to get JSON response instead of redirect                                                                
      441 +                                                                                                                                         
      442  ### Domain Isolation
      443  
      444  JWT tokens are domain-scoped. If user has JWT for `oilyourhair.com`, they cannot use it on `fitnessclub.com`.

● Write(auth_module/GOOGLE_OAUTH_SETUP.md)
  ⎿  Wrote 234 lines to auth_module/GOOGLE_OAUTH_SETUP.md                  
      1 # Google OAuth Setup Guide
      2 This guide shows how to configure Google OAuth for the auth module.
      3 ## Step 1: Create Google Cloud Project
      4 1. Go to [Google Cloud Console](https://console.cloud.google.com/)
      5 2. Create a new project or select existing one
      6 3. Note your project name
     … +224 lines (ctrl+o to expand)