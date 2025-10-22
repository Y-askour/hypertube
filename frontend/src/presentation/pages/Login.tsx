import React, { useState } from "react";
import { FcGoogle } from "react-icons/fc";
import { Si42 } from "react-icons/si";
import ForgotPassword from "./ForgetPassword";
import { useAuth } from "../contexts/authContext";

export default function SignIn({
  setOpen,
}: {
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  const [forgotPassword, setForgotPassword] = useState<boolean>(false);

  // --- States for form inputs and loading ---
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const { loginWithGoogle, loginWith42, login } = useAuth();

  const handleOAuthSignIn = (provider: "google" | "42") => {
    switch (provider) {
      case "google":
        loginWithGoogle();
        break;
      case "42":
        loginWith42();
        break;
      default:
        break;
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await login(username, password);
    } catch (err) {
      console.error("Login failed:", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 flex items-center justify-center 
      bg-gradient-to-br from-black/70 via-zinc-900/80 to-gray-800/70 backdrop-blur-sm 
      z-50">
      <div className="bg-zinc-900 text-white max-w-md w-full rounded-2xl p-6 relative shadow-xl">
        <button
          className="absolute top-2 right-2 text-white text-2xl font-bold hover:text-red-500"
          onClick={() => setOpen(false)}
        >
          &times;
        </button>
        <h2 className="text-2xl font-semibold text-center mb-4">Sign In</h2>

        {/* --- OAuth Buttons --- */}
        <div className="space-y-3 mb-6">
          <button
            onClick={() => handleOAuthSignIn("google")}
            className="flex items-center justify-center w-full gap-3 bg-white text-black py-2 px-4 rounded hover:bg-gray-200 font-semibold"
          >
            Sign in with <FcGoogle size={24} />
          </button>
          <button
            onClick={() => handleOAuthSignIn("42")}
            className="flex items-center justify-center w-full gap-3 bg-gray-800 text-white py-2 px-4 rounded hover:bg-gray-700 font-semibold"
          >
            Sign in with <Si42 size={24} />
          </button>
        </div>

        {/* --- Divider --- */}
        <div className="flex items-center justify-center mb-6">
          <span className="border-t border-zinc-700 w-full"></span>
          <span className="px-3 text-zinc-400 text-sm">or</span>
          <span className="border-t border-zinc-700 w-full"></span>
        </div>

        {/* --- Standard Login Form --- */}
        <form className="space-y-4" onSubmit={handleSubmit}>
          <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="w-full p-3 rounded bg-zinc-800 text-white placeholder-zinc-400 focus:outline-none"
          />
          <div>
            <input
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full p-3 rounded bg-zinc-800 text-white placeholder-zinc-400 focus:outline-none"
            />
            <div className="text-right mt-1">
              <button
                type="button"
                onClick={() => setForgotPassword(true)}
                className="text-sm text-red-400 hover:underline"
              >
                Forgot Password?
              </button>
            </div>
          </div>
          <button
            type="submit"
            disabled={loading}
            className="w-full bg-red-600 hover:bg-red-700 text-white py-3 rounded font-bold transition duration-200 disabled:opacity-50"
          >
            {loading ? "Signing in..." : "Sign In"}
          </button>
        </form>
      </div>
      {forgotPassword && <ForgotPassword setOpen={setForgotPassword} />}
    </div>
  );
}
