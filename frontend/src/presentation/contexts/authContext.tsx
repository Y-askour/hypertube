import React, { createContext, useContext, useState, useCallback } from "react";

interface User {
  username?: string;
  email?: string;
  [key: string]: any;
}

interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (username: string, password: string) => Promise<void>;
  signup : (firstName:string,lastName:string,username:string,password:string,email:string) => Promise<void>;
  loginWithGoogle: () => void;
  loginWith42: () => void;
  resetPassword: (email: string) => Promise<void>;
  logout: () => void;
  loading: boolean;
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  token: null,
  login: async () => {},
  signup:async () => {},
  loginWithGoogle: () => {},
  loginWith42: () => {},
  resetPassword: async () => {},
  logout: () => {},
  loading: false,
});


export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const login = useCallback(async (username: string, password: string) => {
    setLoading(true);
    try {
      console.log("Logging in with:", username, password);
      setUser({ username });
      setToken("dummy-token");
    } finally {
      setLoading(false);
    }
  }, []);

  const signup = useCallback(async (firstName:string,lastName:string,username:string,password:string,email:string) => {
    setLoading(true);
    try {
      console.log("signup with:", firstName,lastName,username,password,email);
      setUser({ username });
      setToken("dummy-token");
    } finally {
      setLoading(false);
    }
  }, []);


  const loginWithGoogle = useCallback(() => {
    console.log("Redirecting to Google login...");
  }, []);

  const loginWith42 = useCallback(() => {
    console.log("Redirecting to 42 login...");
  }, []);

  const resetPassword = useCallback(async (email: string) => {
    console.log("Redirecting to reset password page for:", email);
  }, []);

  const logout = useCallback(() => {
    console.log("Logging out...");
    setUser(null);
    setToken(null);
  }, []);


  const value: AuthContextType = {
    user,
    token,
    login,
    signup,
    loginWithGoogle,
    loginWith42,
    resetPassword,
    logout,
    loading,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};


export const useAuth = () => useContext(AuthContext);
