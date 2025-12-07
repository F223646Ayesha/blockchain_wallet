import { useState } from "react";
import { apiPost } from "../api/api";
import { useNavigate } from "react-router-dom";
import { decryptPrivateKey } from "../crypto/decryptPrivateKey";

export default function Login() {
  const [email, setEmail] = useState("");
  const [otp, setOtp] = useState("");
  const [otpSent, setOtpSent] = useState(false);
  const [message, setMessage] = useState("");

  const navigate = useNavigate();

  // --------------------------------------------------------------
  // 🔹 1. Send OTP to Email
  // --------------------------------------------------------------
  async function handleSendOTP() {
    setMessage("");

    try {
      const res = await apiPost("/auth/send-otp", { email });

      if (!res.success) {
        return setMessage(res.message || "Failed to send OTP");
      }

      setOtpSent(true);
      setMessage("OTP sent to your email!");
    } catch (err) {
      setMessage("Failed to send OTP");
    }
  }

  // --------------------------------------------------------------
  // 🔥 2. Verify OTP & Login
  // --------------------------------------------------------------
  async function handleLogin(e) {
    e.preventDefault();
    setMessage("");

    try {
      const res = await apiPost("/login", { email, otp });

      if (!res.success) {
        return setMessage(res.message || "Invalid OTP");
      }

      const { token, user, wallet } = res.data;

      // Save login session
      localStorage.setItem("token", token);
      localStorage.setItem("user", JSON.stringify(user));
      localStorage.setItem("walletId", wallet.wallet_id);
      localStorage.setItem("publicKey", wallet.public_key);
      localStorage.setItem("privateKeyEnc", wallet.private_key_enc);

      // Decrypt private key
      try {
        const decrypted = decryptPrivateKey(wallet.private_key_enc);
        localStorage.setItem("privateKeyHex", decrypted);
      } catch (error) {
        return setMessage("Private key decryption failed");
      }

      navigate("/");
    } catch (err) {
      setMessage("Login failed");
    }
  }

  // --------------------------------------------------------------
  // UI
  // --------------------------------------------------------------
  return (
    <div className="min-h-screen flex items-center justify-center 
      bg-gradient-to-br from-black via-[#0a0f1f] to-[#020617]">

      <div className="backdrop-blur-xl bg-white/5 border border-white/10
        rounded-2xl p-10 w-full max-w-md shadow-2xl animate-fadeIn">

        <h1 className="text-4xl font-bold text-white text-center mb-2">
          Login to Your Wallet
        </h1>
        <p className="text-gray-400 text-center mb-8">
          Secure OTP-based authentication
        </p>

        {/* Email */}
        <label className="text-gray-300 text-sm">Email</label>
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="w-full mt-1 p-3 rounded-lg bg-black/40 border border-gray-700
                     text-gray-200 focus:ring-2 focus:ring-blue-500 outline-none"
        />

        {/* Send OTP Button */}
        <button
          onClick={handleSendOTP}
          className="w-full mt-4 py-3 rounded-lg bg-blue-600 hover:bg-blue-700 
          text-white font-semibold transition"
        >
          Send OTP
        </button>

        {/* OTP Input */}
        {otpSent && (
          <form onSubmit={handleLogin} className="mt-6 space-y-4">
            <div>
              <label className="text-gray-300 text-sm">Enter OTP</label>
              <input
                type="text"
                value={otp}
                onChange={(e) => setOtp(e.target.value)}
                className="w-full mt-1 p-3 rounded-lg bg-black/40 border border-gray-700
                           text-gray-200 focus:ring-2 focus:ring-green-500 outline-none"
              />
            </div>

            <button
              type="submit"
              className="w-full py-3 rounded-lg bg-green-600 hover:bg-green-700 
              text-white font-semibold transition"
            >
              Login
            </button>
          </form>
        )}

        {/* Message */}
        {message && (
          <div className="mt-6 p-3 rounded-lg bg-blue-900/40 text-blue-300 text-center">
            {message}
          </div>
        )}

        {/* Register Link */}
        <div className="mt-8 text-center">
          <p className="text-gray-400">
            Don’t have an account?
            <button
              onClick={() => navigate("/register")}
              className="text-blue-400 ml-1 hover:text-blue-300 underline underline-offset-4 transition"
            >
              Create one
            </button>
          </p>
        </div>
      </div>
    </div>
  );
}
