import React, { useState } from "react";

export default function ForgotPassword({
  setOpen,
}: {
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  const [email, setEmail] = useState("");
  const [submitted, setSubmitted] = useState(false);

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black z-50">
      <div className="bg-zinc-900 text-white max-w-md w-full rounded-2xl p-6 relative shadow-xl">
        <button
          className="absolute top-2 right-2 text-white text-2xl font-bold hover:text-red-500"
          onClick={() => setOpen(false)}
        >
          &times;
        </button>
        <h2 className="text-2xl font-semibold text-center mb-4">Forgot Password</h2>

        {submitted ? (
          <div className="text-center space-y-4">
            <p className="text-green-400">
              ✅ A password reset link has been sent to your email.
            </p>
            <button
              onClick={() => setOpen(false)}
              className="w-full bg-red-600 hover:bg-red-700 text-white py-3 rounded font-bold"
            >
              Back to Sign In
            </button>
          </div>
        ) : (
          <div className="space-y-4">
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Enter your email"
              required
              className="w-full p-3 rounded bg-zinc-800 text-white placeholder-zinc-400 focus:outline-none"
            />
            <button
              className="w-full bg-red-600 hover:bg-red-700 text-white py-3 rounded font-bold"
            >
              Send Reset Link
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
