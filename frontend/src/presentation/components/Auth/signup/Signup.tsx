import React from "react";
import { FcGoogle } from "react-icons/fc";
import { Si42 } from "react-icons/si";
import { useAuth } from "../../../contexts/authContext";
export default function SignUp({
  setOpen,
}: {
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
}) {

  const {loginWithGoogle,loginWith42,signup} = useAuth();

  const handleOAuthSignUp = (provider: "google" | "42") => {
      switch (provider) {
        case 'google':
          loginWithGoogle();
          break;
        case '42' : 
          loginWith42();
          break;
        default:
          break;
      }
  };

  const handleSubmit = (e)=>{
    // const {firstName,lastName,username,password,email}
    // signup()
    e.preventDefault();
    console.log(e);
  }
  return (
    <div className="flex justify-center items-center h-screen">
      <div className="fixed inset-0 flex items-center justify-center
       bg-gradient-to-br from-black/70 via-zinc-900/80 to-gray-800/70 backdrop-blur-sm z-50">
        <div className="bg-zinc-900 text-white max-w-md w-full rounded-2xl p-6 relative">
          <button
            className="absolute top-2 right-2 text-white text-xl font-bold hover:text-red-400"
            onClick={() => setOpen(false)}
          >
            &times;
          </button>
          <h2 className="text-2xl font-semibold text-center mb-4">
            Create Account
          </h2>

          {/* --- OAuth Buttons --- */}
          <div className="space-y-3 mb-6">
            <button
              onClick={() => handleOAuthSignUp("google")}
              className="flex items-center justify-center w-full gap-3 bg-white text-black py-2 px-4 rounded hover:bg-gray-200 font-semibold"
            >
              Sign up with 
              <FcGoogle size={24} />
            </button>
            <button
              onClick={() => handleOAuthSignUp("42")}
              className="flex items-center justify-center w-full gap-3 bg-gray-800 text-white py-2 px-4 rounded hover:bg-gray-700 font-semibold"
            >
               Sign up with  <Si42 size={24} />
             
            </button>
          </div>

          {/* --- Divider --- */}
          <div className="flex items-center justify-center mb-6">
            <span className="border-t border-zinc-700 w-full"></span>
            <span className="px-3 text-zinc-400 text-sm">or</span>
            <span className="border-t border-zinc-700 w-full"></span>
          </div>

          {/* --- Manual Form --- */}
          <form className="space-y-4" onSubmit={handleSubmit}>
            <input
              type="text"
              placeholder="First Name"
              className="w-full p-3 rounded bg-zinc-800 text-white placeholder-zinc-400 focus:outline-none"
            />
            <input
              type="text"
              placeholder="Last Name"
              className="w-full p-3 rounded bg-zinc-800 text-white placeholder-zinc-400 focus:outline-none"
            />
            <input
              type="text"
              placeholder="Username"
              className="w-full p-3 rounded bg-zinc-800 text-white placeholder-zinc-400 focus:outline-none"
            />
            <input
              type="password"
              placeholder="Password"
              className="w-full p-3 rounded bg-zinc-800 text-white placeholder-zinc-400 focus:outline-none"
            />
            <input
              type="email"
              placeholder="Email Address"
              className="w-full p-3 rounded bg-zinc-800 text-white placeholder-zinc-400 focus:outline-none"
            />
            <button
              type="submit"
              className="w-full bg-red-600 hover:bg-red-700 text-white py-3 rounded font-bold"
            >
              Register
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
