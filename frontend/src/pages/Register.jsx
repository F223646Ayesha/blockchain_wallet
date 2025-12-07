import { useState } from "react";
import { apiPost } from "../api/api";
import { useNavigate } from "react-router-dom";
import { decryptPrivateKey } from "../crypto/decryptPrivateKey";

export default function Register() {
  const navigate = useNavigate();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [cnic, setCnic] = useState("");

  const [otpSent, setOtpSent] = useState(false);
  const [otp, setOtp] = useState("");

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  // ---------------------------------------------------------
  // SEND OTP
  // ---------------------------------------------------------
  async function sendOtp() {
    setError("");
    setSuccess("");
    setLoading(true);

    try {
      const res = await apiPost("/auth/send-otp", { email });

      if (!res.success) {
        setError(res.message || "Failed to send OTP");
        setLoading(false);
        return;
      }

      setOtpSent(true);
      setSuccess("OTP sent to your email!");
    } catch {
      setError("Failed to send OTP");
    }

    setLoading(false);
  }

  // ---------------------------------------------------------
  // REGISTER
  // ---------------------------------------------------------
  async function handleRegister(e) {
    e.preventDefault();
    setLoading(true);
    setError("");
    setSuccess("");

    try {
      const reg = await apiPost("/register", {
        name,
        email,
        cnic,
        otp, // <-- IMPORTANT
      });

      if (!reg.success) {
        setError(reg.message || "Registration failed");
        setLoading(false);
        return;
      }

      const { token, user, wallet } = reg.data;

      // decrypt private key
      let pkHex = "";
      try {
        pkHex = decryptPrivateKey(wallet.private_key_enc);
      } catch {
        setError("Private key decryption failed");
        setLoading(false);
        return;
      }

      // Save to localStorage
      localStorage.setItem("token", token);
      localStorage.setItem("user", JSON.stringify(user));
      localStorage.setItem("walletId", wallet.wallet_id);
      localStorage.setItem("publicKey", wallet.public_key);
      localStorage.setItem("privateKeyEnc", wallet.private_key_enc);
      localStorage.setItem("privateKeyHex", pkHex);

      setSuccess("Account created! Redirecting...");
      setTimeout(() => navigate("/login"), 1500);

    } catch {
      setError("Registration failed");
    }

    setLoading(false);
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 to-black p-6">
      <div className="bg-gray-800/40 backdrop-blur-xl border border-gray-700 shadow-2xl rounded-2xl p-10 w-full max-w-lg text-white">

        <h1 className="text-4xl font-bold text-center mb-2">Create Account</h1>
        <p className="text-gray-300 text-center mb-8">
          Join the decentralized blockchain wallet
        </p>

        <form className="space-y-5" onSubmit={handleRegister}>
          
          {/* Name */}
          <div>
            <label className="text-gray-300 font-medium">Full Name</label>
            <input
              type="text"
              className="w-full p-3 bg-gray-900 border border-gray-700 rounded-lg"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </div>

          {/* Email */}
          <div>
            <label className="text-gray-300 font-medium">Email</label>
            <input
              type="email"
              className="w-full p-3 bg-gray-900 border border-gray-700 rounded-lg"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>

          {/* CNIC */}
          <div>
            <label className="text-gray-300 font-medium">CNIC</label>
            <input
              type="text"
              placeholder="42101-1234567-1"
              className="w-full p-3 bg-gray-900 border border-gray-700 rounded-lg"
              value={cnic}
              onChange={(e) => setCnic(e.target.value)}
              required
            />
          </div>

          {/* Send OTP */}
          {!otpSent && (
            <button
              type="button"
              onClick={sendOtp}
              disabled={loading}
              className="w-full py-3 bg-blue-600 rounded-lg hover:bg-blue-700"
            >
              {loading ? "Sending OTP..." : "Send OTP"}
            </button>
          )}

          {/* OTP + Register */}
          {otpSent && (
            <>
              <div>
                <label className="text-gray-300 font-medium">Enter OTP</label>
                <input
                  type="text"
                  className="w-full p-3 bg-gray-900 border border-gray-700 rounded-lg"
                  value={otp}
                  onChange={(e) => setOtp(e.target.value)}
                  required
                />
              </div>

              <button
                type="submit"
                disabled={loading}
                className="w-full py-3 bg-green-600 rounded-lg hover:bg-green-700"
              >
                {loading ? "Creating Account..." : "Register"}
              </button>
            </>
          )}
        </form>

        {error && (
          <div className="mt-4 p-3 bg-red-600/20 border border-red-500 rounded-lg">
            {error}
          </div>
        )}

        {success && (
          <div className="mt-4 p-3 bg-green-600/20 border border-green-500 rounded-lg">
            {success}
          </div>
        )}

      </div>
    </div>
  );
}
